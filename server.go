package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/rsms/gotalk"
)

var (
	Clients = make(map[*gotalk.Sock]int)
	socksmu sync.RWMutex
)

func onAccept(s *gotalk.Sock) {
	// Keep track of connected sockets
	socksmu.Lock()
	defer socksmu.Unlock()
	Clients[s] = 1

	s.CloseHandler = func(s *gotalk.Sock) {
		socksmu.Lock()
		defer socksmu.Unlock()
		delete(Clients, s)
	}
}

func Broadcast(name string, in interface{}) {
	socksmu.RLock()
	defer socksmu.RUnlock()
	for s, _ := range Clients {
		s.Notify(name, in)
	}
}

func StartServer() {

	http.Handle("/gotalk", gotalk.WebSocketHandler(nil, onAccept))

	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{
				Asset:    Asset,
				AssetDir: AssetDir,
				Prefix:   "web",
			},
		),
	)

	var port string = ":8080"

	if len(*ServerPort) > 0 {
		port = fmt.Sprintf(":%s", ServerPort)
	}

	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
