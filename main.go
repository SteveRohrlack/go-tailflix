package main

import (
	"fmt"
	"os"
	"tailflix/cli"
	"tailflix/file"
	"tailflix/remind"
)

func main() {

	// validate args

	if len(os.Args) != 2 {
		fmt.Println("please specify exactly one file")
		os.Exit(-1)
	}

	// get input file handle

	inputFile, err := file.Open(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer inputFile.Close()

	// show preview

	fmt.Printf("%v", file.Preview(inputFile, 3))

	// watch file

	fileContentChannel := make(chan byte, 128)
	errorChannel := make(chan error)
	remindChannel := make(chan interface{})
	continueChannel := make(chan interface{})
	timeoutChannel := make(chan interface{})

	defer close(errorChannel)
	defer close(fileContentChannel)
	defer close(remindChannel)
	defer close(continueChannel)
	defer close(timeoutChannel)

	const remindInterval = 30
	const remindTimeout = 30

	outputBuffer := make([]byte, 0)
	reminderActive := false

	// start watching the file
	go file.Watch(inputFile, fileContentChannel, errorChannel)

	// start the reminder
	go remind.After(remindInterval, remindChannel)

	for {
		select {

		case err := <-errorChannel:
			// print error and exit
			fmt.Println("whoops", err)
			os.Exit(-1)

		case byte := <-fileContentChannel:
			if reminderActive {
				// when a reminder is active, write to the outputBuffer
				outputBuffer = append(outputBuffer, byte)
				break
			}

			// print received byte
			fmt.Printf("%v", string(byte))

		case <-remindChannel:
			// set state
			reminderActive = true

			// print reminder
			fmt.Println("---------------------------")
			fmt.Println("|                         |")
			fmt.Println("| Are you still watching? |")
			fmt.Println("|                         |")
			fmt.Println("---------------------------")

			// wait for the user to press the "enter" key
			// use a go routine so the fileContentChannel can continue to receive data
			go cli.PromptEnterKey(continueChannel, errorChannel)

			// immediatelly start another go routine that writes to the timeoutChannel
			// if the user doesnt aknowledge the reminder, a timeout occurs
			go remind.After(remindTimeout, timeoutChannel)

		case <-continueChannel:
			// set state
			reminderActive = false

			// restart reminder
			go remind.After(remindInterval, remindChannel)

			println("Here's what you've missed:")

			// print the output buffer
			if len(outputBuffer) == 0 {
				fmt.Println("[nothing]")
			} else {
				fmt.Printf("%s\n", outputBuffer)
				fmt.Println("[continuing]")
			}

			outputBuffer = nil

		case <-timeoutChannel:
			if !reminderActive {
				// when there's no active reminder, this signal can be ignored
				break
			}

			fmt.Println("I'll close this since you stopped watching.")
			os.Exit(1)
		}

	}
}
