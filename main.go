package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func main() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Connected remote.")
		})

		http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var m message
			err := decoder.Decode(&m)
			if err != nil {
				fmt.Println("Error decoding message:", err)
			}

			if m.Text != "" {
				fmt.Println("Remote:", m.Text)
			}
		})

		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	for {
		text := readFromConsole("")
		switch text {
		case "exit":
			os.Exit(0)
		case "connect":
			text = readFromConsole("Enter ip address")
			address := text
			fmt.Println("Connecting to:", text, "...")
			body, err := connectTo(text)
			if err != nil {
				fmt.Println("Err:", err)
				return
			}
			io.Copy(os.Stdout, body)
			fmt.Println()
			fmt.Println()
			for {
				text = readFromConsole("You")
				if text == "--main" {
					break
				}
				msg := message{From: "local", Text: text}
				sendMsg(address, msg)
			}
		}
	}
}

func sendMsg(address string, msg message) {
	jsonMsg, _ := json.Marshal(msg)
	http.Post("http://"+address+":3000/msg", "application/json", bytes.NewBuffer(jsonMsg))
}

func connectTo(address string) (io.ReadCloser, error) {
	resp, err := http.Get("http://" + address + ":3000")
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func readFromConsole(inst string) string {
	reader := bufio.NewReader(os.Stdin)
	if inst == "" {
		inst = "Enter command"
	}
	fmt.Print(inst + ": ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, " \n")
	return text
}
