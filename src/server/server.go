package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"
	"unsafe"

	"github.com/google/uuid"

	"e.coding.net/xverse-git/xmedia/xmit-lib/go/xmit"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

type (
	DemoServer struct {
		server        *xmit.Server
		callbacks     *xmit.ServerCallbacks
		serverAddress uint32
		dataStream    *xmit.ServerStream
		ctrlStream    *xmit.ServerStream
	}
)

var (
	seq      uint64
	seqrecv  uint64
	server2a *zmq.Socket
	last_seq uint64
	lost_seq uint64
	fail_flag      uint64
	interrupt_time int
)

func (demo *DemoServer) signalHttp(rsp http.ResponseWriter, req *http.Request) {
	req.Close = true
	offer, _ := io.ReadAll(req.Body)

	rsp.Header().Set("Content-Type", "application/octet-stream")
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Credentials", "true")
	rsp.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	if _, err := rsp.Write(demo.Accept(offer)); err != nil {
		fmt.Println("HTTP ERR:", err)
	}
}

func (demo *DemoServer) echoHandler(connId []byte, msgType uint16, data []byte) {
	switch msgType {
	case 1:
		demo.ctrlStream.Send(data)
		demo.server.Close(connId)
		fmt.Println("Close Connection:", string(connId))
	case 2:

		fmt.Println("inter: ", interrupt_time)
		fmt.Println("fail: ", fail_flag)

		demo.dataStream.Send(data)
		c := data[24:len(data)]
		md5_ := md5.Sum(c)
		md5recv := uint64(md5_[0])
		md5 := *(*uint64)(unsafe.Pointer(&data[8]))
		if md5 != md5recv {
			fmt.Println("md5byte", md5_)
			fmt.Println("mdedata", md5)
			fmt.Println("md5recv", md5recv)
			fmt.Println("mismatch")
			fmt.Println("data:", c)
		}
		seq := *(*uint64)(unsafe.Pointer(&data[0]))
		fmt.Println("Msg: (", seq, ", TYPE: ", msgType, ", SIZE: ", len(data), " TIME: ", time.Now().UTC(), ") Echoed")
	}
}

func (demo *DemoServer) recvHandler(connId []byte, msgType uint16, data []byte) {
	switch msgType {
	case 1:
		demo.ctrlStream.Send(data)
		demo.server.Close(connId)
		fmt.Println("Close Connection:", string(connId))
	case 2:

		fmt.Println("inter: ", interrupt_time)
		fmt.Println("fail: ", fail_flag)

		seqrecv := *(*uint64)(unsafe.Pointer(&data[0]))
		c := data[24:len(data)]
		md5_ := md5.Sum(c)
		md5recv := uint64(md5_[0])
		md5 := *(*uint64)(unsafe.Pointer(&data[8]))
		sendtime := *(*uint64)(unsafe.Pointer(&data[16]))

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

		fmt.Println("Msg: (", seqrecv, ", TYPE: ", msgType, ", SIZE: ", len(data), " RTT: ", (msg2a.Timestamprecv-sendtime)/1000000, ") Echoed")
		msg2a.Msgtype = uint64(msgType)
		msg2a.Size = uint64(len(data))
		_msg, err := proto.Marshal(msg2a)
		if err != nil {
			println("err: ", err)
		}
		server2a.SendBytes(_msg, zmq.DONTWAIT)
	}
}

func (demo *DemoServer) connectedHandler(connId []byte) {
	fail_flag = 5
	fmt.Println("Connected!")
}


func (demo *DemoServer) failtoconnect(connId []byte) {
<<<<<<< HEAD
<<<<<<< HEAD
=======
	fmt.Println("fail to connect!")
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
=======
	fmt.Println("fail to connect!")
>>>>>>> 4c1b76f... conn-2
	fail_flag = 1
}

func (demo *DemoServer) interrupt(connId []byte, intervalms int){
	interrupt_time = intervalms + 3
<<<<<<< HEAD
<<<<<<< HEAD
=======
	fmt.Println("interrupt time :", interrupt_time)
>>>>>>> d358fbdfa809829caa76f31a472f7c9fa700425a
=======
	fmt.Println("interrupt time :", interrupt_time)
>>>>>>> 4c1b76f... conn-2
}


func (demo *DemoServer) Start(serverAddress uint32, transMode uint64) {
	http.HandleFunc("/signal", demo.signalHttp)
	go http.ListenAndServe(":51913", nil)

	demo.serverAddress = serverAddress
	demo.server, _ = xmit.GetServer(serverAddress, 0, 0)
	// case 1, no back
	// case 2, echo
	switch transMode {
	case 1:
		demo.callbacks = &xmit.ServerCallbacks{
			OnMessage:       demo.recvHandler,
			OnConnected:     demo.connectedHandler,
			OnFailtoconnect: demo.failtoconnect,
			OnInterrupted:   demo.interrupt,
		}
	case 2:
		demo.callbacks = &xmit.ServerCallbacks{
			OnMessage:   demo.echoHandler,
			OnConnected: demo.connectedHandler,
			OnFailtoconnect: demo.failtoconnect,
			OnInterrupted:   demo.interrupt,
		}

	}

	xmit.Start()

	xmit.SignalServerAddress = serverAddress

	server2a, _ = zmq.NewSocket(zmq.PUSH)
	server2a.Connect("tcp://127.0.0.1:6666")
	seq = 0
	last_seq = 0
	lost_seq = 0

	go func(demo *DemoServer) {
		for {
			xmit.GetSignal().SendAnswer(demo.Accept(xmit.GetSignal().RecvOffer()))
		}
	}(demo)

	xmit.Join()
}

func (demo *DemoServer) Accept(offer []byte) (answer []byte) {
	uuid := uuid.New()
	connIdArray := [len(uuid)]byte(uuid)
	connId := connIdArray[:]
	answer, err := demo.server.Accept(connId, offer, demo.callbacks)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	demo.ctrlStream, _ = demo.server.NewStream(connId, 1, xmit.Mode_kSemiReliableOrdered, xmit.Priority_kNormal, 0)
	demo.dataStream, _ = demo.server.NewStream(connId, 2, xmit.Mode_kSemiReliableOrdered, xmit.Priority_kNormal, 0)
	return answer
}
