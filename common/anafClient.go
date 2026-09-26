package common

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type AnafClient struct {
	config *AnafClientConfig
	httpClient *http.Client
	cache *Cache
	cuisToRequest *ConcurentMap[int, string]
}

func NewAnafClient(cache *Cache) *AnafClient {
	anafConfig := &config.AnafClientConfig

	anafClient := &AnafClient{
		config: anafConfig,
		httpClient: &http.Client{
			Timeout: time.Duration(anafConfig.RequstTimeoutInSeconds) * time.Second,
		},
		cache: cache,
		cuisToRequest: NewConcurentMap[int, string](100),
	}

	if config.CacheConfig.AnafEnabled {
		go anafClient.TvaRequestsWorker()
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
		log.Println("error: ", err.Error())
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
		log.Println("Status: ", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil
	}

	var response TvaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil
	}

	return &response
}

func (this *AnafClient) MakeTvaRequest() {
	if this.cuisToRequest.IsEmpty() {
		return
	}

	log.Println("making tva request")

	cuisToRequest := this.cuisToRequest.Move()

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
}

func (this *AnafClient) TvaRequestsWorker() {
	for {
		time.Sleep(time.Duration(this.config.TvaRequestWorkerIntervalInSeconds) * time.Second)
		this.MakeTvaRequest()
	}
}

func (this *AnafClient) getTva(cui int, codInmatriculare string) *TvaInfo {
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
