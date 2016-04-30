package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func processCommand(command string, outputQueueUrl string) bool {
	if strings.HasPrefix(command, "worker") {
		commandsSlice := strings.Split(command, " ")
		if strings.Compare(commandsSlice[0], "worker") == 0 {
			if strings.Compare(commandsSlice[1], "-s") == 0 {
				initialize(commandsSlice[2])
				for true {
					//insertChannel := make(chan string, 1000)
					counter := 0
					result, value := receiveMessage()
					if result == true {
						fmt.Println("value" + value)
						valueSlice := strings.Split(value, " ")
						//Check if the id is present in dynamodb
						id := valueSlice[0]
						//Picking the config from the config table

						if checkItem(commandsSlice[2], id) {
							fmt.Println("Already present, ignore")
						} else {
							writeItem(commandsSlice[2], id, valueSlice[2])
							//if it is a sleep task, sleep
							if strings.Compare(valueSlice[1], "sleep") == 0 {
								sleepTime, _ := strconv.Atoi(valueSlice[2])
								time.Sleep(time.Duration(sleepTime) * time.Millisecond)
								sendMessage(outputQueueUrl)
								//fmt.Println("Reached in if")
							} else if strings.Compare(valueSlice[1], "ffmpeg") == 0 {
								binary, lookErr := exec.LookPath("ls")
								if lookErr != nil {
									panic(lookErr)
								}
								fileName := strconv.Itoa(counter)
								args := []string{"wget", valueSlice[2], fileName + ".png"}
								env := os.Environ()
								syscall.Exec(binary, args, env)

							}
							if counter == 600 {
								break
							}
							//fmt.Println("Outside if")
						}
					}
					fmt.Println("Crossed this line")
				}
				binary, lookErr := exec.LookPath("ls")
				if lookErr != nil {
					panic(lookErr)
				}
				env := os.Environ()
				args1 := []string{"ffmpeg", "-framerate 1/5", "-i %03d.png", "-c:v", "libx264", "r 30", "worker1.mp4"}
				syscall.Exec(binary, args1, env)
				putItem("worker1.mp4", "prajaktvideos")
				sendMessage(outputQueueUrl)
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for true {
		command, _, _ := reader.ReadLine()
		outputQueueUrl := getItem("config", "1")
		if processCommand(string(command), outputQueueUrl) {
			fmt.Println("Nice progress")
		} else {
			break
		}

	}
}
