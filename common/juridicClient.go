package common

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type JuridicClient struct {
	cache *Cache
	repository *Repository
	httpClient *http.Client
}

func NewJuridicClient(cache *Cache, repository *Repository) *JuridicClient {
	return &JuridicClient{
		cache: cache,
		repository: repository,
		httpClient: &http.Client{
			Timeout: time.Duration(config.DosareJuridiceClientConfig.RequstTimeoutInSeconds) * time.Second,
		},
	}
}

type CautareDosare struct {
	XMLName    xml.Name `xml:"PagedSearchDocket"`
	XMLNS      string   `xml:"xmlns,attr"`

	NumarDosar *string   `xml:"numarDosar,omitempty"`
	Obiect     *string   `xml:"obiectDosar,omitempty"`
	NumeParte  *string   `xml:"numeParte,omitempty"`
	Institutie *string  `xml:"institutie,omitempty"`
	DataStart  *string  `xml:"dataStart,omitempty"`
	DataStop   *string  `xml:"dataStop,omitempty"`

	Page int `xml:"page"`
	RowsPerPage int `xml:"rowsPerPage"`
}

type SOAPBody struct {
	Body interface{}
}

type SOAPEnvelope struct {
	XMLNSXSI string `xml:"xmlns:xsi,attr"`
	XMLNSXSD string `xml:"xmlns:xsd,attr"`
	XMLNSSoap string `xml:"xmlns:soap12,attr"`
	XMLName xml.Name `xml:"soap12:Envelope"`
	Soap    string   `xml:"-"`
	Body    SOAPBody `xml:"soap12:Body"`
}

const domain = "http://tempuri.org/"
const endpoint = "http://portalquery.just.ro/QueryDocket.asmx"

func (this *JuridicClient) CreateJuridicBody(body interface{}) []byte {
	reqBody := SOAPEnvelope{
		XMLNSXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		XMLNSXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNSSoap: "http://www.w3.org/2003/05/soap-envelope",
		Body: SOAPBody{
			Body: body,
		},
	}

	data, err := xml.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	data = append(
		[]byte(`<?xml version="1.0" encoding="utf-8"?>`),
		data...
	)

	return data
}

func (this *JuridicClient) sendJuridicRequest(action string, data []byte) []byte {
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(data),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	resp, err := this.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

type DosarParte struct {
	XMLName xml.Name `xml:"DosarParteDocket" json:"-"`
	Nume string `xml:"nume"`
	CalitateParte string `xml:"calitateParte"`
}

type Parti struct {
	XMLName xml.Name `xml:"parti" json:"-"`
	DosareParte []DosarParte `xml:"DosarParteDocket"`
}

type DosarSedinta struct {
	XMLName xml.Name `xml:"DosarSedintaDocket" json:"-"`
	Complet string `xml:"complet"`
	Data string `xml:"data"`
	Ora string `xml:"ora"`
	Solutie string `xml:"solutie"`
	SolutieSumar string `xml:"solutieSumar"`
	DataPronuntare string `xml:"dataPronuntare"`
}

type Sedinte struct {
	XMLName xml.Name `xml:"sedinte" json:"-"`
	DosareSedinta []DosarSedinta `xml:"DosarSedintaDocket"`
}

type Dosar struct {
	XMLName xml.Name `xml:"DosarDocket" json:"-"`
	Parti Parti
	Sedinte Sedinte
	Numar string `xml:"numar"`
	Data string `xml:"data"`
	Institutie string `xml:"institutie"`
	Departament string `xml:"departament"`
	Obiect string `xml:"obiect"`
	CategorieCazNume string `xml:"categorieCazNume"`
	StadiuProcesualNume string `xml:"stadiuProcesualNume"`
}

type CautareDosareResult struct {
	XMLName xml.Name `xml:"PagedSearchDocketResult" json:"-"`
	Dosare []Dosar`xml:"DosarDocket"`
}

type CautareDosareResponse struct {
	XMLName    xml.Name `xml:"PagedSearchDocketResponse"`
	CautareDosareResult CautareDosareResult
}

type CautareDosareResponseSOAPBody struct {
	Body CautareDosareResponse `xml:"PagedSearchDocketResponse"`
}

type CautareDosareResponseSOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    string   `xml:"-"`
	Body    CautareDosareResponseSOAPBody `xml:"Body"`
}

func (this *JuridicClient) normalizeNumeFirma(numeFirma string) string {
	result := []string{}

	for _, word := range strings.Fields(numeFirma) {
		wordWithoutPoints := strings.ReplaceAll(word, ".", "")
		wordWithoutPoints = strings.ToUpper(wordWithoutPoints)

		if slices.Contains(allFormeJuridice, wordWithoutPoints) {
			word = wordWithoutPoints
		}

		result = append(result, word)
	}

	return strings.Join(result, " ")
}

func (this *JuridicClient) GetDosare(codInmatriculare string) []Dosar {
	dosare, found := this.cache.GetJuridic(codInmatriculare)
	if found {
		return dosare
	}

	numeFirma := this.repository.GetNumeFirma(codInmatriculare)
	if numeFirma == "" {
		panic(fmt.Errorf("Firma doesn't exist"))
	}

	numeFirma = this.normalizeNumeFirma(numeFirma)
	log.Println("Getting dosare for: ", numeFirma)

	requestData := this.CreateJuridicBody(
		CautareDosare{
			XMLNS:      domain,
			NumarDosar: nil,
			Obiect:     nil,
			NumeParte:  &numeFirma,
			Page: 0,
			RowsPerPage: 200,
		},
	)

	responseData := this.sendJuridicRequest("CautareDosare", requestData)

	var response CautareDosareResponseSOAPEnvelope
	err := xml.Unmarshal(responseData, &response)
	if err != nil {
		panic(err)
	}

	dosare = response.Body.Body.CautareDosareResult.Dosare

	this.cache.SetForJuridic(codInmatriculare, dosare)

	return dosare
}
