package vv104

import (
	"bufio"
	"os"
)

func (state *State) readCommandsFromStdIn() {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		state.chans.commandsFromStdin <- sc.Text()
	}
}
