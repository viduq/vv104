package vv104

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State struct {
	Config      Config
	connState   ConnState
	ssn         SeqNumber
	rsn         SeqNumber
	chans       AllChans
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	dt_act_sent UFormat // for notification of state machine if a startdt_act or stopdt_act was sent
	tickers     tickers
}
type tickers struct {
	t1ticker time.Ticker
	t2ticker time.Ticker
	t3ticker time.Ticker
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
			Mode:            "",
			Ipv4Addr:        "",
			Port:            0,
			Casdu:           0,
			AutoAck:         false,
			K:               0,
			W:               0,
			T1:              0,
			T2:              0,
			T3:              0,
			IoaStructured:   false,
			InteractiveMode: false,
			UseLocalTime:    false,
		},
		connState: 0,
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
		if state.Config.InteractiveMode {
			go readCommandsFromStdIn(state.chans.commandsFromStdin) // never exits
		}
	}

	for {
		state.chans.received = make(chan Apdu, state.Config.W)
		state.chans.toSend = make(chan Apdu, state.Config.K)
		state.ctx, state.cancel = context.WithCancel(context.Background())
		// always start evaluateInteractiveCommands, we need it to control automatic sending, even if InteractiveMode is off
		go state.evaluateInteractiveCommands()
		go state.startConnection()

		<-state.ctx.Done()
		state.wg.Wait()
		// fmt.Println("Restart!")
		time.Sleep(1500 * time.Millisecond)
	}
}

func (state *State) connectionStateMachine() {
	var apduToSend Apdu
	var apduReceived Apdu
	isServer := state.Config.Mode == "server"
	isClient := state.Config.Mode == "client"

	state.wg.Add(1)
	defer state.wg.Done()

	state.connState = STOPPED
	fmt.Println("Entering state STOPPED")

	if isClient {
		// we need to trigger stardt_act here, it will trigger a notification for the blocking received channel, to jump over it
		state.chans.commandsFromStdin <- "startdt_act"
	}

	for {
		select {

		// block until apdu is received. some apdus are used as internal notifications with special type ids (are not sent)
		case apduReceived = <-state.chans.received:
			if apduReceived.Asdu.TypeId < INTERNAL_STATE_MACHINE_NOTIFIER {
				// real apdu received, not an internal notification
				fmt.Println("<<RX:", apduReceived)
				state.tickers.t3ticker.Reset(time.Duration(state.Config.T3) * time.Second)
			}

			if apduReceived.Apci.UFormat == TestFRAct {
				// always reply to testframes
				apduToSend = NewApdu()
				apduToSend.Apci.FrameFormat = UFormatFrame
				apduToSend.Apci.UFormat = TestFRCon
				state.chans.toSend <- apduToSend
				continue
			}

			// state machine
			switch state.connState {

			case STOPPED:
				if isServer {
					if apduReceived.Apci.UFormat == StartDTAct {
						// startdt act received
						apduToSend = NewApdu()
						apduToSend.Apci.FrameFormat = UFormatFrame
						apduToSend.Apci.UFormat = StartDTCon
						state.chans.toSend <- apduToSend
						state.connState = STARTED
						fmt.Println("Entering state STARTED")
					}

				}
				if isClient && (state.dt_act_sent == StartDTAct) {
					state.dt_act_sent = 0
					state.connState = PENDING_STARTED
					fmt.Println("Entering state PENDING_STARTED")
				}

			case PENDING_STARTED:
				if apduReceived.Apci.UFormat == StartDTCon {

					fmt.Println("Entering state STARTED")
					state.connState = STARTED

				}

			case STARTED:
				if isServer {
					if apduReceived.Apci.UFormat == StopDTAct {
						// stopdt act received
						// todo if unconfirmed frames
						// state.ConnState = PENDING_UNCONFIRMED_STOPPED
						apduToSend.Apci.FrameFormat = UFormatFrame
						apduToSend.Apci.UFormat = StopDTCon
						state.chans.toSend <- apduToSend
						state.connState = STOPPED
						fmt.Println("Entering state STOPPED")

					}
				}
				if isClient {
					if state.dt_act_sent == StopDTAct {
						// we have sent stopdt act as a client
						state.dt_act_sent = 0
						// todo if unconfirmed frames
						// state.ConnState = PENDING_UNCONFIRMED_STOPPED

						state.connState = PENDING_STOPPED
						fmt.Println("Entering state PENDING_STOPPED")
					}
				}
			case PENDING_STOPPED:
				if apduReceived.Apci.UFormat == StopDTCon {
					// we have sent stopdt act as a client OR received Stopdt con (whichever comes first)
					// bug: we could receive stopdt_con twice without having sent stopdt_act
					fmt.Println("Entering state STOPPED")
					state.connState = STOPPED

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
