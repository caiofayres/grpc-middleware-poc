package main

import (
	"fmt"
	"os"
	"tcc/server"
	"tcc/client"
)

func main() {
	//flag for start client or server
	args := os.Args[1:]

	if args[0] == "client"{
		fmt.Println("start client")
		client.SendMessage(args[1])
	} else if args[0] == "server" {
		fmt.Println("start server")
		server.StartServer()
	}

}