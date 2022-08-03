package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/pion/randutil"

	xmit "e.coding.net/xverse-git/xmedia/xmit-lib/go/xmit"
)

type (
	DemoClient struct {
		connId        []byte
		client        *xmit.Client
		callbacks     *xmit.ClientCallbacks
		serverAddress uint32
		dataStream    *xmit.ClientStream
		ctrlStream    *xmit.ClientStream
	}
)

var (
	seq            uint64
	seqrecv        uint64
	server2a       *zmq.Socket
	last_seq       uint64
	lost_seq       uint64
	fail_flag      uint64
	interrupt_time int
<<<<<<< HEAD

=======
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
)

func (demo *DemoClient) recvHandler(msgType uint16, data []byte) {
	switch msgType {
	case 1:
		demo.ctrlStream.Send(data)
		xmit.Stop()
		fmt.Println("Close Connection")
	case 2:
		seqrecv := *(*uint64)(unsafe.Pointer(&data[0]))
		c := data[24:len(data)]
		md5_ := md5.Sum(c)
		md5recv := uint64(md5_[0])
		md5 := *(*uint64)(unsafe.Pointer(&data[8]))
		sendtime := *(*uint64)(unsafe.Pointer(&data[16]))
		fmt.Println("Msg: (", seqrecv, ", TYPE: ", msgType, ", SIZE: ", len(data), " RTT: ", (uint64(time.Now().UnixNano())-sendtime)/1000000, ") Echoed")
		if last_seq != seqrecv-1 {
			lost_seq = lost_seq + 1
			println("lost_ratio: ", lost_seq)
			println("lost_seq: ", last_seq+1)
		}
		last_seq = seqrecv
		msg2a := &Msg2A{}
		msg2a.Timestampsend = sendtime
		msg2a.Seq = seqrecv
		msg2a.Timestamprecv = uint64(time.Now().UnixNano())
		msg2a.Mismatch = false
		if md5 != md5recv {
			// fmt.Println(c)
			fmt.Println("md5data", md5)
			fmt.Println("md5recv", md5recv)
			fmt.Println("mismatch")
			fmt.Println(c)
			msg2a.Mismatch = true
		}
		msg2a.Msgtype = uint64(msgType)
		msg2a.Size = uint64(len(data))
		_msg, err := proto.Marshal(msg2a)
		if err != nil {
			println("err: ", err)
		}
		server2a.SendBytes(_msg, zmq.DONTWAIT)
	}
}

func (demo *DemoClient) connectedHandler() {
	fail_flag = 5
	fmt.Println("Connected!")
}

func (demo *DemoClient) failtoconnect() {
	fail_flag = 1
}

<<<<<<< HEAD
func (demo *DemoClient) interrupt(intervalms int){
	interrupt_time = intervalms+3
=======
func (demo *DemoClient) interrupt(intervalms int) {
	interrupt_time = intervalms + 3
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
}

func (demo *DemoClient) sizeSet(size int) {
	// 把生成的对象和seq绑定在一起
	message := FixSeq(size)

	fmt.Println("inter: ", interrupt_time)
	fmt.Println("fail: ", fail_flag)

	// Send the message as text
	data := []byte(message)
	c := data[24:size]
	// fmt.Println(c)
	md5_ := md5.Sum(c)
	md5 := uint64(md5_[0])
	sendtime := uint64(time.Now().UnixNano())
	//fmt.Println("generate md5: ", md5)
	*(*uint64)(unsafe.Pointer(&data[8])) = md5
	//fmt.Println("md5send: ", md5_, "size: ", size)
	*(*uint64)(unsafe.Pointer(&data[0])) = seq
	*(*uint64)(unsafe.Pointer(&data[16])) = sendtime
	// data hash产生10个数
	seq++
	if err := demo.dataStream.Send(data); err != nil {
		fmt.Println("SEND ERR:", err)
	}

	// xmit.Join()
}



func (demo *DemoClient) Start(serverAddress uint32) {
	rand.Seed(time.Now().UTC().UnixNano())
	demo.serverAddress = serverAddress
	demo.client, _ = xmit.GetClient(xmit.Transport_kUdp)
	demo.callbacks = &xmit.ClientCallbacks{
		OnMessage:       demo.recvHandler,
		OnConnected:     demo.connectedHandler,
		OnFailtoconnect: demo.failtoconnect,
		OnInterrupted:   demo.interrupt,
	}

	xmit.Start()

	xmit.SignalServerAddress = serverAddress
	offer, err := demo.client.Offer(demo.callbacks)

<<<<<<< HEAD

=======
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
	// fmt.Println(offer)
	if err == nil {
		fmt.Println("OFFER:", string(offer))
		fmt.Printf("flag")
	} else {
		fmt.Println("OFFER ERR:", err)
		return
	}

	xmit.GetSignal().SendOffer(offer)
	answer := xmit.GetSignal().RecvAnswer()
	err = demo.client.Answer(answer)
	if err != nil {
		fmt.Println("ANSWER ERR:", err)
		return
	}

	// NewStream(msgType uint16, mode int, priority int, lifetimeMs int) (stream *ClientStream, err error)
<<<<<<< HEAD

=======
	demo.ctrlStream, _ = demo.client.NewStream(1, 3, 1, 0)
	demo.dataStream, _ = demo.client.NewStream(2, 3, 1, 0)
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a

	server2a, _ = zmq.NewSocket(zmq.PUSH)
	server2a.Connect("tcp://127.0.0.1:6666")
	seq = 0
	last_seq = 0
	lost_seq = 0
	// seq := uint64(0)
	// for range time.NewTicker(3 * time.Second).C {
	// 	message := RandSeq(rand.Intn(11111))

	// 	// Send the message as text
	// 	data := []byte(message)
	// 	*(*uint64)(unsafe.Pointer(&data[0])) = seq
	// 	seq++
	// 	if err = demo.dataStream.Send(data); err != nil {
	// 		fmt.Println("SEND ERR:", err)
	// 	}
	// }

}

<<<<<<< HEAD
<<<<<<< HEAD
func datastreamclose () {
=======
func datastreamctrl () {
>>>>>>> 4c1b76f... conn-2
	demo.ctrlStream, _ = demo.client.NewStream(1, 3, 1, 0)
}

func datastreamdata () {
	// demo.ctrlStream, _ = demo.client.NewStream(1, 3, 1, 0)
	demo.dataStream, _ = demo.client.NewStream(2, 3, 1, 0)
}
=======
// func datastreamctrl() {

// }

// func datastreamdata() {
// 	// demo.ctrlStream, _ = demo.client.NewStream(1, 3, 1, 0)

// }
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a

func RandSeq(n int) string {
	val, err := randutil.GenerateCryptoRandomString(n, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		fmt.Println("RAND ERR:", err)
	}

	return val
}

func FixSeq(n int) string {
	a := make([]rune, n)
	for i := range a {
		a[i] = 'a'
	}
	return string(a)
}
