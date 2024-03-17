package main

import (
	"gotimer/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	migratorApp := app.GetMigratorApp()
	schedulerApp := app.GetSchedulerApp()
	webServer := app.GetWebServer()

	migratorApp.Start()
	schedulerApp.Start()
	defer schedulerApp.Stop()

	webServer.Start()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_ = http.ListenAndServe(":9999", nil)
		})
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
}
