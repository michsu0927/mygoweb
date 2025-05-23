package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"web/src/controller"
	"web/src/db"
	"web/src/lib"
	"web/src/router"
	"web/src/sess"

	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	// e := echo.New()
	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// Routes
	// e.GET("/", hello)
	db.Init()

	defer db.Close()

	e := router.Init()

	middle := sess.Init()
	e.Use(middle)
	e.Use(sess.Handler)

	//middle := sess.Init()
	//e.Use(middle)
	//fmt.Printf("Environment variable PORT: %s\n", os.Getenv("PORT"))
	// Get env PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if PORT is not set
		fmt.Println("Environment variable PORT not set, using default port 3000")
	}

	go func() {
		// Create a channel to signal when RunTask is finished
		runTaskDone := make(chan bool)
		expiredTaskDone := make(chan bool)
		lib.Log("controller run task is running")
		//update for time racing
		//for range time.Tick(1 * time.Second) {
		for {
			//time.Sleep(1 * time.Second)
			// Run RunTask in a goroutine to avoid blocking
			go func() {
				rowNum := 1
				rowNumStr := os.Getenv("runTaskNum")
				var err error
				if rowNumStr == "" {
					rowNumStr = "1" // Default port if PORT is not set
					fmt.Println("Environment variable workId not set, using default")
				} else {
					rowNum, err = strconv.Atoi(rowNumStr)
					if err != nil {
						rowNum = 1
						fmt.Println("Environment variable workId is not a valid integer, using default")
					}
				}

				controller.RunTask(rowNum)
				runTaskDone <- true // Signal that RunTask has finished
			}()
			// Wait for RunTask to finish
			<-runTaskDone
			//lib.Log("runTaskDone finished")

			//if expireTime not hh:ii:ss expireTimeParts !=3 , ExpiredTask will not run
			expireTimeStr := os.Getenv("expireTime")
			if expireTimeStr == "" {
				expireTimeStr = "00:00:00"
				fmt.Println("Environment variable expireTime not set, using default expireTime 00:00:00")
			}

			// Run ExpiredTask at the time specified by expireTime everyday.
			expireTimeParts := strings.Split(expireTimeStr, ":")
			if len(expireTimeParts) == 2 {
				expireHour, _ := strconv.Atoi(expireTimeParts[0])
				expireMinute, _ := strconv.Atoi(expireTimeParts[1])
				// expireSecond, _ := strconv.Atoi(expireTimeParts[2])

				currentTime := time.Now()
				// if currentTime.Hour() == expireHour &&
				// 	currentTime.Minute() == expireMinute &&
				// 	currentTime.Second() == expireSecond {
				if currentTime.Hour() == expireHour &&
					currentTime.Minute() == expireMinute {

					lib.Log("ExpiredTask is running")
					// Run ExpiredTask in a goroutine to avoid blocking
					go func() {
						controller.ExpiredTask()
						expiredTaskDone <- true // Signal that ExpiredTask has finished
					}()

					// Wait for ExpiredTask to finish
					<-expiredTaskDone
					lib.Log("ExpiredTask finished")
				}
			}
			var runTaskSecs time.Duration = 1 * time.Second
			runTaskSecsStr := os.Getenv("runTaskSecs")
			var err error
			if runTaskSecsStr == "" {
				runTaskSecs = 1 * time.Second // Default port if PORT is not set
				fmt.Println("Environment variable workId not set, using default workId 1 second")
			} else {
				var runTaskSecsInt int
				runTaskSecsInt, err = strconv.Atoi(runTaskSecsStr)
				if err != nil {
					runTaskSecsInt = 1
					fmt.Println("Environment variable workId is not a valid integer, using default workId 1")
				}
				runTaskSecs = time.Duration(runTaskSecsInt) * time.Second
			}
			time.Sleep(runTaskSecs)
		}
	}()

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
