package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	host string
	port int

	repository *Repository
	cache *Cache
	juridicClient *JuridicClient
}

func NewServer(host string, port int, repository *Repository) *Server {
	return &Server {
		host: host,
		port: port,
		repository: repository,
		cache: newCache(),
		juridicClient: NewJuridicClient(),
	}
}

func (self *Server) Start() {
	fmt.Println("starting server...")

	router := httprouter.New()

	static := self.cache.HtmlCache(http.FileServer(http.Dir("./public")))

	router.Handler("GET", "/public/*filepath", http.StripPrefix("/public/", static))

	router.GET("/", self.serveHtmlFunc("./public/index.html"))

	router.GET("/profile/:numar_inmatriculare", self.serveHtmlFunc("./public/profile.html"))

	router.GET("/search/:nume_partial/:page_number", self.serveHtmlFunc("./public/search.html"))
	router.GET("/top/:page_number", self.serveHtmlFunc("./public/topfirme.html"))
	router.GET("/admins/:cod_inmatriculare/:admin/:page_number", self.serveHtmlFunc("./public/administratori.html"))
	router.GET("/dosareJuridice/:cod_inmatriculare/:page_number", self.serveHtmlFunc("./public/dosareJuridice.html"))
	router.GET("/dosarJuridic/:cod_inmatriculare/:numar_dosar", self.serveHtmlFunc("./public/dosarJuridic.html"))
	router.GET("/error", self.serveHtmlFunc("./public/errorPage.html"))

	router.GET("/firme/:nume_partial/:page_number", self.cache.ApiPagedSearchCache(self.getFirme))
	router.GET("/firma/:numar_inmatriculare", self.cache.DefaultApiCache(self.getFirma))
	router.GET("/topFirme/:page_number", self.cache.ApiPagedCache(self.getTopFirme))
	router.GET("/adminsFirme/:cod_inmatriculare/:admin/:page_number", self.cache.ApiPagedCache(self.getAdminsFirme))
	router.GET("/dosareJuridiceFirma/:cod_inmatriculare/:page_number", self.cache.ApiPagedCache(self.getDosareJuridiceFirma))
	router.GET("/dosarJuridicFirma/:cod_inmatriculare/:numar_dosar", self.cache.DefaultApiCache(self.getDosarJuridicFirma))

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, p any) {
		err, ok := p.(error)
		if ok {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "too generic", http.StatusUnprocessableEntity)
				return
			}
		}

		fmt.Printf("panic: %v\n%s", p, debug.Stack())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	addr := net.JoinHostPort(self.host, strconv.Itoa(self.port))

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (self *Server) returnSuccess(w http.ResponseWriter, content any) {
	w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(content)
}

func (self *Server) serveHtmlFunc(path string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	servFunc := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, path)
	}

	return self.cache.HtmlRouterCache(servFunc)
}

func getPageNumber(ps httprouter.Params) int {
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

	return pageNumber
}

func (self *Server) getFilters(r *http.Request, ps httprouter.Params) (int, *FirmeFilters, *FirmeOrdering){
	pageNumber := getPageNumber(ps)
	fmt.Println(pageNumber)

	query := r.URL.Query()

	fmt.Println(query)

	filters := FirmeFilters {
		numePartial: normalize(ps.ByName("nume_partial")),
		judet: query.Get("judet"),
		status: query.Get("status"),
		formaJuridica: query.Get("forma_juridica"),
	}

	ordering := FirmeOrdering{}

	sort_by := query.Get("sort_by")

	if sort_by != "" && filters.numePartial != "" {
		panic(fmt.Errorf("can't have both"))
	}

	if sort_by != "" {
		sort := "asc"

		if query.Get("sort_order") == "desc" {
			sort = "desc"
		}

		ordering.sortBy = sort_by
		ordering.sortOrder = sort
	}

	return pageNumber, &filters, &ordering
}

func (self *Server) getFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageNumber, filters, _ := self.getFilters(r, ps)

	firme := self.repository.GetFirme(filters, pageNumber)

	self.returnSuccess(w, firme)
}

var allowedFormeJuridiceForBilanturi = []string{"SRL", "SA", "SNC", "SCS", "SCA", "RA"}

func (self *Server) getTopFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageNumber, filters, ordering := self.getFilters(r, ps)

	firme := self.repository.GetTopFirme(filters, ordering, pageNumber)

	self.returnSuccess(w, firme)
}

func (self *Server) getFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("numar_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	fmt.Println("firma:" + numar_inmatriculare)

	firma := self.repository.GetFirma(numar_inmatriculare);

	self.returnSuccess(w, firma)
}

