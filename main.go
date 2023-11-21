package main

import (
	"fmt"
	"os"
	"tcc/client"
	"tcc/server"
)

func main() {
	args := os.Args[1:]

	if args[0] == "client"{
		fmt.Println("start client")
		if args[1]	== "r" {
			client.GetPerson(args[2])
		}else if args[1] == "w" {
			client.NewPerson(args[2], args[3], args[4])
		}
	} else if args[0] == "server" {
		fmt.Println("start server")
		server.StartServer()
	}

}