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

	state.Wg.Add(1)
	defer state.Wg.Done()

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

			<-state.Ctx.Done() // todo? other criteria?
			return

		case <-state.Ctx.Done():
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

	state.Wg.Add(1)
	defer state.Wg.Done()

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

			<-state.Ctx.Done() // todo? other criteria?
			return

		case <-state.Ctx.Done():
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
	state.Wg.Add(1)
	defer state.Wg.Done()

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
				state.Cancel()
				return
			}

			bytesbuf.Write(buf[:recvLen]) // Read from conn directly into bytesbuf?
			var receivedApdu Apdu
			receivedApdu, err = ParseApdu(&bytesbuf)
			bytesbuf.Reset()
			if err != nil {
				fmt.Println("error parsing:", err)
				fmt.Println("bytes:", bytesbuf)
				continue
			}

			if receivedApdu.Apci.FrameFormat == IFormatFrame {
				// check if received ssn is okay
				if receivedApdu.Apci.Ssn != state.recvAck.seqNumber {
					fmt.Println("Error received SSn does not match internal state, received ssn: ", receivedApdu.Apci.Ssn, "state: ", state.recvAck.seqNumber)
				}

				// each received I-Format must be acknowledged
				// this should be done directly after receiving (not in another goroutine, because of race conditions) (?)
				state.recvAck.queueApdu(receivedApdu)
				if state.recvAck.openFrames == 1 {
					// was 0 before, new open frame
					state.tickers.t2tickerReceivedItems.Reset(time.Duration(state.Config.T2) * time.Second)
				}
				weMustAck, seqNumberToAck := state.recvAck.checkForAck(state.Config.W)
				if weMustAck {
					// fmt.Println("we must ack received items because w values open")
					sframe := NewApdu()
					sframe.Apci.FrameFormat = SFormatFrame
					sframe.Apci.Rsn = seqNumberToAck
					state.Chans.ToSend <- sframe
				}
			}
			if receivedApdu.Apci.FrameFormat == IFormatFrame || receivedApdu.Apci.FrameFormat == SFormatFrame {
				// each received I- or S-Format acknowledges some of our sent frames
				state.sendAck.ackApdu(receivedApdu.Apci.Rsn, state.tickers.t2tickerSentItems, time.Duration(state.Config.T2)*time.Second)
			}

			state.Chans.Received <- receivedApdu

		case <-state.Ctx.Done():
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
	state.Wg.Add(1)
	defer state.Wg.Done()

	for {
		select {

		case apduToSend = <-state.Chans.ToSend:

			buf, err = apduToSend.Serialize(*state)
			// fmt.Println(buf)
			if err != nil {
				fmt.Println("error serializing apdu", err)
				continue
			}

			if apduToSend.Apci.UFormat == StopDTAct || apduToSend.Apci.UFormat == StartDTAct {
				// notify state machine
				state.dt_act_sent = apduToSend.Apci.UFormat
				apduNotify := NewApdu()
				apduNotify.Asdu.TypeId = INTERNAL_STATE_MACHINE_NOTIFIER
				state.Chans.Received <- apduNotify
			}

			if state.ConnState != STARTED {
				if apduToSend.Apci.FrameFormat == IFormatFrame {
					fmt.Println("IEC 104 connection is not started. Can not send I-Format")
					continue
				}
			} else {
				// started
				if state.sendAck.openFrames >= state.Config.K {
					// we must not send anymore, wait for acknowledgement
					fmt.Println("we must not send anymore, wait for acknowledgement TODO")
					// TODO block on a channel
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
				state.Cancel()
				return
			}
			if apduToSend.Apci.FrameFormat == SFormatFrame || apduToSend.Apci.FrameFormat == IFormatFrame {
				// by sending an s- or i-format we have acknowledged items
				state.recvAck.ackApdu(apduToSend.Apci.Rsn, state.tickers.t2tickerReceivedItems, time.Duration(state.Config.T2)*time.Second)
			}
			if apduToSend.Apci.FrameFormat == IFormatFrame {
				// each sent frame must be ack'ed by the communication partner in a certain time
				state.sendAck.queueApdu(apduToSend)

				if state.sendAck.openFrames == 1 {
					// was 0 before, new open frame
					state.tickers.t2tickerSentItems.Reset(time.Duration(state.Config.T2) * time.Second)
				}

			}

		case <-state.Ctx.Done():
			fmt.Println("sendingRoutine received Done(), returns")
			return
		}
	}
}

func (state *State) timerRoutine() {
	fmt.Println("timerRoutine started")
	defer fmt.Println("timerRoutine returned")
	state.Wg.Add(1)
	defer state.Wg.Done()

	state.tickers.t1ticker = time.NewTicker(time.Duration(state.Config.T1) * time.Second)
	state.tickers.t2tickerReceivedItems = time.NewTicker(time.Duration(state.Config.T2) * time.Second)
	state.tickers.t2tickerReceivedItems.Stop()
	state.tickers.t2tickerSentItems = time.NewTicker(time.Duration(state.Config.T2) * time.Second)
	state.tickers.t2tickerSentItems.Stop()
	state.tickers.t3ticker = time.NewTicker(time.Duration(state.Config.T3-4) * time.Second)

	for {
		select {

		// case <-state.tickers.t1ticker.C:
		// 	fmt.Println("t1 TIMEOUT")
		case <-state.tickers.t2tickerReceivedItems.C:
			if state.recvAck.openFrames > 0 {
				// fmt.Println("we must ack received items because t2 timeout")

				sframe := NewApdu()
				sframe.Apci.FrameFormat = SFormatFrame
				sframe.Apci.Rsn = state.recvAck.seqNumber
				state.Chans.ToSend <- sframe
			}

		case <-state.tickers.t2tickerSentItems.C:
			fmt.Println("the communication partner did not acknowledge in the specified time, quitting...")
			state.Cancel()

		case <-state.tickers.t3ticker.C:
			// fmt.Println("t3 TIMEOUT")
			state.Chans.CommandsFromStdin <- "testfr_act"

		case <-state.Ctx.Done():
			fmt.Println("timerRoutine received Done(), returns")
			return
		}
	}
}

func checkIpV4Address(ipAddr string) bool {

	return net.ParseIP(ipAddr) != nil
}
