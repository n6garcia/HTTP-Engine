package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	method string
	route  string

	// Map string->string header->value
}

func main() {

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}

	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go handle(conn)
	}

}

func handle(con net.Conn) {
	defer con.Close()

	req := request(con)

	respond(con, req)

}

func request(con net.Conn) Request {
	req := Request{}

	i := 0
	scanner := bufio.NewScanner(con)
	for scanner.Scan() {
		ln := scanner.Text()

		if ln == "" {
			break
		}

		if i == 0 {
			fields := strings.Fields(ln)
			req.method = fields[0]
			req.route = fields[1]
		} else {
			strs := strings.Split(ln, ":")
			fmt.Println(strs)
		}

		i++
	}

	/* scan optional HTTP Body
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if ln == "" {
			break
		}
		i++
	}*/

	return req
}

func respond(con net.Conn, req Request) {
	body := `<b>Hello!</b>`

	fmt.Fprint(con, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(con, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(con, "Content-Type: text/html\r\n")
	fmt.Fprint(con, "\r\n")
	fmt.Fprint(con, body)
}
