package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BilantToRequest struct {
	an int
	codInmatriculare string
	grup string
}

type AnafClient struct {
	config *AnafClientConfig
	httpClient *http.Client

	cache *Cache
	parser *Parser

	cuisToRequest *ConcurentMap[int, string]
	bilanturiToRequest *ConcurentMap[int, BilantToRequest]
}

func NewAnafClient(cache *Cache, repository *Repository) *AnafClient {
	anafConfig := &config.AnafClientConfig

	anafClient := &AnafClient{
		config: anafConfig,
		httpClient: &http.Client{
			Timeout: time.Duration(anafConfig.RequstTimeoutInSeconds) * time.Second,
		},
		cache: cache,
		parser: NewParser(repository),
		cuisToRequest: NewConcurentMap[int, string](100),
		bilanturiToRequest: NewConcurentMap[int, BilantToRequest](100),
	}

	if config.CacheConfig.AnafEnabled {
		go withRestart(anafClient.TvaRequestsWorker)
		go withRestart(anafClient.BilanturiRequestsWorker)
	}
	
	return anafClient
}

const anafTvaUrl = "https://webservicesp.anaf.ro/api/PlatitorTvaRest/v9/tva"

type TvaInfo struct {
	Tva  bool
	Caen string
}

type TvaInregistrareTva struct {
	Tva bool `json:"scpTVA"`
}

type TvaDateGenerale struct {
	Cui int `json:"cui"`
	Caen string `json:"cod_caen"`
}

type TvaFound struct {
	DateGenerale TvaDateGenerale `json:"date_generale"`
	InregistrareScopTva TvaInregistrareTva `json:"inregistrare_scop_Tva"`
}

type TvaResponse struct {
	Found []TvaFound `json:"found"`
}

type TvaRequest struct {
	Cui int `json:"cui"`
	Data string `json:"data"`
}

func (this *AnafClient) ConstructTvaBody(cuisToRequest map[int]string) []TvaRequest {
	requestBody := []TvaRequest{}

	now := time.Now()
	for cui, _ := range cuisToRequest {
		requestBody = append(requestBody, TvaRequest{
			Cui: cui,
			Data: now.Format("2006-01-02"),
		})
	}

	return requestBody
}

func (this *AnafClient) SendTvaRequest(request []TvaRequest) *TvaResponse {
	data, err := json.Marshal(request)
	if err != nil {
		logger.Error("error: ", err.Error())
		return nil
	}

	req, err := http.NewRequest(
		http.MethodPost,
		anafTvaUrl,
		bytes.NewBuffer(data),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := this.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Error("Status: ", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil
	}

	var response TvaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil
	}

	return &response
}

func (this *AnafClient) MakeTvaRequest() {
	cuisToRequest := this.cuisToRequest.Clone()
	defer this.cuisToRequest.RemoveAll(cuisToRequest)

	logger.Info("making tva request")

	request := this.ConstructTvaBody(cuisToRequest)

	response := this.SendTvaRequest(request)
	if response == nil {
		return
	}

	if len(response.Found) == 0 {
		return
	}

	for _, item := range response.Found {
		tvaInfo := TvaInfo {
			Tva: item.InregistrareScopTva.Tva,
			Caen: item.DateGenerale.Caen,
		}

		cui := item.DateGenerale.Cui

		this.cache.SetTva(cui, &tvaInfo)
		this.cache.RemoveProfileCache(cuisToRequest[cui])
	}

	time.Sleep(1 * time.Second)
}

func (this *AnafClient) TvaRequestsWorker() {
	for {
		time.Sleep(time.Duration(this.config.TvaRequestWorkerIntervalInSeconds) * time.Second)
		this.MakeTvaRequest()
	}
}

const anafBilanturiUrl = "https://webservicesp.anaf.ro/bilant"

type AnafBilantEntry struct {
	Indicator string `json:"indicator"`
	ValIndicator int `json:"val_indicator"`
}

type AnafBilanturiResponse struct {
	An int `json:"an"`
	Cui int `json:"cui"`
	Caen int `json:"caen"`
	I []AnafBilantEntry `json:"i"`
}

func (this *AnafClient) sendBilanturiRequest(cui int, bilanturiToRequest BilantToRequest) *AnafBilanturiResponse {
	params := url.Values{}
	params.Set("cui", strconv.Itoa(cui))
	params.Set("an", strconv.Itoa(bilanturiToRequest.an))

	fullUrl := anafBilanturiUrl + "?" + params.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		fullUrl,
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := this.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Error("Status: ", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil
	}

	var response AnafBilanturiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil
	}

	return &response
}

func (this *AnafClient) ProcessBilantResponse(response *AnafBilanturiResponse) string {
	header := "CUI,CAEN"
	values := strconv.Itoa(response.Cui) + "," + strconv.Itoa(response.Caen)

	for _, entry := range response.I {
		header += "," + entry.Indicator
		values += "," + strconv.Itoa(entry.ValIndicator)
	}

	return strings.Join([]string{header, values}, "\n")
}

func (this *AnafClient) MakeBilanturiRequest() {
	cui, bilantToRequest := this.bilanturiToRequest.Get()
	defer this.bilanturiToRequest.Remove(cui)

	logger.Info("making bilanturi request")

	response := this.sendBilanturiRequest(cui, bilantToRequest)
	if response == nil {
		return
	}

	if len(response.I) == 0 {
		logger.Error("Empty list")
		return
	}

	processed := this.ProcessBilantResponse(response)

	parsed := this.parser.ParseByGrup(strings.NewReader(processed), bilantToRequest.an, bilantToRequest.grup)

	this.cache.SetBilant(cui, bilantToRequest.an, parsed)
	this.cache.RemoveProfileCache(bilantToRequest.codInmatriculare)

	time.Sleep(1 * time.Second)
}

func (this *AnafClient) BilanturiRequestsWorker() {
	for {
		time.Sleep(time.Duration(this.config.BilanturiRequstWorkerIntervalInSeconds) * time.Second)
		this.MakeBilanturiRequest()
	}
}

func (this *AnafClient) GetTva(cui int, codInmatriculare string) *TvaInfo {
	if !config.CacheConfig.AnafEnabled {
		return nil
	}

	tvaResponse, found := this.cache.GetTva(cui)
	if found {
		return tvaResponse
	}

	this.cuisToRequest.Set(cui, codInmatriculare)

	return nil
}

func (this *AnafClient) GetBilant(cui int, codInmatriculare string, an int, grup string) map[string]int {
	if !config.CacheConfig.AnafEnabled {
		return nil
	}

	bilant, found := this.cache.GetBilant(cui, an)
	if found {
		return bilant
	}

	this.bilanturiToRequest.Set(cui, BilantToRequest{
		an: an,
		codInmatriculare: codInmatriculare,
		grup: grup,
	})

	return nil
}
