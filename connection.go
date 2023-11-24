package vv104

import (
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
	var err error
	var l net.Listener
	ipAndPort := state.Config.Ipv4Addr + ":" + fmt.Sprint(state.Config.Port)
	l, err = net.Listen("tcp", ipAndPort)
	if err != nil {
		panic(err)
	}

	// state.wg.Add(1)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 2)
			continue
		}
		fmt.Println("Connected from: ", conn.RemoteAddr())
		go state.handleConnection(conn)
	}

}

func (state *State) startClient() {

}

func (state *State) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 256) // todo: read whole tcp frame
	state.wg.Add(1)
	defer state.wg.Done()

	// err := conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	// if err != nil {
	// 	panic(err)
	// }
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
				fmt.Println("Restart because of error reading")
				state.cancel()
				return
			}
			fmt.Println(buf[:recvLen])

		case <-state.ctx.Done():
			fmt.Println("handleConnection received Done(), quitting")
			return
		}
	}
}

func checkIpV4Address(ipAddr string) bool {

	return net.ParseIP(ipAddr) != nil
}
