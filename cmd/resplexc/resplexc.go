package resplexc

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/WeenyWorks/resplex/lib/regheader"
	"github.com/spf13/cobra"
	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
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

func entry(cmd *cobra.Command, args []string) {
	log.Println("Starting...")
	conn, err := kcp.DialWithOptions("127.0.0.1:6007", nil, 10, 3)
	if err != nil {
		for {
			log.Println("failed to connect to server: ")
			time.Sleep(10 * time.Second)
			conn, err = kcp.DialWithOptions("127.0.0.1:6007", nil, 10, 3)
			if err == nil {
				break
			}
		}
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
		for {
			log.Println("failed to create session: ", err)
			session, err = smux.Server(conn, nil)
			if err == nil {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			if err == io.ErrClosedPipe ||
			   err == net.ErrClosed {
				for {
					log.Println("lose connection ", err, "reconnecting")
					conn, err = kcp.DialWithOptions("127.0.0.1:6007", nil, 10, 3)
					if err != nil {
						for {
							log.Println("failed to connect to server: ")
							time.Sleep(10 * time.Second)
							conn, err = kcp.DialWithOptions("127.0.0.1:6007", nil, 10, 3)
							if err == nil {
								break
							}
						}
					}

					session, err = smux.Server(conn, nil)
					if err == nil {
						break
					}
					time.Sleep(10 * time.Second)
				}
			}
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
