package main

import (
	"net"
	"net/http"
)

func main() {

	l, _ := net.Listen("tcp", ":8080")

	http.Serve(l, nil)

}
