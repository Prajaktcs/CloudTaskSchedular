package main

import (
	"fmt"
	"strings"
)

/*

 */
func process_local_load(wordload_file string) {
	insert_channel := make(chan string)

}

/*
 Function to parse the command entered in the line
*/
func process_command(command string) bool {
	if strings.HasPrefix(command, "client") {
		commands_slice := strings.Split(command, " ")
		if strings.Compare(commands_slice[0], "client") == 0 {
			if strings.Compare(commands_slice[1], "-s") == 0 || strings.Compare(commands_slice[3], "-s") == 0 {
				if strings.Compare(commands_slice[1], "-s") == 0 {
					if strings.Compare(commands_slice[2], "LOCAL") == 0 && strings.Compare(commands_slice[3], "-w") == 0 {
						//Pass the workload file
						process_local_load(commands_slice[4])
					}
				}
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

func main() {
	/*
		run the client till user does not hit q
	*/
	var command string
	for true {
		fmt.Scanln(&command)
		if process_command(command) {
			fmt.Println("Nice progress")
		} else {
			break
		}

	}
}
