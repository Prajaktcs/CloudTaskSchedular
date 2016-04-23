package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func processCommand(command string) bool {
	fmt.Println(command)
	if strings.HasPrefix(command, "worker") {
		commandsSlice := strings.Split(command, " ")
		if strings.Compare(commandsSlice[0], "worker") == 0 {
			if strings.Compare(commandsSlice[1], "-s") == 0 {
				initialize(commandsSlice[2])
				for true {
					//insertChannel := make(chan string, 1000)
					result, value := receiveMessage()
					if result == true {
						fmt.Println(value)
					}
					fmt.Println("Crossed this line")
				}
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for true {
		command, _, _ := reader.ReadLine()
		if processCommand(string(command)) {
			fmt.Println("Nice progress")
		} else {
			break
		}

	}
}
