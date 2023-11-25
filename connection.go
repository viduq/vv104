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
		panic("can not start, config mode is nor server nor client")
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
			go state.handleConnection(conn)

			<-state.ctx.Done() // todo? other criteria?
			return

		case <-state.ctx.Done():
			fmt.Println("startServer received Done(), returns")
			return
		}
	}

}

func (state *State) startClient() {

}

func (state *State) handleConnection(conn net.Conn) {
	fmt.Println("handleConnection started")
	defer fmt.Println("handleConnection returned")
	defer conn.Close()
	var bytesbuf bytes.Buffer
	buf := make([]byte, 256) // todo: read whole tcp frame
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
				fmt.Println("Restart because of error reading, handleConnection returns")
				state.cancel()
				return
			}

			bytesbuf.Write(buf[:recvLen]) // Read from conn directly into bytesbuf?
			fmt.Println(bytesbuf)
			apdu, err := ParseApdu(&bytesbuf)
			if err != nil {
				fmt.Println("error parsing:", err)
			} else {
				fmt.Println("<<RX:", apdu)
			}

		case <-state.ctx.Done():
			fmt.Println("handleConnection received Done(), returns")
			return
		}
	}
}

func checkIpV4Address(ipAddr string) bool {

	return net.ParseIP(ipAddr) != nil
}
