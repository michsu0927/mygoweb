package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Log logs a message.
// parametter Text log string ,parametter logFile if not give using logFileName(default:log/app.log) to save or using logFile to save. if logFile is "" it will log to console.
func Log(Text string, logFile ...string) {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	logFileName := fmt.Sprintf("log/%s.log", dateStr)
	var writer io.Writer

	if len(logFile) > 0 {

		if logFile[0] == "" {
			writer = os.Stdout
		} else {
			logFileName = fmt.Sprintf("log/%s%s.log", dateStr, logFile[0])
			file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("Open log file error: ", err)
				writer = os.Stdout
			} else {
				defer file.Close()
				writer = io.MultiWriter(os.Stdout, file)
			}

		}
	} else {
		file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Open log file error: ", err)
			writer = os.Stdout
		} else {
			defer file.Close()
			writer = io.MultiWriter(os.Stdout, file)
		}
	}

	logger := log.New(writer, "", log.LstdFlags)

	logMessage := fmt.Sprintf("%s : %s", now.Format("2006-01-02 15:04:05"), Text)
	logger.Println(logMessage)
}
