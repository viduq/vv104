package vv104

import (
	"container/ring"
	"context"
	"sync"
	"time"
)

type State struct {
	Config           Config
	ConnState        ConnState
	Chans            AllChans
	Objects          *Objects
	Running          bool // true when trying to connect/waiting for connection
	TcpConnected     bool
	Ctx              context.Context
	Wg               sync.WaitGroup
	Cancel           context.CancelFunc
	dt_act_sent      UFormat // for notification of state machine if a startdt_act or stopdt_act was sent
	manualDisconnect bool
	sendAck          ack
	recvAck          ack
	tickers          tickers
}
type tickers struct {
	t1ticker              *time.Ticker
	t2tickerReceivedItems *time.Ticker
	t2tickerSentItems     *time.Ticker

	t3ticker *time.Ticker
}

type ConnState int

type AllChans struct {
	CommandsFromStdin chan string
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
	objects := NewObjects()
	return State{
		Config:    Config{},
		ConnState: 0,
		sendAck:   ack{},
		recvAck:   ack{},
		Chans:     AllChans{},
		Objects:   objects,
		Wg:        sync.WaitGroup{},
		Ctx:       nil,
		Cancel: func() {
		},
		dt_act_sent:      0,
		manualDisconnect: false,
		Running:          false,
		TcpConnected:     false,
		tickers:          tickers{},
	}
}

func (state *State) Start() {

	initLogger(state.Config)
	printConfig(state.Config)
	state.Running = true

	for !state.manualDisconnect {

		state.Chans.Received = make(chan Apdu, state.Config.W)
		state.Chans.ToSend = make(chan Apdu, state.Config.K)
		state.Ctx, state.Cancel = context.WithCancel(context.Background())

		state.sendAck = newAck(state.Config.K)
		state.recvAck = newAck(state.Config.K)

		if state.Config.InteractiveMode {
			state.Chans.CommandsFromStdin = make(chan string, 30)
			go state.readCommandsFromStdIn()
		}

		// always start evaluateInteractiveCommands, we need it to control automatic sending, even if InteractiveMode is off
		go state.evaluateInteractiveCommands()
		go state.startConnection()

		<-state.Ctx.Done()
		state.Wg.Wait()
		if !state.manualDisconnect {
			logInfo.Println("Restart!")

		}
		time.Sleep(1500 * time.Millisecond)
	}
	defer logDebug.Println("Start() returned")
	// disconnect was done purposely, exit
	state.manualDisconnect = false
	state.Running = false
}

func (state *State) connectionStateMachine() {
	var apduToSend Apdu
	var apduReceived Apdu
	isServer := state.Config.Mode == "server"
	isClient := state.Config.Mode == "client"

	state.Wg.Add(1)
	defer state.Wg.Done()

	state.ConnState = STOPPED
	logInfo.Println("Entering state STOPPED")

	if isClient {
		// we need to trigger stardt_act here, it will trigger a notification for the blocking received channel, to jump over it
		state.Chans.CommandsFromStdin <- "startdt_act"
	}

	for {
		select {

		// block until apdu is received. some apdus are used as internal notifications with special type ids (are not sent)
		case apduReceived = <-state.Chans.Received:
			if (apduReceived.Apci.FrameFormat != IFormatFrame) || apduReceived.Asdu.TypeId < INTERNAL_STATE_MACHINE_NOTIFIER {
				// real apdu received, not an internal notification
				logInfo.Println("<<RX:", state.Objects.objNameOrIoa(apduReceived.Asdu), apduReceived)
				state.tickers.t3ticker.Reset(time.Duration(state.Config.T3) * time.Second)
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
						logInfo.Println("Entering state STARTED")
					}

				}
				if isClient && (state.dt_act_sent == StartDTAct) {
					state.dt_act_sent = 0
					state.ConnState = PENDING_STARTED
					logInfo.Println("Entering state PENDING_STARTED")
				}

			case PENDING_STARTED:
				if apduReceived.Apci.UFormat == StartDTCon {

					logInfo.Println("Entering state STARTED")
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
						logInfo.Println("Entering state STOPPED")

					}
				}
				if isClient {
					if state.dt_act_sent == StopDTAct {
						// we have sent stopdt act as a client
						state.dt_act_sent = 0
						// todo if unconfirmed frames
						// state.ConnState = PENDING_UNCONFIRMED_STOPPED

						state.ConnState = PENDING_STOPPED
						logInfo.Println("Entering state PENDING_STOPPED")
					}
				}
			case PENDING_STOPPED:
				if apduReceived.Apci.UFormat == StopDTCon {
					// we have sent stopdt act as a client OR received Stopdt con (whichever comes first)
					// bug: we could receive stopdt_con twice without having sent stopdt_act
					logInfo.Println("Entering state STOPPED")
					state.ConnState = STOPPED

				}
			}

		case <-state.Ctx.Done():
			logDebug.Println("serverStateMachine received ctx.Done, returns")
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

