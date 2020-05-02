package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseUrl string = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}

	return apiClient
}

func (api *APIClient) header(method, endpoint string, body []byte) map[string]string {
	timtestamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(timtestamp)
	message := timtestamp + method + endpoint + string(body)

	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timtestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (api *APIClient) sendRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return
	}

	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}

	endpoint := baseURL.ResolveReference(apiURL).String()
	log.Printf("Action: sendRequest, Endpoint: %s", endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))

	if err != nil {
		return
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}

	res, err := api.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	res, err := api.sendRequest("GET", url, map[string]string{}, nil)
	log.Printf("URL: %s, Response: %s", url, string(res))
	if err != nil {
		log.Printf("url: %s, Error: %s", url, string(res))

		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(res, &balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
