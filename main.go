package main

import (
	"encoding/binary"
	"net"
	"strings"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var (
	host  = flag.String("host", "localhost", "hostname to listen to")
	token = flag.String("token", "myToken", "token to access mini monkey")
	port  = flag.String("port", "1773", "which port to listen to")
	room  = flag.String("room", "stats", "which room to use")
)

func compile(args []string) time.Duration {
	gcc := "/usr/bin/g++"

	t0 := time.Now()
	cmd := exec.Command(gcc, args...)
	t1 := time.Now()

	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return t1.Sub(t0)
}

func check(err error, t string) {
	if err != nil {
		fmt.Println(t)
		fmt.Println(err.Error())
		panic(err)
	}
}

func response(conn net.Conn) {
	buf := make([]byte, 1024)
	conn.Read(buf)
	fmt.Println("got: ", string(buf))
}

func sendStats(args []string, duration time.Duration) {
	l, err := net.Dial("tcp", *host+":"+*port)
	check(err, "connection error")
	defer l.Close()

	fmt.Println("connected to " + *host + " on port " + *port)

	_, err = l.Write(auth(*token))
	check(err, "token error")
	response(l)

	err = binary.Write(l, binary.LittleEndian, enter(*room))
	check(err, "room error")
	response(l)

	msg := strings.Join(args, " ")
	err = binary.Write(l, binary.LittleEndian, publish(msg))
	check(err, "unable to publish args")
	response(l)

	err = binary.Write(l, binary.LittleEndian, publishDuration(duration))
	check(err, "unable to publish duration")
	response(l)
}

func main() {
	flag.Parse()

	args := os.Args[1:]
	duration := compile(args)
	sendStats(args, duration)
}
