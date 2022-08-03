import zmq
import time
import math
import sys
import msg_pb2
import os
import mmap
import contextlib

def accept_data2u(mode: int, ip: str):
    socket_server = zmq.Context().socket(zmq.REQ)
    socket_server.connect("tcp://" + ip + ":5558")
    _msg2server = msg_pb2.Msg()
    _msg2server.sequence = mode
    _msg2server.payload = ip
    serial = _msg2server.SerializeToString() 
    socket_server.send(serial)
    print("send msg")
    mmsg = socket_server.recv()
    print(mmsg)

def accept_data_from_driver(param_dict: dict):

    _msg = msg_pb2.Msg()
    ub_dict = param_dict
    for case_id in ub_dict:
        server_mode = ub_dict[case_id]['server_mode']
        server_ip = ub_dict[case_id]['server_ip']

        print("Connecting to server...")  
        accept_data2u(server_mode, server_ip)
        print("finish")
        context = zmq.Context()
        socket = context.socket(zmq.PUSH)
        socket.connect ("tcp://127.0.0.1:5555")

        socket2 = zmq.Context().socket(zmq.PUSH)
        socket2.connect("tcp://127.0.0.1:7777")

        for user_id in ub_dict[case_id]:
            user_behavior = ub_dict[case_id][user_id]
            interval_seq = user_behavior['interval_seq']
            size_seq = user_behavior['size_seq']
            user_ip = user_behavior['user_ip']
            disconnect_time_seq = user_behavior['disconnect_time_seq']

            for idx, interval in enumerate(interval_seq):
                interval_ms = interval * (10 ** -3)
                current_size = size_seq[idx]
                #stop = socket2.recv()
                #if stop :
                #    print("1")
                if os.path.exists("test.dat"):
                    break
                time.sleep(interval_ms)
                current_delay = 0
                if current_size == -1:
                    current_delay = disconnect_time_seq[0]
                    disconnect_time_seq = disconnect_time_seq[1:]
                    _msg.sequence = 0

                _msg.payload = user_ip
                #print("socket") 
                serial = _msg.SerializeToString() 
                socket.send(serial)
                time.sleep(current_delay * (10 ** -2))
                #print("send")
                #message = socket.recv()
            
                # python 如何处理报错？！
                #print(message) 
                # send_data_2_go_server(current_size, user_ip)
            
# def send_data_2_go_server(size: int, user_ip: str):

#     _msg.sequence = size
#     _msg.payload = user_ip
#     print("socket") 
#     serial = _msg.SerializeToString() 
#     socket.send(serial)
#     message = socket.recv()
#     print(message)  

# message = socket.recv()  
# print("Received reply: ", message)