// queueApdu adds i-formats to the ring, because they need to be ack'ed
func (ack *ack) queueApdu(apdu Apdu) {
	// check if received ssn is okay
	if apdu.Apci.Ssn != ack.seqNumber {
		logError.Println("Error received SSn does not match internal state, received ssn: ", apdu.Apci.Ssn, "state: ", ack.seqNumber)
	}

	ack.ring = ack.ring.Next()
	ack.ring.Value = seqNumberAndTimetag{
		seqNumber: apdu.Apci.Ssn,
		timetag:   time.Now(),
	}

	ack.seqNumber = incrementSeqNumber(ack.seqNumber)
	ack.openFrames++

}

// ackApdu is called when we send an i- or s-format and acknowledge received frames
// or if we receive an i- or s-format which acknowledges sent frames
func (ack *ack) ackApdu(seqNumber SeqNumber, t2ticker *time.Ticker, t2 time.Duration) {
	var stillUnacked int = 0

	// we go back in the ring to find the ack'ed sequence number
	// the more we have to go back, the more frames are still unack'ed
	for stillUnacked = 0; stillUnacked < ack.openFrames; stillUnacked++ {
		// fmt.Printf("%t\n", ack.ring.Value)
		if ack.ring.Value.(seqNumberAndTimetag).seqNumber == seqNumber-1 {
			// all until this seq number are acknowledged
			// we might have received more already, which are still open (stillUnAcked)
			// fmt.Println("all acked until", seqNumber)
			// fmt.Println("still unacked:", still_unacked)
			ack.openFrames = stillUnacked
			if stillUnacked > 0 {

				timetag := ack.ring.Value.(seqNumberAndTimetag).timetag
				// fmt.Println("ttag", timetag)
				// fmt.Println("tnow", time.Now())

				frameIsUnackedTime := time.Now().Sub(timetag)
				// fmt.Println("frameisUnacked", frameIsUnackedTime)
				// fmt.Println("t2:", t2)

				frameMustBeAckedIn := t2 - frameIsUnackedTime
				// fmt.Println("will be acked in ", frameMustBeAckedIn)
				t2ticker.Reset(frameMustBeAckedIn)
			} else {
				// we ack'ed all items, stop ticker
				// TODO: Ticker should be mutexed, because we might stop it here, although it was just started by a received I frame!
				t2ticker.Stop()
			}

			break
		}

		ack.ring = ack.ring.Prev()

	}
}

func (ack *ack) checkForAck(maxOpenFrames int) (bool, SeqNumber) {
	if ack.openFrames >= maxOpenFrames {
		seqNumber := ack.seqNumber
		// fmt.Println("seq number to ack:", seqNumber)
		// fmt.Println("openFrames:", ack.openFrames)

		return true, seqNumber
	}
	return false, 0
}
