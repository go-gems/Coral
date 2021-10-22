package Coral

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Conn struct {
	*websocket.Conn
	ID       uuid.UUID
	locker   sync.Mutex
	actions  map[string]func(content interface{}) error
	channels map[string]*Channel
}

func newConn(conn *websocket.Conn) *Conn {
	return &Conn{
		Conn:     conn,
		ID:       uuid.New(),
		actions:  map[string]func(content interface{}) error{},
		channels: map[string]*Channel{},
		locker:   sync.Mutex{},
	}
}

func (c *Conn) Close() {
	_ = c.Conn.Close()
	globalChannel.RemoveMember(c.ID.String())
	onConnectionCloseAction(c)

}

func (c *Conn) On(t string, callback func(content interface{}) error) {
	c.actions[t] = callback
}

func (c *Conn) Emit(t string, content interface{}) error {
	barray, err := json.Marshal(message{ID: c.ID.String(), Type: t, Content: content})
	if err != nil {
		return err
	}
	c.locker.Lock()
	defer c.locker.Unlock()
	if err := c.WriteMessage(websocket.TextMessage, barray); err != nil {
		return err
	}
	return nil
}
func (c *Conn) Broadcast(t string, content interface{}) {
	for id, conn := range globalChannel.GetMembers() {
		if c.ID.String() != id {
			go conn.Emit(t, content)
		}
	}
}

func (c *Conn) GetChannels() {

}

func (c *Conn) BroadcastToChannel(name string, t string, content interface{}) {
	if channel, ok := channels[name]; ok {
		for id, conn := range channel.connections {
			if c.ID.String() != id {
				go conn.Emit(t, content)
			}
		}
	}

}
