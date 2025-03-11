package main

import (
	"AntonKisselev/YP_GO_final/db"
	"AntonKisselev/YP_GO_final/task"
	"AntonKisselev/YP_GO_final/tests"
	"encoding/json"
	"io"
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
		Date:   date,
	}
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := task.NextDate(nowTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

type responseTasksList struct {
	Tasks []task.Task `json:"tasks"`
}

func handlerApiTask(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		task := task.Task{}
		err = json.Unmarshal(body, &task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		errSave := task.Save()
		if errSave != nil {
			errStr, err := json.Marshal(struct {
				Error string `json:"error"`
			}{
				Error: errSave.Error(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(errStr)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		} else {
			idStr, err := json.Marshal(struct {
				ID int64 `json:"id"`
			}{ID: task.Id})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(idStr)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("error: invalid method"))
}
func handlerApiGetTasks(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		search := req.URL.Query().Get("search")
		tasks, errGetAll := task.GetAll(search)
		if errGetAll != nil {
			errStr, err := json.Marshal(struct {
				Error string `json:"error"`
			}{
				Error: errGetAll.Error(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(errStr)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		}
		resp, err := json.Marshal(responseTasksList{
			Tasks: tasks,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(resp)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("error: invalid method"))
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
	http.HandleFunc("/api/task", handlerApiTask)
	http.HandleFunc("/api/tasks", handlerApiGetTasks)
	log.Println("http server started on :" + strconv.Itoa(port))
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
