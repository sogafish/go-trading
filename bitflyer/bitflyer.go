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

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
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

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	t.Timestamp = "2020-05-01T00:00:00.6005284Z" // poop
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action: DateTime, err=%s", err.Error())
	}

	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	res, err := api.sendRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	log.Printf("URL: %s, Response: %s", url, string(res))
	if err != nil {
		log.Printf("url: %s, Error: %s", url, string(res))

		return nil, err
	}
	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil {
		return nil, err
	}

	return &ticker, nil
}
