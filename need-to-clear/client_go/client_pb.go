package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	client, _ := zmq.NewSocket(zmq.REQ)
	client.Connect("tcp://127.0.0.1:5555")

	_msg := &Msg{
		Sequence: 5,
		Payload:  "ctm",
	}
	out, err := proto.Marshal(_msg)
	if err != nil {
		log.Fatalln("FAILED:", err)
	}
	fmt.Println(out)
	client.SendBytes(out, zmq.SNDMORE)

}
