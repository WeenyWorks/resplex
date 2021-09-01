package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"

	"github.com/WeenyWorks/resplex/lib/visheader"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9898")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	vh := &visheader.VisHeader{
		MachineID: "TESTMACHINEIDXCJ",
		Port:      22,
	}
	vhs, err := vh.Marshal()
	if err != nil {
		panic(err)
	}
	lenBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBuf, uint16(len(vhs)))
	n, err := conn.Write(lenBuf)
	if n != 2 {
		panic("failed to write len")
	}
	log.Println(vhs)
	n, err = conn.Write(vhs)
	stdio := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	go func() {
		for {
			_, err := io.Copy(conn, stdio)
			if err != nil {
				break
			}
		}
	}()
	for {
		io.Copy(stdio, conn)
		if err != nil {
			break
		}
	}
}
