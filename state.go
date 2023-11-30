package vv104

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"time"
)

type State struct {
	Config      Config
	ConnState   ConnState
	sendAck     ack
	recvAck     ack
	Chans       AllChans
	Wg          sync.WaitGroup
	Ctx         context.Context
	Cancel      context.CancelFunc
	dt_act_sent UFormat // for notification of state machine if a startdt_act or stopdt_act was sent
	tickers     tickers
}
type tickers struct {
	t1ticker *time.Ticker
	t2ticker *time.Ticker
	t3ticker *time.Ticker
}

type ConnState int

type AllChans struct {
	commandsFromStdin chan string
	Received          chan Apdu
	ToSend            chan Apdu
}

type ack struct {
	seqNumber  SeqNumber
	openFrames int
	ring       *ring.Ring
}

type seqNumberAndTimetag struct {
	seqNumber SeqNumber
	timetag   time.Time
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
		Config:    Config{Mode: "", Ipv4Addr: "", Port: 0, Casdu: 0, AutoAck: false, K: 0, W: 0, T1: 0, T2: 0, T3: 0, IoaStructured: false, InteractiveMode: false, UseLocalTime: false},
		ConnState: 0,
		sendAck:   ack{},
		recvAck:   ack{},
		Chans:     AllChans{},
		Wg:        sync.WaitGroup{},
		Ctx:       nil,
		Cancel: func() {
		},
		dt_act_sent: 0,
		tickers:     tickers{},
	}
}

func (state *State) Start() {
	if state.Config.InteractiveMode {
		state.Chans.commandsFromStdin = make(chan string, 30)
		if state.Config.InteractiveMode {
			go readCommandsFromStdIn(state.Chans.commandsFromStdin) // never exits
		}
	}

	for {
		state.Chans.Received = make(chan Apdu, state.Config.W)
		state.Chans.ToSend = make(chan Apdu, state.Config.K)
		state.Ctx, state.Cancel = context.WithCancel(context.Background())

		state.sendAck = newAck(state.Config.K)
		state.recvAck = newAck(state.Config.K)

		// always start evaluateInteractiveCommands, we need it to control automatic sending, even if InteractiveMode is off
		go state.evaluateInteractiveCommands()
		go state.startConnection()

		<-state.Ctx.Done()
		state.Wg.Wait()
		// fmt.Println("Restart!")
		time.Sleep(1500 * time.Millisecond)
	}
}

func (state *State) connectionStateMachine() {
	var apduToSend Apdu
	var apduReceived Apdu
	isServer := state.Config.Mode == "server"
	isClient := state.Config.Mode == "client"

	state.Wg.Add(1)
	defer state.Wg.Done()

	state.ConnState = STOPPED
	fmt.Println("Entering state STOPPED")

	if isClient {
		// we need to trigger stardt_act here, it will trigger a notification for the blocking received channel, to jump over it
		state.Chans.commandsFromStdin <- "startdt_act"
	}

	for {
		select {

		// block until apdu is received. some apdus are used as internal notifications with special type ids (are not sent)
		case apduReceived = <-state.Chans.Received:
			if (apduReceived.Apci.FrameFormat != IFormatFrame) || apduReceived.Asdu.TypeId < INTERNAL_STATE_MACHINE_NOTIFIER {
				// real apdu received, not an internal notification
				fmt.Println("<<RX:", apduReceived)
				state.tickers.t3ticker.Reset(time.Duration(state.Config.T3-4) * time.Second)
			}

			if apduReceived.Apci.UFormat == TestFRAct {
				// always reply to testframes
				apduToSend = NewApdu()
				apduToSend.Apci.FrameFormat = UFormatFrame
				apduToSend.Apci.UFormat = TestFRCon
				state.Chans.ToSend <- apduToSend
				continue
			}

			// state machine
			switch state.ConnState {

			case STOPPED:
				if isServer {
					if apduReceived.Apci.UFormat == StartDTAct {
						// startdt act received
						apduToSend = NewApdu()
						apduToSend.Apci.FrameFormat = UFormatFrame
						apduToSend.Apci.UFormat = StartDTCon
						state.Chans.ToSend <- apduToSend
						state.ConnState = STARTED
						fmt.Println("Entering state STARTED")
					}

				}
				if isClient && (state.dt_act_sent == StartDTAct) {
					state.dt_act_sent = 0
					state.ConnState = PENDING_STARTED
					fmt.Println("Entering state PENDING_STARTED")
				}

			case PENDING_STARTED:
				if apduReceived.Apci.UFormat == StartDTCon {

					fmt.Println("Entering state STARTED")
					state.ConnState = STARTED

				}

			case STARTED:
				if isServer {
					if apduReceived.Apci.UFormat == StopDTAct {
						// stopdt act received
						// todo if unconfirmed frames
						// state.ConnState = PENDING_UNCONFIRMED_STOPPED
						apduToSend.Apci.FrameFormat = UFormatFrame
						apduToSend.Apci.UFormat = StopDTCon
						state.Chans.ToSend <- apduToSend
						state.ConnState = STOPPED
						fmt.Println("Entering state STOPPED")

					}
				}
				if isClient {
					if state.dt_act_sent == StopDTAct {
						// we have sent stopdt act as a client
						state.dt_act_sent = 0
						// todo if unconfirmed frames
						// state.ConnState = PENDING_UNCONFIRMED_STOPPED

						state.ConnState = PENDING_STOPPED
						fmt.Println("Entering state PENDING_STOPPED")
					}
				}
			case PENDING_STOPPED:
				if apduReceived.Apci.UFormat == StopDTCon {
					// we have sent stopdt act as a client OR received Stopdt con (whichever comes first)
					// bug: we could receive stopdt_con twice without having sent stopdt_act
					fmt.Println("Entering state STOPPED")
					state.ConnState = STOPPED

				}
			}

		case <-state.Ctx.Done():
			fmt.Println("serverStateMachine received ctx.Done, returns")
			return
		}
	}
}

