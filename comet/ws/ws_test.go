package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

var reqBase = "127.0.0.1:8080"

func TestRoomMsg(t *testing.T) {
	require := require.New(t)

	// hxz
	sessionId1 := login(t, "ECEF70EDC523A11847BF96D6400ED98E", "139b893b345c93848a4e24ec3fb52a92f16463cd")
	log.Println(sessionId1)
	wsConn := wsLoop(sessionId1)
	go func() {
		for i := 0; i < 1000; i++ {
			msgBytes := wrapMsg(i)
			err := wsConn.WriteMessage(websocket.TextMessage, msgBytes)
			require.Nil(err)

			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func TestLoginAndOut(t *testing.T) {
	go func() {
		for i := 0; i < 1000; i++ {
			// wsj
			sessionId2 := login(t, "2342342234423423", "c5cf44aeda970f1266d8572b9e8255c3593096db")
			wsConn := wsLoop(sessionId2)
			time.Sleep(300 * time.Millisecond)
			err := wsConn.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	time.Sleep(30 * time.Second)
}

func wsLoop(sessionId string) *websocket.Conn {
	wsUrl := "ws://" + reqBase + "/ws"
	log.Println(wsUrl)

	var reqHeader = make(http.Header)
	reqHeader.Add("Cookie", sessionId)

	var err error
	wsConn, _, err := websocket.DefaultDialer.Dial(wsUrl, reqHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		for {
			// reset deadline for every type of message (control or data)
			if err := wsConn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				log.Println("failed to set read deadline", "err", err)
			}
			var msg []byte
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Server closed the connection")
				} else {
					log.Println("Failed to read from server", "err", err)
				}
				err := wsConn.Close()
				if err != nil {
					log.Println("Failed to close", "err", err)
				}
				return
			}
			log.Println("recv from server:", string(msg))
		}
	}()
	return wsConn
}

func login(t *testing.T, uuid string, token string) (sessionId string) {
	e := httpexpect.New(t, "http://"+reqBase)
	headers := map[string]string{
		"FZM-DEVICE":      "Android",
		"FZM-DEVICE-NAME": "Xiaomi Max 2",
		"FZM-UUID":        uuid,
		"FZM-AUTH-TOKEN":  token,
	}
	params := map[string]interface{}{
		"type": 2,
		"cid":  "",
	}
	resp := e.POST("/user/tokenLogin").WithJSON(params).WithHeaders(headers).Expect()
	cookie := resp.Status(http.StatusOK).Cookie("session-login")
	sessionId = cookie.Raw().Name + "=" + cookie.Raw().Value
	return
}

func wrapMsg(i int) []byte {
	msg := map[string]interface{}{
		"content": fmt.Sprintf("%v", i),
	}
	proto := map[string]interface{}{
		"avatar":      "https://zb-chat.oss-cn-shanghai.aliyuncs.com/chatList/picture/20181207/20181207105523088_13.jpg",
		"channelType": 2,
		"eventType":   0,
		"msg":         msg,
		"msgId":       "4b883aef-b977-41d9-bdd7-c8100940907e",
		"msgType":     1,
		"targetId":    "98",
		"user_level":  1,
		"isSnap":      2,
	}
	bytes, _ := json.Marshal(proto)
	return bytes
}
