package latoken_go_sdk

import (
	"encoding/json"
	"fmt"
	stompws "github.com/drawdy/stomp-ws-go"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
	"time"
)

type Header struct {
	ApiKey    string `json:"X-LA-APIKEY"`
	Digest    string `json:"X-LA-DIGEST,omitempty"`
	SigData   string `json:"X-LA-SIGDATA"`
	Signature string `json:"X-LA-SIGNATURE"`
}

type WSOrderAns struct {
	Payload   []WsOrderUpdate `json:"payload"`
	Nonce     int             `json:"nonce"`
	Timestamp int64           `json:"timestamp"`
}

type Rate struct {
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
}

type WsOrderUpdate struct {
	Id            string `json:"id"`
	User          string `json:"user"`
	ChangeType    string `json:"changeType"`
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
	DeltaFilled   string `json:"deltaFilled"`
	Timestamp     int64  `json:"timestamp"`
	RejectError   string `json:"rejectError"`
	RejectComment string `json:"rejectComment"`
	Creator       string `json:"creator"`
	CreatorId     string `json:"creatorId"`
}

func (lc *LatokenClient) GetOrdersChan(orderC chan *WsOrderUpdate) (chan struct{}, chan struct{}, error) {
	info, err := lc.GetUserInfo()
	if err != nil {
		return nil, nil, err
	}
	endPoint := "/user/" + info.Id + "/v1/order"
	c, _, err := websocket.DefaultDialer.Dial(WSUrl, nil)
	c.SetReadLimit(650000)
	if err != nil {
		return nil, nil, err
	}

	signatureTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := GetSignature(lc.SecretKey, []byte(signatureTime))

	a := stompws.Message{
		Command: "CONNECT",
		Headers: stompws.Headers{
			"X-LA-APIKEY", lc.APIKey,
			"X-LA-DIGEST", "HMAC-SHA256",
			"X-LA-SIGDATA", signatureTime,
			"X-LA-SIGNATURE", fmt.Sprintf("%x", signature),
		},
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s\n%s\n\x00\n", a.Command, a.Headers.String())))
	if err != nil {
		return nil, nil, err
	}
	_, _, err = c.ReadMessage()
	if err != nil {
		fmt.Printf("getorderschan: can't read message`: %v\n", err)
	}
	go keepAlive(c, 10*time.Second)
	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SUBSCRIBE\nid:%d\ndestination:%s\nack:auto\n\n\x00\n", 0, endPoint)))
	if err != nil {
		return nil, nil, err
	}
	_, _, err = c.ReadMessage()
	if err != nil {
		fmt.Printf("getorderschan: can't read message`: %v\n", err)
	}
	stopC := make(chan struct{})
	doneC := make(chan struct{})
	go func() {
		select {
		case <-stopC:
			err = c.Close()
			if err != nil {
				fmt.Printf("getorderschan: can't close channel: %v\n", err)
			}
			doneC <- struct{}{}
		}
	}()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("getorderschan: can't read message`: %v\n", err)
				if strings.Contains(err.Error(), "clos") {
					break
				}
				continue
			}
			for idx, r := range message {
				if r == '{' {
					message = message[idx : len(message)-1]
					break
				}
			}
			var ans WSOrderAns
			err = json.Unmarshal(message, &ans)
			if err != nil {
				log.Printf("getorderschan: can't unmarshal message`: %v\n", err)
				continue
			}
			for _, update := range ans.Payload {
				orderC <- &update
			}
		}
	}()

	return stopC, doneC, nil
}

func (lc *LatokenClient) GetRate(base, quote string, update chan *Rate) (chan struct{}, chan struct{}, error) {
	endPoint := "/v1/rate/" + base + "/" + quote
	fmt.Println(endPoint)
	//endPoint := "/v1/ticker"
	c, _, err := websocket.DefaultDialer.Dial(WSUrl, nil)
	c.SetReadLimit(6555350)
	if err != nil {
		return nil, nil, err
	}

	//signatureTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	//signature := GetSignature(lc.SecretKey, []byte(signatureTime))

	a := stompws.Message{
		Command: "CONNECT",
		Headers: stompws.Headers{
			stompws.HK_ACCEPT_VERSION, "1.2",
			stompws.HK_HEART_BEAT, "1000,1000",
			stompws.HK_HOST, WSUrl,
		},
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s\n%s\n\x00\n", a.Command, a.Headers.String())))
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return nil, nil, err
	}
	_, mes, err := c.ReadMessage()
	fmt.Println(string(mes))
	if err != nil {
		fmt.Printf("gettickerschan: can't read message1`: %v\n", err)
	}
	keepAlive(c, 500*time.Millisecond)
	//go keepAlive(c, 10*time.Second)
	//err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SUBSCRIBE\nid:%d\ndestination:%s\nack:auto\n\n\x00\n", 1, endPoint)))
	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SUBSCRIBE\nid:%d\ndestination:%s\nsubscription:4\nack:client\n\n\x00\n", 0, endPoint)))
	if err != nil {
		return nil, nil, err
	}
	stopC := make(chan struct{})
	doneC := make(chan struct{})
	//lastResponse := time.Now()
	go func() {
		select {
		case <-stopC:
			fmt.Println("got signal")
			/*err = c.Close()
			if err != nil {
				fmt.Printf("gettickerchan: can't close channel: %v\n", err)
			}*/
			doneC <- struct{}{}
			return

		}
	}()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
			}
			if err != nil {
				log.Printf("gettickerchan: can't read message2`: %v\n", err)
				if strings.Contains(err.Error(), "use of closed network connection") {
					continue
				}
				if strings.Contains(err.Error(), "clos") {
					stopC <- struct{}{}
					fmt.Println("sent signal")
					break
				}
				continue
			}
			fmt.Println(string(message))
			if len(message) < 2 {
				//deadline := time.Now().Add(2 * time.Second)
				//err := c.WriteControl(websocket.PingMessage, []byte("send data"), deadline)
				fmt.Println("here")
				//fmt.Println(message)
				//if lastResponse.UnixMilli()-time.Now().UnixMilli() < -750 {
				//	err = c.WriteMessage(websocket.TextMessage, []byte("9"))
				//	if err != nil {
				//		log.Printf("control err: %v\n", err)
				//	}
				//	lastResponse = time.Now()
				//}
				continue
			}
			fmt.Printf("MY MES: %s\n", string(message))
			for idx, r := range message {
				if r == '{' {
					message = message[idx : len(message)-1]
					break
				}
			}
			//fmt.Println(string(message))
			var ans struct {
				Payload   []*Rate `json:"payload,omitempty"`
				Nonce     int     `json:"nonce,omitempty"`
				Timestamp int64   `json:"timestamp"`
			}
			err = json.Unmarshal(message, &ans)
			if err != nil {
				log.Printf("gettickerchan: can't unmarshal message %s: %v\n", message, err)
				continue
			}
			update <- ans.Payload[0]
		}
		return
	}()
	return stopC, doneC, nil
}

/*func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				_ = c.Close()
				return
			}
		}
	}()
}*/

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	/*c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})*/

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			err := c.WriteMessage(websocket.TextMessage, []byte("\n"))
			if err != nil {
				return
			}
		}
	}()
}