func incrementSeqNumber(seqNumber SeqNumber) SeqNumber {

	return SeqNumber((int(seqNumber) + 1) % 32768)
}

func newAck(length int) ack {
	ack := ack{}
	ack.openFrames = 0
	ack.seqNumber = 0
	ack.ring = ring.New(length)

	for i := 0; i < ack.ring.Len(); i++ {
		ack.ring.Value = seqNumberAndTimetag{
			seqNumber: 0,
			timetag:   time.Time{},
		}
		ack.ring = ack.ring.Next()
	}
	return ack
}

// queueApdu is called when we receive i-formats, that need to be ack'ed
func (ack *ack) queueApdu(apdu Apdu) {
	ack.ring.Value = seqNumberAndTimetag{
		seqNumber: apdu.Apci.Ssn,
		timetag:   time.Now(),
	}

	ack.seqNumber = incrementSeqNumber(ack.seqNumber)
	ack.openFrames++
	ack.ring = ack.ring.Next()

}

// ackApdu is called when we send an i- or s-format and acknowledge received frames
func (ack *ack) ackApdu(seqNumber SeqNumber, t2ticker *time.Ticker, t2 time.Duration) {
	var still_unacked int = 0

	for i := 0; i < ack.ring.Len(); i++ {
		fmt.Println(ack.ring)
		ack.ring = ack.ring.Next()
	}

	// we go back in the ring to find the ack'ed sequence number
	// the more we have to go back, the more frames are still unack'ed
	for still_unacked = 0; still_unacked < ack.openFrames; still_unacked++ {
		if ack.ring.Value.(seqNumberAndTimetag).seqNumber == seqNumber {
			// all acknowledged
			fmt.Println("all acked", seqNumber)
			fmt.Println("still unacked:", still_unacked)
			ack.openFrames = still_unacked
			timetag := ack.ring.Value.(seqNumberAndTimetag).timetag
			frameIsUnackedTime := time.Now().Sub(timetag)
			frameMustBeAckedIn := t2 - frameIsUnackedTime
			fmt.Println("will be acked in ", frameMustBeAckedIn)
			t2ticker.Reset(frameMustBeAckedIn * time.Second)

			break
		} else {
			ack.ring = ack.ring.Prev()

		}

	}
}

func (ack *ack) checkForAck(maxOpenFrames int) (bool, SeqNumber) {
	if ack.openFrames > maxOpenFrames {
		seqNumber := ack.seqNumber
		fmt.Println("seq number to ack:", seqNumber)
		fmt.Println("openFrames:", ack.openFrames)

		return true, seqNumber
	}
	return false, 0
}
