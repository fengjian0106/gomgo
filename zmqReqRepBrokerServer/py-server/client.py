import zmq
import os
import time

# ZeroMQ Context
context = zmq.Context()

# Define the socket using the "Context"
sock = context.socket(zmq.REQ)
sock.connect("tcp://127.0.0.1:5559")

pid = os.getpid()

while True:
    # Send a "message" using the socket
    sock.send("Hello from python client [" + str(pid) + "]")
    print sock.recv()
    time.sleep(1)
