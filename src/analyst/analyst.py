import zmq
import time
import math
import sys
import msg2a_pb2
import signal

context = zmq.Context()
socket = context.socket(zmq.PULL)
socket.bind("tcp://*:6666")

#   string timestamp = 1;
#   uint64 seq = 2;
#   bool   mismatch = 3;
#   uint64 msgtype = 4;
#   uint64 size =5;
seq_dict = {}
lost_seq = {}
LONG_dict = {}
seq_all = {}
max_rtt = 0
min_rtt = 100
i = 0
last_seq = -1

def handler(signum, frame):
        res = input("do U want to exit and print the log? y/n")
        if res == "yes":
            print(lost_seq, len(lost_seq)/last_seq)
            #sorted(seq_dict(), key=lambda item:item[0])
            for k in sorted(seq_dict.keys()):
                print("from", (k-1)*3, "to", k*3-1,seq_dict[k])
            for k in sorted(LONG_dict.keys()):
                print("from", (k-1)*100, "to", k*100-1,LONG_dict[k])
            exit(1)
            fo = open("1.txt", "w")
            fo.write(str(seq_all))
            fo.close()
        else:
            print("f**k")

signal.signal(signal.SIGINT, handler)

_msg = msg2a_pb2.Msg2A()

while True:
    try:
        #print("wait for client ...")
        message = socket.recv()
        #print("no wait")
        _msg.ParseFromString(message)
        if _msg.size < 500:
            print(lost_seq, len(lost_seq)/last_seq)
            #sorted(seq_dict(), key=lambda item:item[0])
            for k in sorted(seq_dict.keys()):
                    print("from", k * 3 - 1 , "to", (k + 1) * 3 - 1, seq_dict[k])
            for k in sorted(LONG_dict.keys()):
                print("from", k * 100 + 1, "to", (k + 1)*100-1, LONG_dict[k])
            break
        
        if _msg.mismatch :
            print("mistake")
            print(lost_seq)
            #sorted(seq_dict(), key=lambda item:item[0])
            for k in sorted(seq_dict.keys()):
                    print("from", (k-1)*3, "to", k*3 - 1, seq_dict[k])

            with open ("./test.txt", "w") as f:
                f.write(_msg)

                print("mistake")
                break
            
        RTT = (_msg.timestamprecv - _msg.timestampsend) / 1000000
        min_rtt = min(RTT, min_rtt)
        max_rtt = max(RTT, max_rtt)
        seq_all[_msg.seq] = RTT
        if RTT > 100:
            LONG = RTT // 100
            if LONG in LONG_dict.keys():
                LONG_dict[LONG] = LONG_dict[LONG] + 1
            else:
                LONG_dict[LONG] = 1
        else:
            SORT = (RTT - 0) // 3
            if SORT in seq_dict.keys():
                seq_dict[SORT] = seq_dict[SORT] + 1
            else:
                seq_dict[SORT] = 1
        if last_seq + 1 != _msg.seq:
            i = i + 1
            print(last_seq + 1)
            lost_seq[i] = last_seq + 1
        last_seq = _msg.seq
        print("message from client:", "RTT:", (_msg.timestamprecv - _msg.timestampsend) / 1000000, "Seq:", _msg.seq,  "msgType:", _msg.msgtype, "Size:", _msg.size)
    except Exception as e:
        print('exception:',e)
        sys.exit()

    
