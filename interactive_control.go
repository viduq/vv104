package vv104

import (
	"bufio"
	"os"
	"strings"
)

func (state State) readCommandsFromStdIn() {
	logDebug.Println("readCommandsFromStdIn started ******")
	sc := bufio.NewScanner(os.Stdin)
	state.Wg.Add(1)
	defer state.Wg.Done()

	for sc.Scan() {
		select {
		default:
			state.Chans.CommandsFromStdin <- sc.Text()

		case <-state.Ctx.Done():
			logDebug.Println("readCommandsFromStdIn received Done(), returns")
			return

		}
	}
}

func (state *State) evaluateInteractiveCommands() {
	logDebug.Println("evaluateCommandsFromStdIn started")
	defer logDebug.Println("evaluateCommandsFromStdIn returned")
	state.Wg.Add(1)
	defer state.Wg.Done()

	for {
		select {

		case input := <-state.Chans.CommandsFromStdin:
			inputSplit := strings.Split(input, " ")
			state.evaluateInputSplit(inputSplit)

		case <-state.Ctx.Done():
			logDebug.Println("evaluateCommandsFromStdIn received Done(), returns")
			return
		}
	}
}

func (state *State) evaluateInputSplit(inputSplit []string) {
	var apdu Apdu
	switch inputArgsCount := len(inputSplit); {
	case inputArgsCount == 1:
		switch inputSplit[0] {
		case "restart":
			logInfo.Println("called restart")
			state.Cancel()

		case "exit":
			logInfo.Println("Exiting")
			state.Cancel()
			state.Wg.Wait()
			os.Exit(1)

		case "disconnect":
			logInfo.Println("Disconnecting")
			state.manualDisconnect = true
			state.Cancel()

		case "startdt_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StartDTAct
			state.Chans.ToSend <- apdu

		case "startdt_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StartDTCon
			state.Chans.ToSend <- apdu

		case "stopdt_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StopDTAct
			state.Chans.ToSend <- apdu

		case "stopdt_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = StopDTCon
			state.Chans.ToSend <- apdu

		case "testfr_act":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = TestFRAct
			state.Chans.ToSend <- apdu

		case "testfr_con":
			apdu = NewApdu()
			apdu.Apci.FrameFormat = UFormatFrame
			apdu.Apci.UFormat = TestFRCon
			state.Chans.ToSend <- apdu

		case "sp": // temporarily
			sp := NewApdu()
			infoObj := newInfoObj()
			infoObj.Ioa = 12345
			infoObj.Value = IntVal(1)
			sp.Apci.FrameFormat = IFormatFrame
			sp.Asdu.TypeId = M_SP_NA_1
			sp.Asdu.CauseTx = Spont
			sp.Asdu.Casdu = 1
			sp.Asdu.addInfoObject(infoObj)

			state.Chans.ToSend <- sp
		}

	case inputArgsCount > 2:

	}

}
