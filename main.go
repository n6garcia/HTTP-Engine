package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	method string
	route  string

	headers map[string]string
}

type Response struct {
	code int
	desc string

	headers map[string]string

	body []byte
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
	req := Request{headers: make(map[string]string)}

	i := 0
	scanner := bufio.NewScanner(con)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)

		if ln == "" {
			break
		}

		if i == 0 {
			fields := strings.Fields(ln)
			req.method = fields[0]
			req.route = fields[1]
		} else {
			strs := strings.SplitN(ln, ":", 2)
			if len(strs) != 2 {
				break
			}

			req.headers[strs[0]] = strs[1]

		}

		i++
	}

	/*
		for key, val := range req.headers {
			fmt.Println(key, " ", val)
		}*/

	return req
}

func respond(con net.Conn, req Request) {

	// if route is special "index.html"
	if req.route == "/" {
		body := `<b>Hello!</b>`
		fmt.Fprint(con, "HTTP/1.1 200 OK\r\n")
		fmt.Fprintf(con, "Content-Length: %d\r\n", len(body))
		fmt.Fprint(con, "Content-Type: text/html\r\n")
		fmt.Fprint(con, "\r\n")
		fmt.Fprint(con, body)
	} else {

		res := Response{headers: make(map[string]string)}

		res.code = 200
		res.desc = "OK"

		bs, err := os.ReadFile("./public_html/favicon.ico")
		if err != nil {
			panic(err)
		}

		res.body = bs
		res.headers["Content-Length"] = strconv.Itoa(len(bs))
		res.headers["Content-Type"] = "image/vnd.microsoft.icon"

		res.Write(con)

	}

}

func (res Response) Write(con net.Conn) {

	fmt.Fprintf(con, "HTTP/1.1 %d %s\r\n", res.code, res.desc)
	fmt.Printf("HTTP/1.1 %d %s\r\n", res.code, res.desc)
	for key, val := range res.headers {
		fmt.Fprintf(con, "%s: %s\r\n", key, val)
		fmt.Printf("%s: %s\r\n", key, val)
	}
	fmt.Fprintf(con, "\r\n")
	fmt.Printf("\r\n")
	con.Write(res.body)

}
