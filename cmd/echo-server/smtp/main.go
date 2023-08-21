package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a new buffered reader and writer for the connection
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Greeting message
	writer.WriteString("220 Welcome to My SMTP Server\r\n")
	writer.Flush()

	var from, to string
	var data []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		line = strings.TrimRight(line, "\r\n")
		parts := strings.Fields(line)
		var command string
		if len(parts) == 0 {
			command = ""
		} else {
			command = strings.ToUpper(parts[0])
		}

		switch command {
		case "HELO":
			writer.WriteString("250 Hello " + parts[1] + "\r\n")
			writer.Flush()
		case "MAIL":
			from = parts[2]
			writer.WriteString("250 Sender " + from + " OK\r\n")
			writer.Flush()
		case "RCPT":
			to = parts[2]
			writer.WriteString("250 Recipient " + to + " OK\r\n")
			writer.Flush()
		case "DATA":
			writer.WriteString("354 Enter message, ending with \".\" on a line by itself\r\n")
			writer.Flush()

			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading:", err)
					return
				}

				line = strings.TrimRight(line, "\r\n")
				if line == "." {
					break
				}
				data = append(data, line)
			}

			fmt.Println("Received message:")
			fmt.Println("From:", from)
			fmt.Println("To:", to)
			fmt.Println("Data:", strings.Join(data, "\n"))
			writer.WriteString("250 Message accepted for delivery\r\n")
			writer.Flush()
		case "QUIT":
			writer.WriteString("221 Closing connection. Goodbye!\r\n")
			writer.Flush()
			return
		default:
			writer.WriteString("500 Command not recognized\r\n")
			writer.Flush()
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:2525")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("SMTP server listening on localhost:25")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
