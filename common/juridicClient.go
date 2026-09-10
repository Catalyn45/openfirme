package common

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type JuridicClient struct {
	httpClient *http.Client
	mu sync.Mutex
}

func NewJuridicClient() *JuridicClient {
	return &JuridicClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CautareDosare struct {
	XMLName    xml.Name `xml:"CautareDosare"`
	XMLNS      string   `xml:"xmlns,attr"`
	NumarDosar string   `xml:"numarDosar"`
	Obiect     string   `xml:"obiectDosar"`
	NumeParte  string   `xml:"numeParte"`
	Institutie *string  `xml:"institutie,omitempty"`
	DataStart  *string  `xml:"dataStart,omitempty"`
	DataStop   *string  `xml:"dataStop,omitempty"`
}

type SOAPBody struct {
	Body interface{}
}

type SOAPEnvelope struct {
	XMLNSXSI string `xml:"xmlns:xsi,attr"`
	XMLNSXSD string `xml:"xmlns:xsd,attr"`
	XMLNSSoap string `xml:"xmlns:soap,attr"`
	XMLName xml.Name `xml:"soap:Envelope"`
	Soap    string   `xml:"-"`
	Body    SOAPBody `xml:"soap:Body"`
}

const domain = "portalquery.just.ro"
const endpoint = "http://" + domain + "/query.asmx"

func (this *JuridicClient) CreateJuridicBody(body interface{}) []byte {
	reqBody := SOAPEnvelope{
		XMLNSXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		XMLNSXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNSSoap: "http://schemas.xmlsoap.org/soap/envelope/",
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
	this.mu.Lock()
	defer this.mu.Unlock()

	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(data),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", `"` + domain + `/` + action + `"`)

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
	XMLName xml.Name `xml:"DosarParte" json:"-"`
	Nume string `xml:"nume"`
	CalitateParte string `xml:"calitateParte"`
}

type Parti struct {
	XMLName xml.Name `xml:"parti" json:"-"`
	DosareParte []DosarParte `xml:"DosarParte"`
}

type DosarSedinta struct {
	XMLName xml.Name `xml:"DosarSedinta" json:"-"`
	Complet string `xml:"complet"`
	Data string `xml:"data"`
	Ora string `xml:"ora"`
	Solutie string `xml:"solutie"`
	SolutieSumar string `xml:"solutieSumar"`
	DataPronuntare string `xml:"dataPronuntare"`
}

type Sedinte struct {
	XMLName xml.Name `xml:"sedinte" json:"-"`
	DosareSedinta []DosarSedinta `xml:"DosarSedinta"`
}

type Dosar struct {
	XMLName xml.Name `xml:"Dosar" json:"-"`
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
	XMLName xml.Name `xml:"CautareDosareResult" json:"-"`
	Dosare []Dosar`xml:"Dosar"`
}

type CautareDosareResponse struct {
	XMLName    xml.Name `xml:"CautareDosareResponse"`
	CautareDosareResult CautareDosareResult
}

type CautareDosareResponseSOAPBody struct {
	Body CautareDosareResponse `xml:"CautareDosareResponse"`
}

type CautareDosareResponseSOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    string   `xml:"-"`
	Body    CautareDosareResponseSOAPBody `xml:"Body"`
}


var formeJuridiceMap = []FormaJuridicaMap {
	FormaJuridicaMap{initialForma: " S.R.L.", juridicForma: " SRL"},
	FormaJuridicaMap{initialForma: " P.F.A.", juridicForma: " PFA"},
	FormaJuridicaMap{initialForma: " P.F.", juridicForma: " PF"},
	FormaJuridicaMap{initialForma: " S.A.", juridicForma: " SA"},
	FormaJuridicaMap{initialForma: " I.I.", juridicForma: " II"},
	FormaJuridicaMap{initialForma: " C.A.", juridicForma: " CA"},
}

func (this *JuridicClient) normalizeNumeFirma(numeFirma string) string {
	for _, formaJuridica := range formeJuridiceMap {
		numeFirma = strings.ReplaceAll(numeFirma, formaJuridica.initialForma, formaJuridica.juridicForma)
	}

	return numeFirma
}

func (this *JuridicClient) GetDosare(numeFirma string) []Dosar {
	numeFirma = this.normalizeNumeFirma(numeFirma)

	requestData := this.CreateJuridicBody(
		CautareDosare{
			XMLNS:      domain,
			NumarDosar: "",
			Obiect:     "",
			NumeParte:  numeFirma,
		},
	)

	responseData := this.sendJuridicRequest("CautareDosare", requestData)

	var response CautareDosareResponseSOAPEnvelope
	err := xml.Unmarshal(responseData, &response)
	if err != nil {
		panic(err)
	}

	return response.Body.Body.CautareDosareResult.Dosare
}
