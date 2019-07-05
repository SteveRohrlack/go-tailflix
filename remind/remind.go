package remind

import (
	"time"
)

/*
After sends a signal to the signal channel when the given timeout (in seconds) is reached.
The function than terminates.
*/
func After(seconds int, signal chan<- interface{}) {
	var elapsed time.Duration
	started := time.Now()

	duration := time.Duration(seconds) * time.Second

	for {
		elapsed = time.Now().Sub(started)

		if elapsed > duration {
			signal <- nil
			break
		}

		time.Sleep(1 * time.Second)
	}
}
