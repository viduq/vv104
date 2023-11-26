package vv104

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State struct {
	Config    Config
	ConnState ConnState
	ssn       SeqNumber
	rsn       SeqNumber
	chans     AllChans
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

type ConnState int

type AllChans struct {
	commandsFromStdin chan string
	received          chan Apdu
	toSend            chan Apdu
}

const (
	STOPPED ConnState = iota
	STARTED
	PENDING_UNCONFIRMED_STOPPED
	PENDING_STARTED
	PENDING_STOPPED
	STOP_CONN
	START_CONN
)

func NewState() State {
	return State{
		Config: Config{
			Mode:          "",
			Ipv4Addr:      "",
			Port:          0,
			Casdu:         0,
			AutoAck:       false,
			K:             0,
			W:             0,
			T1:            0,
			T2:            0,
			T3:            0,
			IoaStructured: false,
			UseLocalTime:  false,
		},
		ConnState: 0,
		ssn:       0,
		rsn:       0,
		chans:     AllChans{},
		// wg:                    sync.WaitGroup{},
	}
}

func (state *State) Start() {
	state.Config.ParseFlags()
	if state.Config.InteractiveMode {
		state.chans.commandsFromStdin = make(chan string, 30)
		go readCommandsFromStdIn(state.chans.commandsFromStdin) // never exits
		go state.evaluateCommandsFromStdIn()
	}

	for {
		state.chans.received = make(chan Apdu, state.Config.W)
		state.chans.toSend = make(chan Apdu, state.Config.K)
		state.ctx, state.cancel = context.WithCancel(context.Background())
		go state.startConnection()

		if state.Config.InteractiveMode {
			go state.evaluateCommandsFromStdIn()
		}

		<-state.ctx.Done()
		state.wg.Wait()
		fmt.Println("Restart!")
		time.Sleep(1500 * time.Millisecond)
	}
}

func (state *State) serverStateMachine() {
	var apduReceived Apdu
	var apduToSend Apdu

	state.wg.Add(1)
	defer state.wg.Done()

	state.ConnState = STOPPED
	fmt.Println("Entering state STOPPED")

	for {
		select {

		case apduReceived = <-state.chans.received:
			fmt.Println("<<RX:", apduReceived) // put after next if statement, if testfr shall not be printed

			if apduReceived.Apci.UFormat == TestFRAct {
				// always reply to testframes
				apduToSend = NewApdu()
				apduToSend.Apci.FrameFormat = UFormatFrame
				apduToSend.Apci.UFormat = TestFRCon
				state.chans.toSend <- apduToSend
				continue
			}

			// state machine
			switch state.ConnState {

			case STOPPED:
				if apduReceived.Apci.UFormat == StartDTAct {
					// startdt act received
					apduToSend = NewApdu()
					apduToSend.Apci.FrameFormat = UFormatFrame
					apduToSend.Apci.UFormat = StartDTCon
					state.chans.toSend <- apduToSend
					state.ConnState = STARTED
					fmt.Println("Entering state STARTED")

				}
			case STARTED:
				if apduReceived.Apci.UFormat == StopDTAct {
					// stopdt act received
					apduToSend.Apci.FrameFormat = UFormatFrame
					apduToSend.Apci.UFormat = StopDTCon
					state.chans.toSend <- apduToSend
					state.ConnState = STOPPED
					fmt.Println("Entering state STOPPED")

				}
			}

		case <-state.ctx.Done():
			fmt.Println("serverStateMachine received ctx.Done, returns")
			return
		}
	}
}

// for {
// 	select {

// 	default:

// 		switch state.ConnState {
// 		case START_CONN:
// 			state.ConnState = STOPPED
// 			fmt.Println("Entering state STOPPED")

// 		case STOPPED:
// 			apduReceived = <-state.chans.received
// 			fmt.Println("<<RX:", apduReceived)

// 			if apduReceived.Apci.FrameFormat == FrameFormat(StartDTAct) {
// 				// startdt act received
// 				apduToSend = NewApdu()
// 				apduToSend.Apci.FrameFormat = UFormatFrame
// 				apduToSend.Apci.UFormat = StartDTAct

// 				state.chans.toSend <- apduToSend
// 				state.ConnState = STARTED
// 				fmt.Println("Entering state STARTED")

// 			}

// 		case STARTED:
// 			apduReceived = <-state.chans.received
// 			fmt.Println("apdu received in state machine")

// 		case PENDING_UNCONFIRMED_STOPPED:

// 		case STOP_CONN:

// 		}

// 	case <-state.ctx.Done():
// 		fmt.Println("serverStateMachine received ctx.Done, returns")
// 		return

// 	}
// }
// }

func incrementSeqNumber(seqNumber SeqNumber) SeqNumber {

	return SeqNumber((int(seqNumber) + 1) % 32768)
}
