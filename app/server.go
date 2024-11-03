package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const seperator = "\r\n"

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Unable to accept the connection", err.Error())
			continue
		}
		go handleConnection(conn)

	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)

	requests := parser(reader)

	if len(requests) < 1 {
		fmt.Println("No requests or Something Wrong with Request Format")
	}

	for i := 0; i < len(requests); i++ {
		subRequest := requests[i]

		currentCommand := subRequest[0]

		if strings.ToUpper(currentCommand) == "ECHO" {
			response := echoHandler(subRequest)
			if response == "" {
				fmt.Println("bad response")
			}
			fmt.Println("response", response)
			_, err := c.Write([]byte(response))

			if err != nil {
				fmt.Println("something went wrong while repsondiong", err.Error())
			}
		}
	}
}

func echoHandler(subrequest []string) string {
	if len(subrequest) != 2 {
		fmt.Println("Invalid Request Format")
		return ""
	}

	response := encoder(subrequest[1])
	return response
}

func encoder(response string) string {
	res := "$"
	res += strconv.Itoa(len(response))
	res += seperator
	res += response
	res += seperator
	fmt.Println(res)
	return res

}

func parser(buffer *bufio.Reader) [][]string {
	requests := [][]string{}

	// We are parsing the request into multiple Commands
	// Each command starts with *<int>, <int> represents number of parameters, and there could be multiple requests.

	var buf bytes.Buffer
	_, err := io.Copy(&buf, buffer)
	if err != nil {
		panic(err)
	}

	request := buf.String()
	println("Complete request:", request)

	//Lets get all the subrequests that are there in this main request
	subRequests := strings.Split(request, "*")

	println("subrequests", subRequests)

	for _, subreq := range subRequests {

		if subreq == "" {
			continue
		}

		// Each command is madeup of multiple lines, we are parsing the lines here
		linesOfCommand := strings.Split(subreq, "\r\n")

		println("subreq", subreq)
		println("linesOfCommand", linesOfCommand)

		if len(linesOfCommand) < 2 {
			continue
		}

		//First line in each command represents numberofparameters
		numberOfParamters, err := strconv.Atoi(linesOfCommand[0])

		if err != nil {
			fmt.Println("Invalid Command, unable to read number of parameters", err.Error())
			continue
		}

		commands := []string{}
		for i := 0; i < numberOfParamters; i++ {
			lineIndex := i*2 + 2
			if lineIndex >= len(linesOfCommand) {
				fmt.Println("Command Length Mismatch")
				continue
			}
			commands = append(commands, linesOfCommand[lineIndex])
		}

		requests = append(requests, commands)
	}

	return requests
}
