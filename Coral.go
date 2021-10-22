package Coral

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

const globalChannelName = "_coral.global"

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	onConnectionOpenAction  func(conn *Conn) error
	onConnectionCloseAction func(conn *Conn)
	onJoinAction            func(channel string, conn *Conn)
	onLeaveAction           func(channel string, conn *Conn)
	globalChannel           = newChannel(globalChannelName)
)

func Handle(w http.ResponseWriter, r *http.Request) error {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	conn := newConn(c)
	conn.On("_coral.channel.join", func(content interface{}) error {
		res := content.(map[string]string)
		if _, ok := channels[res["channel"]]; !ok {
			channels[res["channel"]] = newChannel(res["channel"])
		}
		channels[res["channel"]].AddMember(conn.ID.String(), conn)
		conn.Emit("_coral.channel.joined", res["channel"])
		onJoinAction(res["channel"], conn)
		return nil
	})
	conn.On("_coral.channel.leave", func(content interface{}) error {
		res := content.(map[string]string)

		channels[res["channel"]].RemoveMember(conn.ID.String())
		if len(channels[res["channel"]].GetMembers()) == 0 {
			delete(channels, res["channel"])
		}
		conn.Emit("_coral.channel.left", res["channel"])
		onLeaveAction(res["channel"], conn)
		return nil
	})
	defer conn.Close()

	if err := onConnectionOpenAction(conn); err != nil {
		return err
	}
	globalChannel.AddMember(conn.ID.String(), conn)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		var message message = message{ID: conn.ID.String()}
		if err := json.Unmarshal(msg, &message); err != nil {
			return err
		}

		if act, ok := conn.actions[message.Type]; ok {
			if err = act(message.Content); err != nil {
				return err
			}
		}
	}
}

func OnConnection(callback func(conn *Conn) error) {
	onConnectionOpenAction = callback
}
func OnClose(callback func(conn *Conn)) {
	onConnectionCloseAction = callback
}

func onChannelJoin(callback func(channel string, conn *Conn)) {
	onJoinAction = callback
}

func Broadcast(t string, content interface{}) {
	for _, conn := range globalChannel.GetMembers() {
		go conn.Emit(t, content)
	}
}
