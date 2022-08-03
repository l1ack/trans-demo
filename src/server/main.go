package main

import (
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"fmt"

	xmit "e.coding.net/xverse-git/xmedia/xmit-lib/go/xmit"
)

var (
	address = "43.138.64.40"
	demo    DemoServer
)

func IPV4ToUint(ip string) uint32 {
	ip = strings.TrimSpace(ip)
	bits := strings.Split(ip, ".")
	if len(bits) != 4 {
		return 0
	}

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	if b0 > 255 || b1 > 255 || b2 > 255 || b3 > 255 {
		return 0
	}

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Printf("I: syntax: %s <endpoint>\n", os.Args[0])
	// 	return
	// }
	server, _ := zmq.NewSocket(zmq.REP)
	server.Bind("tcp://*:5558")
	fmt.Println("connect to bridge!")

	_msg, err := server.RecvBytes(zmq.SNDMORE)
	server.SendMessage("shoudao")
	//
	newTest := &Msg{}
	println("yes recv")
	err = proto.Unmarshal(_msg, newTest)
	if err != nil {
		fmt.Println(err)
	}
	msg := newTest.Payload
	transMode := newTest.Sequence
	// size := _msg.Sequence
	fmt.Println(msg)
	demo.Start(IPV4ToUint(msg), transMode)

	xmit.Join()
}
