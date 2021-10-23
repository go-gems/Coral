# Coral

Coral is a small library which aims to have the same workflow as SocketIO.

## Requirements
- golang 1.16+

## Getting Started

Install the library in your go project : 

```
go get -u github.com/go-gems/coral
```

Then in your projects you have to implement a web server and on a route call the Coral.handle method.
For example on a vanilla http server : 
```go
package main

import (
	"github.com/go-gems/coral"
	"net/http"
)

func main() {

	// handle coral requests on /ws.
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		if err := Coral.Handle(writer, request); err != nil {
			panic(err)
		}
	})
	
	// serve a public directory on /
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.ListenAndServe(":8080", nil)
}
```
the example folder contains the implement with gin.

You can now listen on every websocket connection:
```go
//event to do when user joins the server.
Coral.OnConnection(func(conn *Coral.Conn) error {
    log.Println("connection established")
    
    //Broadcast to every user.
    Coral.Broadcast("user:joined","Somebody joined!")
    
    //conn Object represents a single user connected.
    //When conn sends a "says" type (or anything you want type)
    conn.On("says", func(content interface{}) error {
        //broadcast to all except sender
         conn.Broadcast("listen", content)
         return nil
    })
    // When conn sends a "ping" message.
    conn.On("ping", func(content interface{}) error {
        //Forward to all except conn.
        return conn.Emit("pong", content)
    })
    return nil
})
```

Javascript side, you can use the Coral.js which is the beginning of a library - and thus will be published as a npm package.
Example usage : 
```HTML
<!doctype html>
<html>
    <body>
        <h1>Coral Demo.</h1>
        <script type="module" src="script.js"></script>
    </body>
</html>
```
and in script.js.
```javascript
import Coral from "./Coral.js";

//coral communication instance
//default connection to same host/ws
const coral = new Coral()
// if needed
// const coral = new Coral({url: "ws://test.com/something"})

//on connection open
coral.onOpen = function () {
    console.log("connection established");
    
    // on message recieved on listen
    coral.on("listen", (msg) => {
        console.log("received message ", msg)
    })
    
    coral.on("user:joined", (msg)=>{
        console.log(msg);
        coral.emit("says","Hello new User.");

    })
    coral.on("pong", (msg) => {
        console.log("PONG", msg)
    })

    coral.emit("ping","Hello");
}

```

And that's all for now!

## Contribution

If you like the project, feel free to contribute, there are plenty of things to do : 
- Transform coral.js as a npm library
- Add Channels management (WIP)

So every PR is welcome.