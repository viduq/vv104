package vv104

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State struct {
	Config                Config
	Iec104ConnectionState ConnectionState
	TcpConnected          bool
	ssn                   SeqNumber
	rsn                   SeqNumber
	chans                 AllChans
	wg                    sync.WaitGroup
	ctx                   context.Context
	cancel                context.CancelFunc
}

type ConnectionState int

type AllChans struct {
	commandsFromStdin chan string
	received          chan Apdu
	toSend            chan Apdu
}

const (
	STOPPED ConnectionState = iota
	STARTED
	PENDING_UNCONFIRMED_STOPPED
	PENDING_STARTED
	PENDING_STOPPED
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
		Iec104ConnectionState: 0,
		TcpConnected:          false,
		ssn:                   0,
		rsn:                   0,
		chans:                 AllChans{},
		// wg:                    sync.WaitGroup{},
	}
}

func (state *State) Start() {
	state.Config.ParseFlags()
	state.chans.commandsFromStdin = make(chan string, 30)
	go readCommandsFromStdIn(state.chans.commandsFromStdin) // never exits

	for {
		state.chans.received = make(chan Apdu, state.Config.W)
		state.chans.toSend = make(chan Apdu, state.Config.K)
		state.ctx, state.cancel = context.WithCancel(context.Background())
		go state.startConnection()

		select {
		case input := <-state.chans.commandsFromStdin:

			if input == "exit" {
				fmt.Println("called exit")
				state.cancel()
				state.wg.Wait()
				fmt.Println("Start() exited")
				return
			}

		case <-state.ctx.Done():
			state.wg.Wait()
			fmt.Println("Restart!")
			time.Sleep(1500 * time.Millisecond)
			continue

		}
	}

}

func incrementSeqNumber(seqNumber SeqNumber) SeqNumber {

	return SeqNumber((int(seqNumber) + 1) % 32768)
}
