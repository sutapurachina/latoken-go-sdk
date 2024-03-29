package latoken_go_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Endpoints
const (
	baseAPIMainURL = "https://api.latoken.com"
	WSUrl          = "wss://api.latoken.com/stomp"
)

// LatokenClient defines API client
type LatokenClient struct {
	APIKey     string
	SecretKey  string
	HttpClient *http.Client
	crs        *CreateOrderService
}
type UserInfo struct {
	Id                  string      `json:"id"`
	Status              string      `json:"status"`
	Role                string      `json:"role"`
	Email               string      `json:"email"`
	Phone               string      `json:"phone"`
	Authorities         []string    `json:"authorities"`
	ForceChangePassword interface{} `json:"forceChangePassword"`
	AuthType            string      `json:"authType"`
	Socials             []string    `json:"socials"`
}
type Order struct {
	Id            string `json:"id"`
	Status        string `json:"status"`
	Side          string `json:"side"`
	Condition     string `json:"condition"`
	Type          string `json:"type"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	ClientOrderId string `json:"clientOrderId"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	Cost          string `json:"cost"`
	Filled        string `json:"filled"`
	Trader        string `json:"trader"`
	Creator       string `json:"creator"`
	CreatorId     string `json:"creatorId"`
	Timestamp     int64  `json:"timestamp"`
}
type Ticker struct {
	Symbol          string `json:"symbol"`
	BaseCurrency    string `json:"baseCurrency"`
	QuoteCurrency   string `json:"quoteCurrency"`
	Volume24H       string `json:"volume24h"`
	Volume7D        string `json:"volume7d"`
	Change24H       string `json:"change24h"`
	Change7D        string `json:"change7d"`
	Amount24H       string `json:"amount24h"`
	Amount7D        string `json:"amount7d"`
	LastPrice       string `json:"lastPrice"`
	LastQuantity    string `json:"lastQuantity"`
	BestBid         string `json:"bestBid"`
	BestBidQuantity string `json:"bestBidQuantity"`
	BestAsk         string `json:"bestAsk"`
	BestAskQuantity string `json:"bestAskQuantity"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
}

type OrderBook struct {
	Ask []struct {
		Price       string `json:"price"`
		Quantity    string `json:"quantity"`
		Cost        string `json:"cost"`
		Accumulated string `json:"accumulated"`
	} `json:"ask"`
	Bid []struct {
		Price       string `json:"price"`
		Quantity    string `json:"quantity"`
		Cost        string `json:"cost"`
		Accumulated string `json:"accumulated"`
	} `json:"bid"`
	TotalAsk string `json:"totalAsk"`
	TotalBid string `json:"totalBid"`
}

func NewLatokenClient(ApiKey, SecretKey string, httpClient *http.Client) *LatokenClient {
	res := &LatokenClient{
		APIKey:     ApiKey,
		SecretKey:  SecretKey,
		HttpClient: httpClient,
	}

	crs := &CreateOrderService{
		client: res,
	}
	res.crs = crs

	return res
}

func NewDefaultLatokenClient(ApiKey, SecretKey string) *LatokenClient {
	return NewLatokenClient(ApiKey, SecretKey, http.DefaultClient)
}

func (lc *LatokenClient) NewCreateOrderService() *CreateOrderService {
	lc.crs.Renovate()
	return lc.crs
}

func (lc *LatokenClient) GetUserInfo() (*UserInfo, error) {
	endpoint := "/v2/auth/user"
	req, err := lc.makeSignedRequest(http.MethodGet, endpoint, "", "", nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if req.Body == nil {
			return
		}
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var info UserInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (lc *LatokenClient) CancelOrder(id string) error {
	endpoint := "/v2/auth/order/cancel"
	bodyContentBytes, err := json.Marshal(struct {
		Id string `json:"id"`
	}{
		Id: id,
	})
	if err != nil {
		return err
	}
	req, err := lc.makeSignedRequest(http.MethodPost, endpoint, "", fmt.Sprintf("id=%s", id), bytes.NewReader(bodyContentBytes))
	if err != nil {
		return err
	}
	defer func() {
		if req.Body == nil {
			return
		}
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return err
	}
	if apiErr.Status != "SUCCESS" {
		return apiErr
	}
	return nil
}

func (lc *LatokenClient) CancelAllOrdersInPair(base, quote string) error {
	endpoint := "/v2/auth/order/cancelAll" + "/" + base + "/" + quote
	req, err := lc.makeSignedRequest(http.MethodPost, endpoint, "", "", nil)
	if err != nil {
		return err
	}
	defer func() {
		if req.Body == nil {
			return
		}
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return err
	}
	if apiErr.Status != "SUCCESS" {
		return apiErr
	}
	return nil
}

func (lc *LatokenClient) CancelAllOrdersInSymbol(base, quote string) error {
	endpoint := "/v2/auth/order/cancelAll" + "/" + base + "/" + quote
	req, err := lc.makeSignedRequest(http.MethodPost, endpoint, "", "", nil)
	if err != nil {
		return err
	}
	defer func() {
		if req.Body == nil {
			return
		}
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return err
	}
	if apiErr.Status != "SUCCESS" {
		return apiErr
	}
	return nil
}

func (lc *LatokenClient) makeSignedRequest(method string, endpoint, query string, signableBodyParams string, body io.Reader) (*http.Request, error) {
	signature := GetSignature(lc.SecretKey, []byte(method+endpoint+signableBodyParams))
	var (
		req *http.Request
		err error
	)
	if query != "" {
		req, err = http.NewRequest(method, baseAPIMainURL+endpoint+"?"+query, body)
	} else {
		req, err = http.NewRequest(method, baseAPIMainURL+endpoint, body)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-LA-APIKEY", lc.APIKey)
	req.Header.Add("X-LA-SIGNATURE", fmt.Sprintf("%x", signature))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-LA-DIGEST", "HMAC-SHA256")
	return req, nil
}

func (lc *LatokenClient) GetAllMyOrders(from int64, limit int) ([]Order, error) {
	endpoint := "/v2/auth/order"
	bodyContentBytes, err := json.Marshal(struct {
		From  int64 `json:"from"`
		Limit int   `json:"limit"`
	}{
		From:  from,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	signableBOdyParams := fmt.Sprintf("from=%d&limit=%d", from, limit)
	req, err := lc.makeSignedRequest(http.MethodGet, endpoint, signableBOdyParams, signableBOdyParams, bytes.NewReader(bodyContentBytes))
	if err != nil {
		return nil, err
	}
	defer func() {
		if req.Body == nil {
			return
		}
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	if resp.StatusCode < 400 {
		var res []Order
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return nil, err
	}
	return nil, apiErr
}

func (lc *LatokenClient) GetOrderBook(base, quote string, limit int) (*OrderBook, error) {
	endpoint := "/v2/book" + "/" + base + "/" + quote
	req, err := http.NewRequest(http.MethodGet, baseAPIMainURL+endpoint+"?"+fmt.Sprintf("limit=%d", limit), nil)
	if err != nil {
		return nil, err
	}
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 400 {
		var res OrderBook
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return nil, err
	}
	return nil, apiErr
}

func (lc *LatokenClient) GetTicker(base, quote string) (*Ticker, error) {
	endpoint := "/v2/ticker/" + base + "/" + quote
	req, err := http.NewRequest(http.MethodGet, baseAPIMainURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 400 {
		var res Ticker
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return nil, err
	}
	return nil, apiErr
}

func (lc *LatokenClient) GetOrder(id string) (*Order, error) {
	endpoint := "/v2/auth/order/getOrder/" + id
	req, err := lc.makeSignedRequest(http.MethodGet, endpoint, "", "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := lc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body == nil {
			return
		}
		err2 := resp.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	fmt.Println(resp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 400 {
		var res Order
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return nil, err
	}
	return nil, apiErr
}
