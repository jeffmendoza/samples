// Copyright 2016 Google, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type wsServer struct {
	da   *DataAccess
	cons *[](chan bool)
}

func (srv wsServer) removeSelf(ws *websocket.Conn, nd chan bool) {
	ws.Close()
	fmt.Println("Closed a connection")

	for ii, vv := range *(srv.cons) {
		if vv == nd {
			*srv.cons = append((*srv.cons)[:ii], (*srv.cons)[ii+1:]...)
			break
		}
	}

	fmt.Printf("Num of connections: %v\n", len(*srv.cons))
}

func (srv wsServer) handleOutgoing(ws *websocket.Conn, nd chan bool) {
	data := srv.da.Get()
	websocket.JSON.Send(ws, data)

	for {
		<-nd
		data = srv.da.Get()
		err := websocket.JSON.Send(ws, data)
		if err != nil {
			break
		}
	}

	srv.removeSelf(ws, nd)
}

func (srv wsServer) handleIncoming(ws *websocket.Conn, nd chan bool) {

	for {
		var data Data
		err := websocket.JSON.Receive(ws, &data)
		fmt.Printf("Received text: %v\n", data.Text)
		if err != nil {
			break
		}
	}

	srv.removeSelf(ws, nd)
}

func (srv wsServer) Serve(ws *websocket.Conn) {
	fmt.Println("New connection")
	// who is this??
	for header, val := range ws.Request().Header {
		fmt.Printf("%v %v\n", header, val)
	}

	//send users
	//if users change, send it
	//read changes to string??

	newCon := make(chan bool, 1)
	*(srv.cons) = append(*(srv.cons), newCon)

	fmt.Printf("Num of connections: %v\n", len(*srv.cons))

	srv.handleOutgoing(ws, newCon)
}

func main() {
	var cons [](chan bool)
	da := NewDataAccess()
	ws := wsServer{&da, &cons}
	go Simulate(&da, &cons)
	http.Handle("/ws", websocket.Handler(ws.Serve))

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static/"))))

	http.ListenAndServe(":8080", nil)
}
