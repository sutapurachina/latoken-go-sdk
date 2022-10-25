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
