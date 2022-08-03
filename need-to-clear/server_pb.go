package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	server, _ := zmq.NewSocket(zmq.REP)
	server.Bind("tcp://127.0.0.1:5555")

	_msg, err := server.RecvBytes(zmq.DONTWAIT)

	newTest := &Msg{}
	err = proto.Unmarshal(_msg, newTest)
	if err != nil {
		log.Fatal("unmarshaling error", err)
	}
	fmt.Println(newTest.Payload, newTest.Sequence)
}
