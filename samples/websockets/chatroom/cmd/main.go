package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.wdf.sap.corp/practice-learning/websockets/chatroom"
)

func main() {
	// the trailing slash indicates the path is rooted subtrees which won't match any other subtree.
	http.HandleFunc("/home/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets/chatroom/home.html")
	})

	room := chatroom.NewRoom()
	go room.Start()

	http.HandleFunc("/chatroom", func(w http.ResponseWriter, r *http.Request) {
		chatroom.ServeWs(room, w, r)
	})
	addr := ":59988"
	logrus.Infof("server started at %s", addr)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}
