package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all

	for{
		conn,err2 := ln.Accept()
		if err2 != nil{
			handleError(err2)
		}
		conns <- conn
	}
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	reader := bufio.NewReader(client)
	for{
		msg,_ := reader.ReadString('\n')
		msgs <- Message{clientid,msg}
	}
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()
	fmt.Println(*portPtr)
	//TODO Create a Listener for TCP connections on the port given above.
	ln,err1 := net.Listen("tcp", ":8030")
	if(err1 != nil){
		handleError(err1)
	}
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	i := 0
	//Start accepting connections
	go acceptConns(ln, conns)

	for {
		select {
		case conn := <-conns:
			clients[i] = conn
			go handleClient(conn,i,msgs)
			fmt.Fprint(conn,"connected\n")
			i++
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:
			fmt.Println(msg)
			sender := msg.sender
			for i := 0; i < len(clients); i++{
				if i == sender {
					continue
				}
				fmt.Fprint(clients[i], msg.message)
			}
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		}
	}
}