func (self *Server) getAdminsFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	fmt.Println("firma:" + numar_inmatriculare)

	pageNumber := getPageNumber(ps)
	fmt.Println(pageNumber)

	admin := ps.ByName("admin")

	firme := self.repository.GetAdminFirme(numar_inmatriculare, admin, pageNumber)

	self.returnSuccess(w, firme)
}

func getDosareForPage(dosare []Dosar, pageNumber int) []Dosar {
	perPage := 20

	start := (pageNumber - 1) * perPage
	if start >= len(dosare) {
		return []Dosar{}
	}

	end := min(len(dosare)-1, start+perPage)

	return dosare[start:end]
}

type DosareJuridiceResponse struct {
	Count int
	Dosare []Dosar
}

type FiltersDosare struct {
	tribunal string
	categorie string
	stadiuProcesual string
}

func (self *Server) FilterDosare(dosare []Dosar, filters *FiltersDosare) []Dosar {
	if filters.tribunal == "" && filters.categorie == "" && filters.stadiuProcesual == "" {
		return dosare
	}

	dosareFiltered := []Dosar{}
	for _, dosar := range dosare {
		if filters.tribunal != "" && dosar.Institutie != filters.tribunal {
			continue
		}

		if filters.categorie != "" && dosar.CategorieCazNume != filters.categorie {
			continue
		}

		if filters.stadiuProcesual != "" && dosar.StadiuProcesualNume != filters.stadiuProcesual {
			continue
		}

		dosareFiltered = append(dosareFiltered, dosar)
	}

	return dosareFiltered
}

func (self *Server) OrderDosare(dosare *[]Dosar, sortBy string, sortOrder string) {
	if sortBy == "" {
		return
	}

	if sortOrder == "" {
		sortOrder = "asc"
	}

	sort.Slice(*dosare, func(i, j int) bool {
		if sortBy == "data" {
			if sortOrder == "asc" {
				return (*dosare)[i].Data < (*dosare)[j].Data
			}

			return (*dosare)[i].Data > (*dosare)[j].Data
		}

		if sortBy == "parti" {
			if sortOrder == "asc" {
				return len((*dosare)[i].Parti.DosareParte) < len((*dosare)[j].Parti.DosareParte)
			}

			return len((*dosare)[i].Parti.DosareParte) > len((*dosare)[j].Parti.DosareParte)
		}

		panic(fmt.Errorf("invalid sortBy"))
	})
}

type FormaJuridicaMap struct {
	initialForma string
	juridicForma string
}

func (self *Server) getDosareJuridiceFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	fmt.Println("firma:" + numar_inmatriculare)

	pageNumber := getPageNumber(ps)
	fmt.Println(pageNumber)

	dosare, found := self.cache.GetJuridic(numar_inmatriculare)
	if !found {
		numeFirma := self.repository.GetNumeFirma(numar_inmatriculare)
		if numeFirma == "" {
			panic(fmt.Errorf("Firma doesn't exist"))
		}

		dosare = self.juridicClient.GetDosare(numeFirma)
		self.cache.SetForJuridic(numar_inmatriculare, dosare)
	}

	query := r.URL.Query()

	filters := &FiltersDosare {
		tribunal: query.Get("tribunal"),
		categorie: query.Get("categorie"),
		stadiuProcesual: query.Get("stadiu_procesual"),
	}

	dosare = self.FilterDosare(dosare, filters)

	sort_by := query.Get("sort_by")
	sort_order := query.Get("sort_order")

	self.OrderDosare(&dosare, sort_by, sort_order)

	self.returnSuccess(w, &DosareJuridiceResponse{
		Count: len(dosare),
		Dosare: getDosareForPage(dosare, pageNumber),
	})
}

func (self *Server) getDosarJuridicFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	numar_dosar := ps.ByName("numar_dosar")
	numar_dosar = strings.ReplaceAll(numar_dosar, "-", "/")

	dosare, found := self.cache.GetJuridic(numar_inmatriculare)
	if !found {
		numeFirma := self.repository.GetNumeFirma(numar_inmatriculare)
		if numeFirma == "" {
			panic(fmt.Errorf("Firma doesn't exist"))
		}

		dosare = self.juridicClient.GetDosare(numeFirma)
		self.cache.SetForJuridic(numar_inmatriculare, dosare)
	}

	index := slices.IndexFunc(dosare, func(dosar Dosar) bool { return dosar.Numar == numar_dosar })
	if index == -1 {
		panic(fmt.Errorf("Dosar doesn't exist"))
	}

	self.returnSuccess(w, &dosare[index])
}
