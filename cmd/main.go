// Package for go-ssh-app testing an idea.
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:2222")

	if err != nil {
		log.Fatal("failed to create listener")
	}
	defer listener.Close()

	log.Print("Server started...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Failed to accept connection...", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	fileBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("could not read id_rsa")
	}

	key, err := ssh.ParsePrivateKey(fileBytes)
	if err != nil {
		log.Fatal("could not parse id_rsa")
	}

	config.AddHostKey(key)

	_, chans, reqs, err := ssh.NewServerConn(c, config)
	if err != nil {
		log.Fatal("failed to initialize server connection: ", err)
	}

	go ssh.DiscardRequests(reqs)

	go func() {
		for newChannel := range chans {
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			connection, _, err := newChannel.Accept()
			if err != nil {
				log.Print("Could not accept channel: ", err)
				continue
			}

			go sshProgram(&connection)
		}
	}()

	go func() {
		for req := range reqs {
			log.Println(req.Type, req.Payload)
			req.Reply(true, nil)
		}
	}()
}

func sshProgram(c *ssh.Channel) {
	log.Printf("New session started for %s", "session")
	for {
		(*c).Write([]byte("Input:"))

		for {
			readBuffer := make([]byte, 0, 1024)
			(*c).Read(readBuffer)
			endOfInput := false

			for _, v := range readBuffer {
				if v != '\n' {
					fmt.Print(v)
					continue
				}

				endOfInput = true
				break
			}

			if endOfInput {
				break
			}
		}
	}
}
