package main

import (
	"AntonKisselev/YP_GO_final/db"
	"AntonKisselev/YP_GO_final/task"
	"AntonKisselev/YP_GO_final/tests"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handlerApiNextDate(w http.ResponseWriter, req *http.Request) {
	now := req.URL.Query().Get("now")
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")

	if now == "" || date == "" || repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	task := task.Task{
		Repeat: repeat,
	}
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := task.NextDate(nowTime, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func main() {
	err := db.CheckDb()
	if err != nil {
		log.Fatal(err)
		return
	}

	webDir := "./web"
	port, err := strconv.Atoi(os.Getenv("TODO_PORT"))
	if err != nil {
		port = tests.Port
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlerApiNextDate)
	log.Println("http server started on :" + strconv.Itoa(port))
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
