package main

import (
	"fmt"
	"log"
)

func main() {
	host := "localhost:3091"
	_, err := StartServer(host)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started:", host)
}
