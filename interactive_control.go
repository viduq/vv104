package vv104

import (
	"bufio"
	"os"
)

func readCommandsFromStdIn(fromStdInchan chan string) {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		fromStdInchan <- sc.Text()
	}
}
