package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

var done = make(chan bool)

func main() {
	input := io.ReadCloser(os.Stdin)
	stdin := readIntoChannel(input, make(chan bool))
	for {
		select {
		case x := <-stdin:
			message := parseMessage(x)
			switch message.Input.Type {
			case analyze:
				go doAnalysis(message)
			case stop:
				done <- true
			}
		case <-done:
			return
		}
	}
}

func doAnalysis(msg message) {
	target := msg.Input.Data
	file, _ := os.OpenFile(target, os.O_RDONLY, 600)
	time.Sleep(1000 * time.Millisecond)
	fileData := make([]byte, 16)
	_, _ = file.Read(fileData)
	pairs := make([]pair, 1)
	pairs[0] = pair{Key: "example", Value: string(fileData)}
	//pairs[1] = pair{Key: "example2", Value: "flurgldy-pop"}
	print, err := json.Marshal(message{Input: msg.Input, Output: output{Pairs: pairs}})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(print))
	done <- true
}

func readIntoChannel(rc io.ReadCloser, done <-chan bool) chan string {
	out := make(chan string)
	go func() {
		reader := bufio.NewScanner(rc)
		for {
			if !reader.Scan(){
				//fmt.Println("done")
				return
			}
			//fmt.Println("line read...")
			select {
			case out <- reader.Text():
			case <-done:
				//fmt.Println("done")
				return
			}
		}
	}()
	return out
}

func parseMessage(in string) message {
	var target message
	json.Unmarshal([]byte(in), &target)
	return target
}

type message struct {
	Input 	input
	Output 	output
}

type input struct {
	Type	messageType
	Data	string
}

type output struct {
	Pairs	[]pair
}

type pair struct {
	Key	string
	Value	string
}

type messageType int

const (
	analyze messageType = iota
	suspend
	resume
	stop
)
