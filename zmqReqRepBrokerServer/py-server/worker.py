import os
import zmq
import random
import time
# https://www.digitalocean.com/community/tutorials/how-to-work-with-the-zeromq-messaging-library

# ZeroMQ Context
context = zmq.Context()

# Define the socket using the "Context"
sock = context.socket(zmq.REP)
sock.connect("tcp://127.0.0.1:5560")

pid = os.getpid()

# Run a simple "Echo" server
while True:
    message = sock.recv()
    delay = random.randint(0, 3)
    print "Delay " + str(delay) + " seconds" + "Recv: " + message
    time.sleep(delay)
    sock.send("Python server [" + str(pid) + "] echo: " + message)
    #print "Echo: " + message
