package latoken_go_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreateOrderResponse struct {
	OrderId string `json:"id"`
	Message string `json:"message"`
}

type APIError struct {
	Result  bool              `json:"result"`
	Message string            `json:"message"`
	Err     string            `json:"error"`
	Errs    map[string]string `json:"errors"`
	Status  string            `json:"status"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("res=%v, mes=%v, err=%v, errs=%v, status=%v\n", err.Result, err.Message, err.Err, err.Errs, err.Status)
}

type CreateOrderService struct {
	client        *LatokenClient
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	Condition     string `json:"condition"`
	Side          string `json:"side"`
	Type          string `json:"type"`
	ClientOrderId string `json:"clientOrderId,omitempty"`
	Price         string `json:"price,omitempty"`
	Quantity      string `json:"quantity"`
}

func (s *CreateOrderService) Renovate() {
	s.BaseCurrency = ""
	s.QuoteCurrency = ""
	s.Condition = "GTC"
	s.Side = ""
	s.Type = ""
	s.ClientOrderId = ""
	s.Price = ""
	s.Quantity = ""
}

func (s *CreateOrderService) CreateSignableString() string {
	res := "baseCurrency" + "=" + s.BaseCurrency + "&" + "quoteCurrency" + "=" + s.QuoteCurrency + "&" + "condition" + "=" + s.Condition + "&" + "side" +
		"=" + s.Side + "&" + "type" + "=" + s.Type

	if s.ClientOrderId != "" {
		res += "&" + "clientOrderId" + "=" + s.ClientOrderId
	}

	if s.Price != "" {
		res += "&" + "price" + "=" + s.Price
	}

	res += "&" + "quantity" + "=" + s.Quantity

	return res
}

func (s *CreateOrderService) SetBase(baseCurrency string) *CreateOrderService {
	s.BaseCurrency = baseCurrency
	return s
}

func (s *CreateOrderService) SetQuote(quoteCurrency string) *CreateOrderService {
	s.QuoteCurrency = quoteCurrency
	return s
}

func (s *CreateOrderService) SetTimeInForce(condition string) *CreateOrderService {
	s.Condition = condition
	return s
}

func (s *CreateOrderService) SetSide(side string) *CreateOrderService {
	s.Side = side
	return s
}

func (s *CreateOrderService) SetType(orderType string) *CreateOrderService {
	s.Type = orderType
	return s
}

func (s *CreateOrderService) SetClientOrderId(cliOrdId string) *CreateOrderService {
	s.ClientOrderId = cliOrdId
	return s
}

func (s *CreateOrderService) SetPrice(price string) *CreateOrderService {
	s.Price = price
	return s
}

func (s *CreateOrderService) SetQuantity(quantity string) *CreateOrderService {
	s.Quantity = quantity
	return s
}

func (s *CreateOrderService) Do() (*CreateOrderResponse, error) {
	endPoint := "/v2/auth/order/place"

	bodyContentBytes, _ := json.Marshal(s)
	signableBodyParams := s.CreateSignableString()
	req, err := s.client.makeSignedRequest(http.MethodPost, endPoint, "", signableBodyParams, bytes.NewReader(bodyContentBytes))
	if err != nil {
		return nil, err
	}
	defer func() {
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err2 := req.Body.Close()
		if err == nil && err2 != nil {
			err = err2
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CreateOrderResponse
	if resp.StatusCode < 400 {
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return &response, err
	}
	var apiErr APIError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return nil, err
	}
	err = apiErr
	return &response, err
}
