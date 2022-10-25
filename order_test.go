package latoken_go_sdk

import (
	"fmt"
	"testing"
)

const (
	APIKey    = "70a111bb-0fbe-4902-8a1d-529099322bde"
	SecretKey = "ZjAwYjZlN2YtZDc2OC00ZmUzLWI3NjktMDg5ZjljOTY3M2Vl"
)

func TestCreateOrderService_Do(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	res, err := latokenClient.NewCreateOrderService().
		SetBase("ETH").
		SetQuote("USDT").
		SetTimeInForce("GTC").
		SetType("LIMIT").
		SetSide("SELL").
		SetPrice("2500").
		SetQuantity("1").Do()
	fmt.Println(err)
	fmt.Println(res)
}
