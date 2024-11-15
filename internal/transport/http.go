package transport

import (
	"encoding/json"
	"goland-tour-api/internal/todo"
	"log"
	"net/http"
)

type TodoItem struct {
	Id   int64  `json:"id"`
	Item string `json:"item"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSvc *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		items, err := todoSvc.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(items)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		var t TodoItem
		//var maxId int64 = 1

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// For extra credit delete method.
		// Kept out for now because it conflicts with the rest of the course
		//for x := range todos {
		//	if todos[x].Id > maxId {
		//		maxId = todos[x].Id
		//	}
		//}
		//t.Id = maxId + 1

		err = todoSvc.Add(t.Item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	})

	// Removed as it conflicts with the rest of the course
	//mux.HandleFunc("DELETE /t", func(w http.ResponseWriter, r *http.Request) {
	//	var t TodoItem
	//	err := json.NewDecoder(r.Body).Decode(&t)
	//	if err != nil {
	//		log.Println(err)
	//		w.WriteHeader(http.StatusBadRequest)
	//	}
	//
	//	for x := range todos {
	//		if todos[x].Id == t.Id {
	//			todos = append(todos[:x], todos[x+1:]...)
	//			w.WriteHeader(http.StatusOK)
	//			return
	//		}
	//	}
	//
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := todoSvc.Search(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
