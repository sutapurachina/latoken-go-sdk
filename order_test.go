package latoken_go_sdk

import (
	"fmt"
	"testing"
)

const (
	APIKey    = "bc6eb3ed-0a32-4f19-8e56-e2359a209492"
	SecretKey = "YmQyZTRhNWItZWE0Ni00N2Q0LTg4MmYtYjlkMTk5MDQ1ZGE5"
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
