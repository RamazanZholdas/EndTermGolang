package loggs

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("Status code: %d\nClients ip address: %s\nMethod: %s\nPath: \"%s\"\nLatency: %s\nExact time: %s\n\n",
			params.StatusCode, params.ClientIP, params.Method, params.Path, params.Latency, params.TimeStamp.Format(time.RFC850))
	})
}

func SetupLogOutput() os.File {
	file, err := os.Create("gin.log")
	if err != nil {
		log.Fatalf("Cannot create a file, \nError:%v", err)
	}
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	return *file
}
