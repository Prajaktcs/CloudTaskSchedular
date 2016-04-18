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
	id   int
	task string
}

/*
	A blocking function to get the
*/
func process_channel(insert_channel chan *task, output_channel chan int) {
	input := <-insert_channel.task // get the input from the channel
	// for now, only supporting sleep tasks
	if strings.HasPrefix(input, "sleep") {
		value, _ := strconv.Atoi(strings.Split(input, " ")[1])
		fmt.Println(value)
		fmt.Println("Going to sleep")
		//setting the time in milliseconds
		time.Sleep(time.Millisecond)
		fmt.Println("Out of sleep")
	}

}

/*
	Add all the workload to the channel
*/
func add_local_load(wordload_file string, insert_channel chan string) int {
	file, _ := os.Open(wordload_file)
	reader := bufio.NewReaderSize(file, 4*1024)
	defer file.Close()
	fmt.Println("Reached here")
	line, isPrefix, err := reader.ReadLine()
	count := 0
	for err == nil && !isPrefix {
		fmt.Println("Inserting line" + string(line))
		insert_channel <- string(line) //Adding to the channel
		line, isPrefix, err = reader.ReadLine()
		count++
	}
	return count
}

/*
 Function to parse the command entered in the line
*/
func process_command(command string) bool {
	fmt.Println(command)
	if strings.HasPrefix(command, "client") {
		commands_slice := strings.Split(command, " ")
		if strings.Compare(commands_slice[0], "client") == 0 {
			if strings.Compare(commands_slice[1], "-s") == 0 || strings.Compare(commands_slice[3], "-s") == 0 {
				if strings.Compare(commands_slice[1], "-s") == 0 {
					if (strings.Compare(commands_slice[2], "LOCAL")) == 0 && (strings.Compare(commands_slice[3], "-t") == 0) && (strings.Compare(commands_slice[5], "-w") == 0) {
						//Pass the workload file
						a, _ := strconv.Atoi(commands_slice[4])
						insert_channel := make(chan string, a)
						output_channel := make(chan int)
						go add_local_load(commands_slice[6], insert_channel)

						go process_channel(insert_channel, output_channel)
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
		if process_command(string(command)) {
			fmt.Println("Nice progress")
		} else {
			break
		}

	}
}
