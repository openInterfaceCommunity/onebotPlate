package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Configuration
// Replace these with your actual Bot Access Token
const (
	AccessToken = "LYldoYM09S5rXtjibBxHUpcR_qv-cjHvbAFuZjDxaCw=" // Replace this!
	Host        = "incitymega.cn"
	Path        = "/bot/v1/onebot/v12/ws"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme:   "wss",
		Host:     Host,
		Path:     Path,
		RawQuery: "access_token=" + AccessToken,
	}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Wait for connection to be ready
	time.Sleep(1 * time.Second)

	// Fetch pending friend requests
	log.Println("Checking for pending friend requests...")
	reqObj := map[string]interface{}{
		"action": "get_friend_requests",
		"echo":   "init_friend_requests",
	}
	reqBytes, _ := json.Marshal(reqObj)
	c.WriteMessage(websocket.TextMessage, reqBytes)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handleMessage(c, message)
		}
	}()

	ticker := time.NewTicker(30 * time.Second) // Keep-alive if needed, though not strictly required by protocol
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Optional: implementation might require heartbeats, but OneBot 12 usually handles this via meta events
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// handleMessage parses and processes incoming JSON events
func handleMessage(c *websocket.Conn, msg []byte) {
	log.Printf("[DEBUG] Raw message received: %s", string(msg))
	var event map[string]interface{}
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Println("json unmarshal:", err)
		return
	}

	eventType, _ := event["type"].(string)

	if eventType == "message" {
		log.Printf("Received message: %s", string(msg))

		data, _ := event["data"].(map[string]interface{})
		detailType, _ := event["detail_type"].(string)

		// Construct reply
		action := map[string]interface{}{
			"action": "send_message",
			"params": map[string]interface{}{
				"detail_type": detailType,
				"message": []map[string]interface{}{
					{
						"type": "text",
						"data": map[string]string{
							"text": "Go Echo: ",
						},
					},
					// We blindly appending the original message content segments if possible,
					// or just simplify and send text. For complex structures, you'd iterate `data["content"]`.
				},
			},
			"echo": "go_echo_reply",
		}

		// Append original content to reply
		content, ok := data["content"].([]interface{})
		if ok {
			params := action["params"].(map[string]interface{})
			msgList := params["message"].([]map[string]interface{})

			// Convert []interface{} back to []map[string]interface{} is tricky in Go without struct,
			// so we just re-marshal or handle simple copy.
			// Ideally use proper structs (like in internal/onebot/model) but we want this example standalone.
			// Convert []interface{} back to []map[string]interface{}
			var newMsgList []map[string]interface{}
			for _, m := range msgList {
				newMsgList = append(newMsgList, m)
			}

			for _, seg := range content {
				if segMap, ok := seg.(map[string]interface{}); ok {
					newMsgList = append(newMsgList, segMap)
				}
			}
			params["message"] = newMsgList
		}

		// Set target ID
		params := action["params"].(map[string]interface{})
		if detailType == "private" {
			params["user_id"] = data["user_id"]
		} else if detailType == "group" {
			params["group_id"] = data["group_id"]
		}

		// Send reply
		replyBytes, _ := json.Marshal(action)
		log.Printf("Sending reply: %s", string(replyBytes))
		c.WriteMessage(websocket.TextMessage, replyBytes)

	} else if eventType == "meta" {
		detailType, _ := event["detail_type"].(string)
		if detailType == "heartbeat" {
			// Ignore heartbeat logs to keep console clean
		} else {
			log.Printf("Meta event: %s", detailType)
		}
	} else if eventType == "request" {
		detailType, _ := event["detail_type"].(string)
		if detailType == "friend" {
			data, _ := event["data"].(map[string]interface{})
			flag, _ := data["flag"].(string)
			userID, _ := data["user_id"].(string)

			log.Printf("Received friend request from %s, auto-accepting...", userID)

			action := map[string]interface{}{
				"action": "set_friend_add_request",
				"params": map[string]interface{}{
					"flag":    flag,
					"approve": true,
				},
				"echo": "auto_approve_friend",
			}

			actionBytes, _ := json.Marshal(action)
			c.WriteMessage(websocket.TextMessage, actionBytes)
		} else if detailType == "group" {
			// Handle group invitation
			data, _ := event["data"].(map[string]interface{})
			flag, _ := data["flag"].(string)
			groupID, _ := data["group_id"].(string)
			groupName, _ := data["group_name"].(string)

			log.Printf("Received group invitation to %s (%s), auto-accepting...", groupName, groupID)

			action := map[string]interface{}{
				"action": "set_group_add_request",
				"params": map[string]interface{}{
					"flag":    flag,
					"approve": true,
				},
				"echo": "auto_approve_group",
			}

			actionBytes, _ := json.Marshal(action)
			c.WriteMessage(websocket.TextMessage, actionBytes)
		}
	} else if eventType == "" {
		// Response
		echo, _ := event["echo"].(string)
		if echo == "init_friend_requests" {
			data, _ := event["data"].([]interface{})
			log.Printf("Found %d pending friend requests", len(data))
			for _, item := range data {
				req := item.(map[string]interface{})
				flag, _ := req["flag"].(string)
				userID, _ := req["user_id"].(string)

				log.Printf("Accepting pending request from %s...", userID)

				action := map[string]interface{}{
					"action": "set_friend_add_request",
					"params": map[string]interface{}{
						"flag":    flag,
						"approve": true,
					},
					"echo": "auto_approve_friend_init",
				}
				actionBytes, _ := json.Marshal(action)
				c.WriteMessage(websocket.TextMessage, actionBytes)
			}
		}
	}
}
