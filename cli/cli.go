package cli

import (
	"bufio"
	"os"
)

/*
PromptEnterKey prompts the user to press the "enter" key and sends a signal to the continueChannel
*/
func PromptEnterKey(continueChannel chan<- interface{}, errorChannel chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	_, _, err := reader.ReadRune()

	if err != nil {
		errorChannel <- err
		return
	}

	continueChannel <- nil
}
