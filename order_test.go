package latoken_go_sdk

import (
	"fmt"
	"testing"
)

const (
	APIKey    = "fad02b9a-6c82-4dae-93f7-1596a3098242"
	SecretKey = "Y2Q3ZDA4MDgtNzA2ZC00MjJiLWI3NjMtMjkwYmNjNzcwNDRl"
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
