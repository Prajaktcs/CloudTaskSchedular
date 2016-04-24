package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func processCommand(command string, outputQueueUrl string) bool {
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
						valueSlice := strings.Split(value, " ")
						//Check if the id is present in dynamodb
						id := valueSlice[0]
						if checkItem(commandsSlice[2], id) {
							fmt.Println("Already present, ignore")
						} else {
							writeItem(commandsSlice[2], id, valueSlice[2])
							//if it is a sleep task, sleep
							if strings.Compare(commandsSlice[1], "sleep") == 0 {
								sleepTime, _ := strconv.Atoi(commandsSlice[2])
								time.Sleep(time.Duration(sleepTime) * time.Millisecond)
								sendMessage(outputQueueUrl)
							}
						}
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
	outputQueueUrl := flag.String("url", "", "Url for the output queue")
	for true {
		command, _, _ := reader.ReadLine()
		if processCommand(string(command), *outputQueueUrl) {
			fmt.Println("Nice progress")
		} else {
			break
		}

	}
}
