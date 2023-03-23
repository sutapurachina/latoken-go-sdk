package latoken_go_sdk

import (
	"fmt"
	"testing"
)

func TestLatokenClient_GetOrdersChan(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	orderC := make(chan *WsOrderUpdate, 100)
	_, _, err := latokenClient.GetOrdersChan(orderC)

	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		a := <-orderC
		fmt.Printf("aaa: %v\n", *a)
	}
}

func TestLatokenClient_GetRate(t *testing.T) {
	latokenClient := NewDefaultLatokenClient(APIKey, SecretKey)
	update := make(chan *Rate, 100)
	//_, doneC, err := latokenClient.GetRate("620f2019-33c0-423b-8a9d-cde4d7f8ef7f", "0c3a106d-bde3-4c13-a26e-3fd2394529e5", update)
	_, doneC, err := latokenClient.GetRate("MATIC", "0c3a106d-bde3-4c13-a26e-3fd2394529e5", update)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case a := <-update:
			fmt.Println("got update")
			fmt.Printf("aaa: %v\n", *a)
		case <-doneC:
			fmt.Println("closed")
			return
		}
	}
}
