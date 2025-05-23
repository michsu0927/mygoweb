package router

import (
	"fmt"
	"io"
	"os"

	//"fmt"

	"web/src/controller"
	"web/src/tpl" // template

	//"text/template" //text template
	"net/http"
	echopprof "web/src/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware" //v4 use v4 middleware
)

func GetIP(c echo.Context) string {

	// Open the log file in append mode
	logFile, err := os.OpenFile("log/ip.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // Handle error as appropriate for your application
	}
	defer logFile.Close()

	// Create a multiwriter to write to both the file and standard output
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Log all request headers to the file
	for name, headers := range c.Request().Header {
		for _, h := range headers {
			// Append header information to the log file
			_, err := fmt.Fprintf(multiWriter, "Header %v: %v\n", name, h)
			if err != nil {
				panic(err)
			}
		}
	}

	ip := ""

	// Prioritize headers in order: X-Forwarded-For, X-Real-IP, Sec-User-Ip
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "Sec-User-Ip"} {
		ip = c.Request().Header.Get(header)
		if ip != "" {
			c.Logger().Info("Using header: %s, IP: %s", header, ip)
			return ip
		}
	}

	if ip != "" {
		// If none of the preferred headers are found, fallback to c.RealIP()
		ip = c.RealIP()
	}

	c.Logger().Info("Using c.RealIP(), IP: %s", ip)
	return ip

}

// Init Router
func Init() *echo.Echo {

	e := echo.New()
	file, err := os.OpenFile("log/ip.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	e.Logger.SetOutput(file)
	e.Pre(middleware.AddTrailingSlash()) //request uri 最後面自動加 /

	t := tpl.Init()

	e.Renderer = t
	/*
		https://echo.labstack.com/guide/routing
		I would like to apply the same middlewares for all /admin/* URLs. Is it possible ?
		g := e.Group("/admin", <your-middleware>)
		g.GET("/secured", <your-handler>)
		Now route /admin/secured will execute your-middleware.
		https://github.com/labstack/echo/issues/613
	*/
	e.GET("/hello/:page/*", controller.Demo)

	e.GET("/hello/*", controller.Demo)

	//e.GET("/", api.Home)
	e.GET("/", func(c echo.Context) error {
		//fmt.Printf("/ Hello")
		//c.Logger().Info("/")
		userIP := GetIP(c)
		output := "Hello, World!" + userIP
		//get data from DB
		// Suggested code may be subject to a license. Learn more: ~LicenseLog:4158231234.
		// DB := db.Manager()
		// DB.q
		return c.String(http.StatusOK, output)
	})

	e.Any("/user/:name", func(c echo.Context) error {
		name := c.Param("name")
		output := fmt.Sprintf("Hello, %s!", name)
		return c.String(http.StatusOK, output)
	})

	e.POST("/api/addTask/", controller.AddTask)

	//tasks
	e.GET("/api/tasks/:taskIDoruserID/*", controller.Tasks)
	e.GET("/api/tasks/*", controller.Tasks)

	//records
	e.GET("/api/records/:userID/*", controller.Records)
	e.GET("/api/records/*", controller.Records)

	//userPointBalance balance
	e.GET("/api/balance/:userID/*", controller.Users)
	e.GET("/api/balance/*", controller.Users)

	//所有 public folder 的文件都可以被 access ,例如 public/robots.txt -> http://yourhost/robots.txt
	e.Static("/", "public")

	//print payload
	e.Any("/api/print_all/", controller.PrintPayload)

	// EX:go tool pprof http://192.168.119.128/
	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	echopprof.Wrap(e)

	return e
}
