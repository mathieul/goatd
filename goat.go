package main

import (
  zmq "github.com/alecthomas/gozmq"
  "flag"
)

const (
  CommandEcho = "echo"
)

var command string
var message string

func init() {
  flag.StringVar(&command, "cmd", CommandEcho, "test connection by sending an echo request")
  flag.StringVar(&message, "msg", "testing...", "message to send")
  flag.Parse()
}

func main() {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  socket.Connect("tcp://127.0.0.1:5000")

  switch command {
  case CommandEcho:
    socket.Send([]byte(message), 0)
    println("Sending:", message)
    socket.Recv(0)
  }
}