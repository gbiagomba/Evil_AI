package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <credentials file> <remote systems file>")
	}

	// Read the credentials file
	creds, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Parse the credentials into a username and password
	parts := strings.Split(string(creds), ":")
	if len(parts) != 2 {
		log.Fatal("Credentials file should contain a single line with a username and password separated by a colon")
	}
	username := parts[0]
	password := parts[1]

	// Read the remote systems file
	remotes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// Split the file into a slice of remote system addresses
	addresses := strings.Split(string(remotes), "\n")

	// Loop over the addresses and connect to each one
	for _, address := range addresses {
		if address == "" {
			continue
		}

		// Create a new SSH client
		client, err := ssh.Dial("tcp", address+":22", &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			log.Println(err)
			continue
		}
		defer client.Close()

		// Open a new session
		session, err := client.NewSession()
		if err != nil {
			log.Println(err)
			continue
		}
		defer session.Close()

		// Create a file to store the output for this remote system
		outputFile, err := os.Create(strings.Replace(address, ":", "_", -1) + ".txt")
		if err != nil {
			log.Println(err)
			continue
		}
		defer outputFile.Close()

		// Redirect the session's stdout to the file
		session.Stdout = outputFile

		// Run the command to search for private ssh keys on the remote system
		if err := session.Run("find / -name 'id_rsa'"); err != nil {
			log.Println(err)
			continue
		}
	}
}
