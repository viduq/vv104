package vv104

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func (state *State) startConnection() {
	if state.Config.Mode == "server" {
		fmt.Println("Starting Server")
		state.startServer()
	} else if state.Config.Mode == "client" {
		fmt.Println("Starting Client")
		state.startClient()
	} else {
		panic("can not start, config mode is neither server nor client")
	}
}

func (state *State) startServer() {
	fmt.Println("startServer started")
	defer fmt.Println("startServer returned")

	var err error

	ipAndPortStr := state.Config.Ipv4Addr + ":" + fmt.Sprint(state.Config.Port)
	ipAndPort, err := net.ResolveTCPAddr("tcp", ipAndPortStr)
	if err != nil {
		panic(err)
	}

	var l *net.TCPListener
	l, err = net.ListenTCP("tcp", ipAndPort)
	if err != nil {
		panic(err)
	}

	l.SetDeadline(time.Now().Add(2 * time.Second))
	defer l.Close()

	state.wg.Add(1)
	defer state.wg.Done()

	for {
		select {
		default:

			conn, err := l.Accept()
			if err != nil {
				if err, ok := err.(*net.OpError); ok && err.Timeout() {
					// it was a timeout
					// fmt.Println("timeout")a
					l.SetDeadline(time.Now().Add(2 * time.Second))
					continue
				}
				// other problem
				fmt.Println("accept error (not timeout)", err)
				continue
			}
			fmt.Println("Connected from: ", conn.RemoteAddr())
			go state.receivingRoutine(conn)
			go state.sendingRoutine(conn)
			go state.connectionStateMachine()
			go state.timerRoutine()

			<-state.ctx.Done() // todo? other criteria?
			return

		case <-state.ctx.Done():
			fmt.Println("startServer received Done(), returns")
			return
		}
	}

}

func (state *State) startClient() {
	fmt.Println("startClient started")
	defer fmt.Println("startClient returned")

	var err error

	ipAndPortStr := state.Config.Ipv4Addr + ":" + fmt.Sprint(state.Config.Port)
	ipAndPort, err := net.ResolveTCPAddr("tcp", ipAndPortStr)
	if err != nil {
		panic(err)
	}

	state.wg.Add(1)
	defer state.wg.Done()

	for {
		select {
		default:

			conn, err := net.DialTCP("tcp", nil, ipAndPort)
			if err != nil {
				fmt.Println("dial error", err)
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Println("Connected to:", conn.RemoteAddr())
			go state.receivingRoutine(conn)
			go state.sendingRoutine(conn)
			go state.connectionStateMachine()
			go state.timerRoutine()

			<-state.ctx.Done() // todo? other criteria?
			return

		case <-state.ctx.Done():
			fmt.Println("startClient received Done(), returns")
			return
		}
	}

}

func (state *State) receivingRoutine(conn net.Conn) {
	fmt.Println("receivingRoutine started")
	defer fmt.Println("receivingRoutine returned")
	defer conn.Close()
	var bytesbuf bytes.Buffer
	buf := make([]byte, 256) // todo: read multiple tcp frames from a whole tcp frame
	state.wg.Add(1)
	defer state.wg.Done()

	for {
		select {

		default:
			err := conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				fmt.Println(err)
			}
			recvLen, err := conn.Read(buf)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					// fmt.Println(err)
					continue
				}
				fmt.Println("Error reading:", err.Error())
				fmt.Println("Restart because of error reading, receivingRoutine returns")
				state.cancel()
				return
			}

			bytesbuf.Write(buf[:recvLen]) // Read from conn directly into bytesbuf?
			var receivedApdu Apdu
			receivedApdu, err = ParseApdu(&bytesbuf)
			bytesbuf.Reset()
			if err != nil {
				fmt.Println("error parsing:", err)
				fmt.Println("bytes:", bytesbuf)
			} else {
				// fmt.Println("<<RX:", apdu)
				state.chans.received <- receivedApdu
			}

		case <-state.ctx.Done():
			fmt.Println("receivingRoutine received Done(), returns")
			return
		}
	}
}

func (state *State) sendingRoutine(conn net.Conn) {
	fmt.Println("sendingRoutine started")
	defer fmt.Println("sendingRoutine returned")
	defer conn.Close()
	var apduToSend Apdu
	var buf []byte
	var err error
	state.wg.Add(1)
	defer state.wg.Done()

	for {
		select {

		case apduToSend = <-state.chans.toSend:

			buf, err = apduToSend.Serialize(*state)
			if err != nil {
				fmt.Println("error serializing apdu", err)
				continue
			}

			if apduToSend.Apci.UFormat == StopDTAct || apduToSend.Apci.UFormat == StartDTAct {
				// notify state machine
				state.dt_act_sent = apduToSend.Apci.UFormat
				apduNotify := NewApdu()
				apduNotify.Asdu.TypeId = INTERNAL_STATE_MACHINE_NOTIFIER
				state.chans.received <- apduNotify
			}

			if state.connState != STARTED {
				if apduToSend.Apci.FrameFormat == IFormatFrame {
					fmt.Println("IEC 104 connection is not started. Can not send I-Format")
					continue
				}
			}
			fmt.Println("TX>>:", apduToSend)
			err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				fmt.Println(err)
			}
			_, err = conn.Write(buf)
			if err != nil {
				fmt.Println("error sending apdu", err)
				fmt.Println("Error sending:", err.Error())
				fmt.Println("Restart because of error sending, sendingRoutine returns")
				state.cancel()
				return
			}

		case <-state.ctx.Done():
			fmt.Println("sendingRoutine received Done(), returns")
			return
		}
	}
}

func (state *State) timerRoutine() {
	fmt.Println("timerRoutine started")
	defer fmt.Println("timerRoutine returned")
	state.wg.Add(1)
	defer state.wg.Done()

	state.tickers.t1ticker = *time.NewTicker(time.Duration(state.Config.T1) * time.Second)
	state.tickers.t2ticker = *time.NewTicker(time.Duration(state.Config.T2) * time.Second)
	state.tickers.t3ticker = *time.NewTicker(time.Duration(state.Config.T3) * time.Second)

	for {
		select {

		case <-state.tickers.t1ticker.C:
			fmt.Println("t1 TIMEOUT")
		case <-state.tickers.t2ticker.C:
			fmt.Println("t2 TIMEOUT")
		case <-state.tickers.t3ticker.C:
			fmt.Println("t3 TIMEOUT")
			state.chans.commandsFromStdin <- "testfr_act"

		case <-state.ctx.Done():
			fmt.Println("timerRoutine received Done(), returns")
			return
		}
	}
}

func checkIpV4Address(ipAddr string) bool {

	return net.ParseIP(ipAddr) != nil
}
