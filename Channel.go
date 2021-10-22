package Coral

var channels = map[string]*Channel{}

type Channel struct {
	name        string
	connections map[string]*Conn
}

func newChannel(name string) *Channel {
	return &Channel{name: name, connections: map[string]*Conn{}}
}
func (c *Channel) AddMember(id string, conn *Conn) {
	c.connections[id] = conn
	conn.channels[c.name] = c
}

func (c *Channel) RemoveMember(id string) {
	delete(c.connections[id].channels, c.name)
	delete(c.connections, id)
}

func (c *Channel) GetMembers() map[string]*Conn {
	return c.connections

}
func (c *Channel) Broadcast(t string, content interface{}) {
	for _, conn := range c.connections {
		go conn.Emit(t, content)
	}
}
