package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/tinkerbaj/redditjobs/router"
	"github.com/tinkerbaj/redditjobs/utils"
)


//Just simple main file
func main() {

	// Start the server
	if err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}

func runServer() error {
	app := gin.Default()

	app.HTMLRender = &utils.TemplRender{}

	router.SetRouter(app)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", 7070),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("Starting server on port %d\n", 7070)

	return server.ListenAndServe()
}
