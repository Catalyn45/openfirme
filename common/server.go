package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	config *ServerConfig

	repository *Repository
	cache *Cache
	juridicClient *JuridicClient
}

func NewServer() *Server {
	return &Server {
		config: &config.ServerConfig,
		repository: NewRepository("file:" + config.DBFilePath + "?mode=ro"),
		cache: newCache(),
		juridicClient: NewJuridicClient(),
	}
}

func (this *Server) Start() {
	router := httprouter.New()

	static := this.cache.HtmlCache(http.FileServer(http.Dir(this.config.PublicDirectory)))

	router.Handler("GET", "/public/*filepath", http.StripPrefix("/public/", static))

	router.GET("/", this.serveHtmlFunc("index.html"))

	router.GET("/search/:nume_partial/:page_number", this.serveHtmlFunc("search.html"))
	router.GET("/api/search/:nume_partial/:page_number", this.cache.ApiPagedSearchCache(this.getFirme))

	router.GET("/profile/:numar_inmatriculare", this.serveHtmlFunc("profile.html"))
	router.GET("/api/profile/:numar_inmatriculare", this.cache.DefaultApiCache(this.getFirma))

	router.GET("/top/:page_number", this.serveHtmlFunc("topfirme.html"))
	router.GET("/api/top/:page_number", this.cache.ApiPagedCache(this.getTopFirme))

	router.GET("/admins/:cod_inmatriculare/:admin/:page_number", this.serveHtmlFunc("administratori.html"))
	router.GET("/api/admins/:cod_inmatriculare/:admin/:page_number", this.cache.ApiPagedCache(this.getAdminsFirme))

	router.GET("/dosareJuridice/:cod_inmatriculare/:page_number", this.serveHtmlFunc("dosareJuridice.html"))
	router.GET("/api/dosareJuridice/:cod_inmatriculare/:page_number", this.cache.ApiPagedCache(this.getDosareJuridiceFirma))

	router.GET("/dosarJuridic/:cod_inmatriculare/:numar_dosar", this.serveHtmlFunc("dosarJuridic.html"))
	router.GET("/api/dosarJuridic/:cod_inmatriculare/:numar_dosar", this.cache.DefaultApiCache(this.getDosarJuridicFirma))

	router.GET("/error", this.serveHtmlFunc("errorPage.html"))
	router.GET("/descarca", this.serveHtmlFunc("descarca.html"))
	router.GET("/despre", this.serveHtmlFunc("despre.html"))

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, p any) {
		err, ok := p.(error)
		if ok {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "too generic", http.StatusUnprocessableEntity)
				return
			}
		}

		log.Printf("panic: %v\n%s", p, debug.Stack())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	addr := net.JoinHostPort(this.config.Host, strconv.Itoa(this.config.Port))

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Println("Starting server...")
	log.Printf("Go to http://%s:%d in your browser.\n", this.config.Host, this.config.Port)
	log.Println("Do not close the window.")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (this *Server) returnSuccess(w http.ResponseWriter, content any) {
	w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(content)
}

func (this *Server) serveHtmlFunc(path string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	path = filepath.Join(this.config.PublicDirectory, path)

	servFunc := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, path)
	}

	return this.cache.HtmlRouterCache(servFunc)
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

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}

func (this *Server) getFilters(r *http.Request, ps httprouter.Params) (int, *FirmeFilters, *FirmeOrdering){
	pageNumber := getPageNumber(ps)

	query := r.URL.Query()

	filters := FirmeFilters {
		judet: query.Get("judet"),
		status: query.Get("status"),
		formaJuridica: query.Get("forma_juridica"),
		domeniu: query.Get("domeniu"),
	}

	numePartial := normalize(ps.ByName("nume_partial"))

	// If numePartial is a number, then we search by cui
	cui, err := strconv.Atoi(numePartial)
	if err == nil {
		filters.cui = cui
	} else {
		filters.numePartial = numePartial
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

func (this *Server) getFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageNumber, filters, _ := this.getFilters(r, ps)

	firme := this.repository.GetFirme(filters, pageNumber)

	this.returnSuccess(w, firme)
}

func (this *Server) getTopFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageNumber, filters, ordering := this.getFilters(r, ps)

	firme := this.repository.GetTopFirme(filters, ordering, pageNumber)

	this.returnSuccess(w, firme)
}

func (this *Server) getFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("numar_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	firma := this.repository.GetFirma(numar_inmatriculare);

	this.returnSuccess(w, firma)
}

func (this *Server) getAdminsFirme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	admin := ps.ByName("admin")

	pageNumber, filters, _ := this.getFilters(r, ps)

	firme := this.repository.GetAdminFirme(numar_inmatriculare, admin, filters, pageNumber)

	this.returnSuccess(w, firme)
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

func (this *Server) FilterDosare(dosare []Dosar, filters *FiltersDosare) []Dosar {
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

func (this *Server) OrderDosare(dosare *[]Dosar, sortBy string, sortOrder string) {
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

func (this *Server) getDosareJuridiceFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	pageNumber := getPageNumber(ps)

	dosare, found := this.cache.GetJuridic(numar_inmatriculare)
	if !found {
		numeFirma := this.repository.GetNumeFirma(numar_inmatriculare)
		if numeFirma == "" {
			panic(fmt.Errorf("Firma doesn't exist"))
		}

		dosare = this.juridicClient.GetDosare(numeFirma)
		this.cache.SetForJuridic(numar_inmatriculare, dosare)
	}

	query := r.URL.Query()

	filters := &FiltersDosare {
		tribunal: query.Get("tribunal"),
		categorie: query.Get("categorie"),
		stadiuProcesual: query.Get("stadiu_procesual"),
	}

	dosare = this.FilterDosare(dosare, filters)

	sort_by := query.Get("sort_by")
	sort_order := query.Get("sort_order")

	this.OrderDosare(&dosare, sort_by, sort_order)

	this.returnSuccess(w, &DosareJuridiceResponse{
		Count: len(dosare),
		Dosare: getDosareForPage(dosare, pageNumber),
	})
}

func (this *Server) getDosarJuridicFirma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numar_inmatriculare := ps.ByName("cod_inmatriculare")
	numar_inmatriculare = strings.ReplaceAll(numar_inmatriculare, "-", "/")

	numar_dosar := ps.ByName("numar_dosar")
	numar_dosar = strings.ReplaceAll(numar_dosar, "-", "/")

	dosare, found := this.cache.GetJuridic(numar_inmatriculare)
	if !found {
		numeFirma := this.repository.GetNumeFirma(numar_inmatriculare)
		if numeFirma == "" {
			panic(fmt.Errorf("Firma doesn't exist"))
		}

		dosare = this.juridicClient.GetDosare(numeFirma)
		this.cache.SetForJuridic(numar_inmatriculare, dosare)
	}

	index := slices.IndexFunc(dosare, func(dosar Dosar) bool { return dosar.Numar == numar_dosar })
	if index == -1 {
		panic(fmt.Errorf("Dosar doesn't exist"))
	}

	this.returnSuccess(w, &dosare[index])
}
