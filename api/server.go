package api

import "net/http"

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) mount() {

}

func (server *Server) Run() {
	http.ListenAndServe(":8080")
}
