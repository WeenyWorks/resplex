package resplex

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/WeenyWorks/resplex/lib/visheader"
	"github.com/spf13/cobra"
)

func handleVistorConn(conn net.Conn, cl *connLake) {
	defer conn.Close()
	l := make([]byte, 2)
	_, err := conn.Read(l)
	if err != nil {
		log.Fatal("failed to read length", err)
		return
	}

	visHeaderBinary := make([]byte, binary.LittleEndian.Uint16(l))
	n, err := conn.Read(visHeaderBinary)
	if err != nil {
		log.Fatal("failed to read vistor header", err)
		return
	}

	if n != int(binary.LittleEndian.Uint16(l)) {
		log.Fatal("short vistor header: ", n)
		return
	}

	log.Println(visHeaderBinary)
	vh := &visheader.VisHeader{}
	err = visheader.Unmarshal(visHeaderBinary, vh)
	if err != nil {
		log.Fatal("failed to parse header: ", err)
		return
	}

	stream, err := cl.GetStreamTo(vh.MachineID)
	if err != nil {
		log.Println("failed to get stream for ", vh.MachineID, ":", err)
		return
	}
	
	go func() {
		for {
			_, err := io.Copy(conn, stream)
			if err != nil {
				break
			}
		}
	}()
	for {
		io.Copy(stream, conn)
		if err != nil {
			break
		}
	}
	log.Println("finished handle conn: ", conn)
	stream.Close()
}

func entry(cmd *cobra.Command, args []string) {
	lnTCP, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9898,
	})
	if err != nil {
		panic(err)
	}
	defer lnTCP.Close()

	cl := NewConnLake()
	go cl.Serve("0.0.0.0:6007")
	
	for {
		conn, err := lnTCP.Accept()
		if err != nil {
			log.Fatal("failed to accept new conn due to: ", err)
			continue
		}

		go handleVistorConn(conn, cl)
	}
}

var ServeCMD = &cobra.Command{
	Use:        "serve",
	Short:      "run as a resplex server",
	Run: entry,
}
