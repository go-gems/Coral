package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	Coral "github.com/go-gems/coral"
	"log"
)

/*
This example is designed to work with Gin but can work with any web server which gives you access to http.ResponseWriter
and http.Request.
 */
func main() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./", false)))

	r.GET("/ws", func(c *gin.Context) {
		//By default coralJS tries to access to /ws but everything is configurable
		// You have to enable the Coral.Handle if you want to handle websocket
		err := Coral.Handle(c.Writer, c.Request)
		if err != nil {
			panic(err)
		}
	})

	// on new join
	Coral.OnConnection(func(conn *Coral.Conn) error {

		// on new message receive
		conn.On("message", func(content interface{}) error {
			conn.Broadcast("forward", content)
			return nil
		})
		return nil
	})

	Coral.OnClose(func(conn *Coral.Conn) {
		log.Println("bye")
	})

	r.Run("localhost:12312")

}
