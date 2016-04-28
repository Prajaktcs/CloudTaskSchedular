package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type task struct {
	id       int
	taskType int
	task     string
}

/*
	Method to process the local queue. Will start threads
	till the channel doesnt block
*/
func execute(taskChannel chan *task, signalChannel chan int, numberOfTasks int, outputChannel chan int) {
	for i := 0; i < numberOfTasks; i++ {
		taskInstance := <-taskChannel
		// Sleep task is 1
		if taskInstance.taskType == 1 {
			task1 := strings.Split(taskInstance.task, " ")
			multiplier, _ := strconv.Atoi(task1[1])
			go sleep(multiplier, outputChannel, taskInstance.id)
		}

	}

	signalChannel <- 1
}

func sleep(multiplier int, outputChannel chan int, taskId int) {
	time.Sleep(time.Duration(multiplier) * time.Millisecond)
	outputChannel <- taskId
}

/*
	A blocking function to schedule the number of tasks
*/
func processChannel(insertChannel chan string, outputChannel chan int, numberOfTasks int) {
	//Running the code till all the channel is processed
	signalChannel := make(chan int, 1)
	signalChannel <- 1
	tasksQueue := make(chan *task, numberOfTasks)
	counter := 0
	for true {
		// block till all the values are executed
		flag := <-signalChannel
		if flag == 1 {
			for i := 0; i < numberOfTasks; i++ {
				input := <-insertChannel // get the input from the channel
				counter += 1
				taskInstance := task{counter, 1, input}
				tasksQueue <- &taskInstance

			}
			execute(tasksQueue, signalChannel, numberOfTasks, outputChannel)
		}
	}

}

/*
	Add all the workload to the channel
*/
func addLocalLoad(workloadFile string, insertChannel chan string) int {
	file, _ := os.Open(workloadFile)
	reader := bufio.NewReaderSize(file, 4*1024)
	defer file.Close()
	fmt.Println("Reached here")
	line, isPrefix, err := reader.ReadLine()
	count := 0
	for err == nil && !isPrefix {
		insertChannel <- string(line) //Adding to the channel
		line, isPrefix, err = reader.ReadLine()
		count++
	}
	return count
}

/**
function to read file, and initialize queue and dynamodb
*/
func processRemoteQueue(queueName string, workloadFile string) {
	tasksList := make([]string, 0)
	file, _ := os.Open(workloadFile)
	reader := bufio.NewReaderSize(file, 4*1024)
	defer file.Close()
	fmt.Println("Reached here")
	line, isPrefix, err := reader.ReadLine()
	count := 0
	for err == nil && !isPrefix {
		tasksList = append(tasksList, string(line))
		line, isPrefix, err = reader.ReadLine()
		count++
	}
	//Create the queue
	createQueue(queueName)
	url := createQueue("output")
	//Here create dynamodb table
	createTable(queueName)
	createTable("config")
	writeItem("config", "1", url)

	start := time.Now()
	count = 0
	for index, element := range tasksList {
		key := strconv.Itoa(index)
		sendMessage(key + " " + element)
		count = index
	}
	//now waiting till we get all the messages on the output queue
	i := 0
	for i <= count {
		isSucccessful, _ := receiveMessage(url)
		if isSucccessful == true {
			i++
			fmt.Println("Received output")
		}
	}
	end := time.Since(start).String()
	fmt.Println(end)
}

/*
 Function to parse the command entered in the line
*/
func processCommand(command string) bool {
	fmt.Println(command)
	if strings.HasPrefix(command, "client") {
		commandsSlice := strings.Split(command, " ")
		if strings.Compare(commandsSlice[0], "client") == 0 {
			if strings.Compare(commandsSlice[1], "-s") == 0 || strings.Compare(commandsSlice[3], "-s") == 0 {
				if strings.Compare(commandsSlice[1], "-s") == 0 {
					if (strings.Compare(commandsSlice[2], "LOCAL")) == 0 && (strings.Compare(commandsSlice[3], "-t") == 0) && (strings.Compare(commandsSlice[5], "-w") == 0) {
						//Pass the workload file
						a, _ := strconv.Atoi(commandsSlice[4])
						insertChannel := make(chan string, 1000)
						outputChannel := make(chan int, 1000)
						count := addLocalLoad(commandsSlice[6], insertChannel)
						start := time.Now()
						go processChannel(insertChannel, outputChannel, a)
						fmt.Println(count)
						for i := 0; i < count; i++ {
							//THis will block till we dont get all the ids
							select {
							case <-outputChannel:
								fmt.Print()
							case <-time.After(time.Second * 3):
								fmt.Println("timeout 2")
							}
						}
						end := time.Since(start).String()
						fmt.Println(end)
					} else {
						processRemoteQueue(commandsSlice[2], commandsSlice[4])
					}
				}
			}
		} else {
			return false
		}

	} else {
		return false
	}
	return true
}

func main() {
	/*
		run the client till user does not hit q
	*/
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
