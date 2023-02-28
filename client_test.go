package latoken_go_sdk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLatokenClient_GetUserId(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)

	info, err := latokenClient.GetUserInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestLatokenClient_CancelOrder(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	res, err := latokenClient.NewCreateOrderService().
		SetBase("ETH").
		SetQuote("USDT").
		SetTimeInForce("GTC").
		SetType("LIMIT").
		SetSide("SELL").
		SetPrice("2500").
		SetQuantity("0.001").Do()
	assert.NoError(t, err, "expected successfully create order")

	err = latokenClient.CancelOrder(res.OrderId)
	assert.NoError(t, err, "expected successfully cancel order")
}

func TestLatokenClient_CancelAllOrders(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	_, err := latokenClient.NewCreateOrderService().
		SetBase("ETH").
		SetQuote("USDT").
		SetTimeInForce("GTC").
		SetType("LIMIT").
		SetSide("SELL").
		SetPrice("2500").
		SetQuantity("0.001").Do()
	assert.NoError(t, err, "expected successfully create order")
	err = latokenClient.CancelAllOrdersInPair("ETH", "USDT")
	assert.NoError(t, err, "expected successfully cancel order")
}

func TestLatokenClient_CancelAllOrdersInSymbol(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	_, err := latokenClient.NewCreateOrderService().
		SetBase("ETH").
		SetQuote("USDT").
		SetTimeInForce("GTC").
		SetType("LIMIT").
		SetSide("SELL").
		SetPrice("2500").
		SetQuantity("0.001").Do()
	assert.NoError(t, err, "expected successfully create order")

	err = latokenClient.CancelAllOrdersInSymbol("ETH", "USDT")
	assert.NoError(t, err, "expected successfully cancel order")
}

func TestLatokenClient_GetAllMyOrders(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	res, err := latokenClient.GetAllMyOrders(time.Now().UnixMilli(), 100)
	assert.NoError(t, err, "expected successfully get orders")

	if res != nil {
		fmt.Println(res)
	}
}

func TestLatokenClient_GetOrderBook(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	res, err := latokenClient.GetOrderBook("ETH", "USDT", 10)
	assert.NoError(t, err, "expected successfully get orders")

	if res != nil {
		fmt.Println(res)
	}
}

func TestLatokenClient_GetTicker(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	res, err := latokenClient.GetTicker("ETH", "USDT")
	assert.NoError(t, err, "expected successfully get orders")

	if res != nil {
		fmt.Println(res.BestAsk)
		fmt.Println(res.BestBid)
	}
}
