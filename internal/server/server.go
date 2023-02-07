package server

import (
	json2 "encoding/json"
	"github.com/sleepycrew/appmonitor-client/internal/monitor"
	"io"
	"net/http"
	"os"
)

type handler func(w http.ResponseWriter, r *http.Request)

func StartServer(monitor monitor.MonitorInterface) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/monitor", monitorHandler(monitor))

	http.ListenAndServe(":"+port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hewo!\n")
}

func monitorHandler(monitor monitor.MonitorInterface) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := monitor.Execute()
		json, _ := json2.Marshal(resp)
		w.Write(json)
	}
}
