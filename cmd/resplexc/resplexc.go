package resplexc

import (
	"io"
	"log"
	"net"

	"github.com/WeenyWorks/resplex/lib/regheader"
	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
	"github.com/spf13/cobra"
)

func handler(stream *smux.Stream) {
	conn, err := net.Dial("tcp4", "127.0.0.1:22")
	if err != nil {
		log.Println("failed to connect to sshd: ", err)
		return
	}
	go func() {
		for {
			_, err := io.Copy(stream, conn)
			if err != nil {
				break
			}
		}
	}()
	func() {
		for {
			_, err := io.Copy(conn, stream)
			if err != nil {
				break
			}
		}
	}()
}

func entry(clmd *cobra.Command, args []string) {
	log.Println("Starting...")
	conn, err := kcp.DialWithOptions("127.0.0.1:6007", nil, 10, 3)
	if err != nil {
		log.Println("failed to connect to server: ")
		panic("FIXME")
	}
	header := &regheader.RegHeader{
		ID:    "TESTMACHINEIDXCJ",
		Ports: []uint16{22},
	}
	buf, err := header.Marshal()
	if err != nil {
		log.Println("falied to marshal regheader ", err)
	}
	log.Println(buf)
	_, err = conn.Write(buf)
	if err != nil {
		log.Println("failed to send header: ", err)
		panic("FIXME")
	}
	session, err := smux.Server(conn, nil)
	if err != nil {
		log.Println("failed to create session: ", err)
		panic("FIXME")
	}
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Println("accept stream failed: ", err)
			continue
		}
		go handler(stream)
	}
}

var RegisterCMD = &cobra.Command{
	Use:        "connect",
	Short:      "run as a resplex client",
	Run: entry,
}
