package vv104

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readCommandsFromStdIn(fromStdInchan chan string) {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		fromStdInchan <- sc.Text()
	}
}

func (state *State) evaluateInteractiveCommands() {
	fmt.Println("evaluateCommandsFromStdIn started")
	defer fmt.Println("evaluateCommandsFromStdIn returned")
	state.wg.Add(1)
	defer state.wg.Done()

	for {
		select {

		case input := <-state.chans.commandsFromStdin:
			inputSplit := strings.Split(input, " ")
			go state.evaluateInputSplit(inputSplit)

		case <-state.ctx.Done():
			fmt.Println("evaluateCommandsFromStdIn received Done(), returns")
			return
		}
	}
}

func (state *State) evaluateInputSplit(inputSplit []string) {
	var apdu Apdu
	switch len(inputSplit) {
	case 1:
		switch inputSplit[0] {
		case "restart":
			fmt.Println("called restart")
			state.cancel()

		case "exit":
			fmt.Println("Exiting")
			state.cancel()
			state.wg.Wait()
			os.Exit(1)

		case "startdt_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StartDTAct
			state.chans.toSend <- apdu

		case "startdt_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StartDTCon
			state.chans.toSend <- apdu

		case "stopdt_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StopDTAct
			state.chans.toSend <- apdu

		case "stopdt_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StopDTCon
			state.chans.toSend <- apdu

		case "testfr_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = TestFRAct
			state.chans.toSend <- apdu

		case "testfr_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = TestFRCon
			state.chans.toSend <- apdu
		}
	case 2:
	}

}
