package ws

import (
	"encoding/json"
	"errors"
	"github.com/Masterminds/semver"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/openaccounting/oa-server/core/auth"
	"github.com/openaccounting/oa-server/core/model/types"
	"log"
	"net/http"
	"sync"
)

const version = "1.0.0"

//var upgrader = websocket.Upgrader{} // use default options
var txSubscriptions = make(map[string][]*websocket.Conn)
var accountSubscriptions = make(map[string][]*websocket.Conn)
var priceSubscriptions = make(map[string][]*websocket.Conn)
var userMap = make(map[*websocket.Conn]*types.User)
var sequenceNumbers = make(map[*websocket.Conn]int)
var locks = make(map[*websocket.Conn]*sync.Mutex)

type Message struct {
	Version        string      `json:"version"`
	SequenceNumber int         `json:"sequenceNumber"`
	Type           string      `json:"type"`
	Action         string      `json:"action"`
	Data           interface{} `json:"data"`
}

func Handler(w rest.ResponseWriter, r *rest.Request) {
	c, err := websocket.Upgrade(w.(http.ResponseWriter), r.Request, nil, 0, 0)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	sequenceNumbers[c] = -1
	locks[c] = &sync.Mutex{}

	defer c.Close()
	for {
		mt, messageData, err := c.ReadMessage()
		if err != nil {
			log.Println("readerr:", err)
			// remove connection from maps
			unsubscribeAll(c)
			break
		}

		if mt != websocket.TextMessage {
			log.Println("Unsupported message type", mt)
			continue
		}

		message := Message{}

		err = json.Unmarshal(messageData, &message)

		if err != nil {
			log.Println("Could not parse message:", string(messageData))
			continue
		}

		log.Printf("recv: %s", message)

		// check version
		err = checkVersion(message.Version)

		if err != nil {
			log.Println(err.Error())
			writeMessage(c, websocket.CloseMessage, websocket.FormatCloseMessage(4001, err.Error()))
			break
		}

		// make sure they are authenticated
		if userMap[c] == nil {
			log.Println("checking message for authentication")
			err = authenticate(message, c)
			if err != nil {
				log.Println("Authentication error " + err.Error())
				writeMessage(c, websocket.CloseMessage, websocket.FormatCloseMessage(4000, err.Error()))
				break
			}
			continue
		}

		err = processMessage(message, c)

		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func getKey(userId string, orgId string) string {
	return userId + "-" + orgId
}

func processMessage(message Message, conn *websocket.Conn) error {
	var dataString string
	err := mapstructure.Decode(message.Data, &dataString)

	if err != nil {
		return err
	}

	key := getKey(userMap[conn].Id, dataString)

	log.Println(message.Action, message.Type, dataString)

	switch message.Action {
	case "subscribe":
		switch message.Type {
		case "transaction":
			subscribe(conn, key, txSubscriptions)
		case "account":
			subscribe(conn, key, accountSubscriptions)
		case "price":
			subscribe(conn, key, priceSubscriptions)
		default:
			return errors.New("Unhandled message type: " + message.Type)
		}
	case "unsubscribe":
		switch message.Type {
		case "transaction":
			unsubscribe(conn, key, txSubscriptions)
		case "account":
			unsubscribe(conn, key, accountSubscriptions)
		case "price":
			unsubscribe(conn, key, priceSubscriptions)
		default:
			return errors.New("Unhandled message type: " + message.Type)
		}
	case "ping":
		sequenceNumbers[conn]++
		response := Message{version, sequenceNumbers[conn], "pong", "pong", nil}
		responseData, err := json.Marshal(response)

		if err != nil {
			return err
		}

		err = writeMessage(conn, websocket.TextMessage, responseData)

		if err != nil {
			unsubscribeAll(conn)
			return err
		}
	}

	return nil
}

func subscribe(conn *websocket.Conn, key string, clientMap map[string][]*websocket.Conn) {
	conns := clientMap[key]
	alreadySubscribed := false

	for _, c := range conns {
		if conn == c {
			alreadySubscribed = true
		}
	}

	if alreadySubscribed == false {
		clientMap[key] = append(clientMap[key], conn)
	}
}

func unsubscribe(conn *websocket.Conn, key string, clientMap map[string][]*websocket.Conn) {
	newConns := clientMap[key][:0]

	for _, c := range clientMap[key] {
		if conn != c {
			newConns = append(newConns, c)
		}
	}
}

func unsubscribeAll(conn *websocket.Conn) {
	for key, conns := range txSubscriptions {
		newConns := conns[:0]
		for _, c := range conns {
			if conn != c {
				newConns = append(newConns, c)
			}
		}
		txSubscriptions[key] = newConns
	}

	for key, conns := range accountSubscriptions {
		newConns := conns[:0]
		for _, c := range conns {
			if conn != c {
				newConns = append(newConns, c)
			}
		}
		accountSubscriptions[key] = newConns
	}

	for key, conns := range priceSubscriptions {
		newConns := conns[:0]
		for _, c := range conns {
			if conn != c {
				newConns = append(newConns, c)
			}
		}
		priceSubscriptions[key] = newConns
	}

	delete(userMap, conn)
	delete(sequenceNumbers, conn)
	delete(locks, conn)
}

func PushTransaction(transaction *types.Transaction, userIds []string, action string) {
	log.Println(txSubscriptions)

	message := Message{version, -1, "transaction", action, transaction}

	for _, userId := range userIds {
		key := getKey(userId, transaction.OrgId)
		for _, conn := range txSubscriptions[key] {
			sequenceNumbers[conn]++
			message.SequenceNumber = sequenceNumbers[conn]
			messageData, err := json.Marshal(message)

			if err != nil {
				log.Println("PushTransaction json error:", err)
				return
			}

			err = writeMessage(conn, websocket.TextMessage, messageData)

			if err != nil {
				log.Println("Cannot PushTransaction to client:", err)
				unsubscribeAll(conn)
			}
		}
	}
}

func PushAccount(account *types.Account, userIds []string, action string) {
	message := Message{version, -1, "account", action, account}

	for _, userId := range userIds {
		key := getKey(userId, account.OrgId)
		for _, conn := range accountSubscriptions[key] {
			sequenceNumbers[conn]++
			message.SequenceNumber = sequenceNumbers[conn]
			messageData, err := json.Marshal(message)

			if err != nil {
				log.Println("PushAccount error:", err)
				return
			}
			err = writeMessage(conn, websocket.TextMessage, messageData)

			if err != nil {
				log.Println("Cannot PushAccount to client:", err)
				unsubscribeAll(conn)
			}
		}
	}
}

func PushPrice(price *types.Price, userIds []string, action string) {
	message := Message{version, -1, "price", action, price}

	for _, userId := range userIds {
		key := getKey(userId, price.OrgId)
		for _, conn := range priceSubscriptions[key] {
			sequenceNumbers[conn]++
			message.SequenceNumber = sequenceNumbers[conn]
			messageData, err := json.Marshal(message)

			if err != nil {
				log.Println("PushPrice error:", err)
				return
			}

			err = writeMessage(conn, websocket.TextMessage, messageData)

			if err != nil {
				log.Println("Cannot PushPrice to client:", err)
				unsubscribeAll(conn)
			}
		}
	}
}

func authenticate(message Message, conn *websocket.Conn) error {
	var id string
	err := mapstructure.Decode(message.Data, &id)

	if err != nil {
		return err
	}

	if message.Action != "authenticate" {
		return errors.New("Authentication required")
	}

	user, err := auth.Instance.Authenticate(id, "")
	if err != nil {
		return err
	}

	userMap[conn] = user

	return nil
}

func checkVersion(clientVersion string) error {
	constraint, err := semver.NewConstraint(clientVersion)

	if err != nil {
		return errors.New("Invalid version")
	}

	serverVersion, _ := semver.NewVersion(version)

	versionMatch := constraint.Check(serverVersion)

	if versionMatch != true {
		return errors.New("Invalid version")
	}

	return nil
}

func writeMessage(conn *websocket.Conn, messageType int, data []byte) error {
	locks[conn].Lock()
	defer locks[conn].Unlock()
	return conn.WriteMessage(messageType, data)
}
