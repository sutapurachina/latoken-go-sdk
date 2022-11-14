package latoken_go_sdk

import (
	"fmt"
	"testing"
)

func TestLatokenClient_GetOrdersChan(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)

	update, _, _, err := latokenClient.GetOrdersChan()

	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		a := <-update
		fmt.Printf("aaa: %v\n", *a)
	}
}

func TestLatokenClient_GetRate(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)

	update, _, doneC, err := latokenClient.GetRate("BTC", "USDT")

	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case <-update:
			fmt.Println("got update")
			//fmt.Printf("aaa: %v\n", *a)
		case <-doneC:
			fmt.Println("closed")
			return
		}
	}
	return
}
