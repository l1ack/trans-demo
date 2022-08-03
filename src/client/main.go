package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"fmt"

	xmit "e.coding.net/xverse-git/xmedia/xmit-lib/go/xmit"
)

var (
	address = "42.194.250.246"
	size    int
	demo    DemoClient
<<<<<<< HEAD
	init1   int
=======
	init1    int
>>>>>>> 4c1b76f... conn-2
	size1   int
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
<<<<<<< HEAD
	
		server, _ := zmq.NewSocket(zmq.PULL)
		server.Bind("tcp://127.0.0.1:5555")
		newTest := &Msg{}
<<<<<<< HEAD
		
		for{
=======

	server, _ := zmq.NewSocket(zmq.PULL)
	server.Bind("tcp://127.0.0.1:5555")
	newTest := &Msg{}
	init1 = 0
	size1 = 0
	_msg, err := server.RecvBytes(zmq.SNDMORE)
	err = proto.Unmarshal(_msg, newTest)
	if err != nil {
		fmt.Println(err)
	}
	msg := newTest.Payload
	// size := newTest.Sequence
	for {
		// size := _msg.Sequence
		//fmt.Println(msg)
		if size1 == 0 || init1 == 0 {
			demo.Start(IPV4ToUint(msg))
			// datastreamctrl()
			// datastreamdata()
			// NewStream(msgType uint16, mode int, priority int, lifetimeMs int) (stream *ClientStream, err error)
		}
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
=======
		init1 = 0
		size1 = 0
		for {
		init1 = 1
>>>>>>> 4c1b76f... conn-2
		_msg, err := server.RecvBytes(zmq.SNDMORE)
		err = proto.Unmarshal(_msg, newTest)
		if err != nil {
			fmt.Println(err)
		}
		msg := newTest.Payload
		size := newTest.Sequence
		// size := _msg.Sequence
		fmt.Println(msg)
		if (size1 == 0 || init1 == 0) {
		demo.Start(IPV4ToUint(msg))
		datastreamctrl()
		datastreamdata()
		// NewStream(msgType uint16, mode int, priority int, lifetimeMs int) (stream *ClientStream, err error)
	}
		fmt.Println(size)
<<<<<<< HEAD
<<<<<<< HEAD
		size1 := int(size)
		datastreamclose()
		datastreamstart()
=======
		size1 = int(size)
>>>>>>> 4c1b76f... conn-2
		if size1 == 0 {
			xmit.Stop()
			//continue
			//size1 = 500
		} else {
			demo.sizeSet(size1)
		}
			// mm, err := server.SendMessage("shoudao")
			// fmt.Println(mm)
<<<<<<< HEAD
		}
=======
		size1 = int(size)
		if size1 == 0 {
			// demo.sizeSet(300)
			time.Sleep(1000000000)
			xmit.Stop()
			//continue
			//size1 = 500
		} else {
			demo.sizeSet(size1)
		}
		init1 = 1
		time.Sleep(1000000000)
		// mm, err := server.SendMessage("shoudao")
		// fmt.Println(mm)
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
=======
>>>>>>> 4c1b76f... conn-2
	}
	xmit.Join()

}
