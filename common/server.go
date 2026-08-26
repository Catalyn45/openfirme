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

	router.GET("/search/:page_number", self.serveHtmlFunc("./public/search.html"))
	router.GET("/profile/:numar_inmatriculare", self.serveHtmlFunc("./public/profile.html"))
	router.GET("/top/:page_number", self.serveHtmlFunc("./public/topfirme.html"))

	router.GET("/firme/:page_number", self.getFirme)
	router.GET("/firma/:numar_inmatriculare", self.getFirma)

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

func (self *Server) getFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	page_number := ps.ByName("page_number")

	pageNumber := 1

	if page_number != "" {
		var err error

		pageNumber, err = strconv.Atoi(page_number)
		if err != nil {
			panic(err)
		}
	}

	if pageNumber > 10 {
		panic(fmt.Errorf("can't have more than 200 results"))
	}

	fmt.Println(pageNumber)

	query := r.URL.Query()

	fmt.Println(query)

	filters := FirmeFilters {
		numePartial: query.Get("nume_partial"),
		judet: query.Get("judet"),
		status: query.Get("status"),
		formaJuridica: query.Get("forma_juridica"),
		dataAfter: query.Get("data_after"),
		dataBefore: query.Get("data_before"),
	}

	ordering := FirmeOrdering{}

	sort_by := query.Get("sort_by")
	if sort_by != "" {
		sort := "asc"

		if query.Get("sort_order") == "desc" {
			sort = "desc"
		}

		ordering.sortBy = sort_by
		ordering.sortOrder = sort
	}

	firme := self.repository.GetFirme(&filters, &ordering, pageNumber)

	self.returnSuccess(w, firme)
}

func (self *Server) getFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("numar_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	fmt.Println("firma:" + numar_inmatriculare)

	firma := self.repository.GetFirma(numar_inmatriculare);

	self.returnSuccess(w, firma)
}
