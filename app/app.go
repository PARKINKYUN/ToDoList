package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rd *render.Render

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var todoMap map[int]*Todo

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	rd.JSON(w, http.StatusOK, list)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	id := len(todoMap) + 1
	todoMap[id] = &Todo{id, name, false, time.Now()}
	rd.JSON(w, http.StatusCreated, todoMap[id])
}

type Success struct {
	Success bool `json:"success"`
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := todoMap[id]; ok {
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusBadRequest, Success{false})
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusBadRequest, Success{false})
	}
}

func MakeHandler() http.Handler {
	todoMap = make(map[int]*Todo)

	rd = render.New()
	router := mux.NewRouter()

	router.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")
	router.HandleFunc("/", indexHandler)

	return router
}
