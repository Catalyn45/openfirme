package common

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	host string
	port int
	repository *Repository
}

func NewServer(host string, port int, repository *Repository) *Server {
	return &Server {
		host: host,
		port: port,
		repository: repository,
	}
}

func (self *Server) Start() {
	fmt.Println("starting server...")

	router := httprouter.New()

	static := http.FileServer(http.Dir("./public"))

	router.Handler("GET", "/public/*filepath", http.StripPrefix("/public/", static))

	router.GET("/", self.serveHtmlFunc("./public/index.html"))

	router.GET("/search", self.serveHtmlFunc("./public/search.html"))
	router.GET("/profile/:numar_inmatriculare", self.serveHtmlFunc("./public/profile.html"))

	router.GET("/firme", self.getFirme)
	router.GET("/firme/:numar_inmatriculare", self.getFirma)

	addr := net.JoinHostPort(self.host, strconv.Itoa(self.port))

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	server.ListenAndServe()
}

func (self *Server) returnSuccess(w http.ResponseWriter, content any) {
	w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(content)
}

func (self *Server) serveHtmlFunc(path string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, path)
	}
}

func (self *Server) getFirme(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()

	fmt.Println(query)

	nume_partial := query.Get("nume_partial")
	judet := query.Get("judet")
	status := query.Get("status")

	data_after := query.Get("data_after")
	data_before := query.Get("data_before")

	firme := self.repository.GetFirme(nume_partial, judet, status, data_after, data_before)

	self.returnSuccess(w, firme)
}

func (self *Server) getFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("numar_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	fmt.Println("firma:" + numar_inmatriculare)

	firma := self.repository.GetFirma(numar_inmatriculare);

	self.returnSuccess(w, firma)
}
