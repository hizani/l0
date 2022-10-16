package server

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"wbintern/l0/service/cache"

	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server
	ch cache.Cache
}

// Initialize new server
func New(address string, ch cache.Cache) *Server {
	var s *Server = &Server{&http.Server{Addr: address}, ch}
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", s.getOrderByIdHandler)
	s.Handler = router
	return s
}

func (s *Server) getOrderByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	od := s.ch.GetOrderById(id)
	switch {
	case od.Uid == "":
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Order %s not found\n", id)
		w.Write([]byte(fmt.Sprintln(msg)))
		log.Println(msg)
	default:
		template, err := template.ParseFiles("service/server/template/getById.html")
		if err != nil {
			msg := fmt.Sprintf("Error while parsing template: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(msg))
			log.Println(msg)
		}
		err = template.Execute(w, od)
		if err != nil {
			msg := fmt.Sprintf("Error while executing template: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(msg))
			log.Println(msg)
		}
	}
}
