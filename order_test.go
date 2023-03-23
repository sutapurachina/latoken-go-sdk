package latoken_go_sdk

import (
	"fmt"
	"testing"
)

const (
	APIKey    = "1580b1ed-69c8-4794-98a7-5ce59f136ef0"
	SecretKey = "M2UyOTUyNzktMzBjMC00MDJhLTgxNjAtNDRlMDVjM2FmMWY0"
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
