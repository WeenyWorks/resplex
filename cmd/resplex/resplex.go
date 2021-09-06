package resplex

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/WeenyWorks/resplex/lib/visheader"
	"github.com/spf13/cobra"
)

type server struct {
	cl *connLake
	registerAddr string
	proxyAddr string
}

func newServer(ra string, pa string) *server {
	return &server{
		cl:           NewConnLake(),
		registerAddr: ra,
		proxyAddr:    pa,
	}
}

func (s *server)serve() {
	proxyListener, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		log.Fatalln("Failed to listen proxy address: ",
			proxyAddr ,  err)
	}
	defer proxyListener.Close()

	go s.cl.Serve(registerAddr)
	for {
		conn, err := proxyListener.Accept()
		if err != nil {
			log.Println("failed to accept new connection:",
				err)
			continue
		}
		go handleVistorConn(conn, s.cl)
	}
}

func handleVistorConn(conn net.Conn, cl *connLake) {
	defer conn.Close()
	l := make([]byte, 2)
	_, err := conn.Read(l)
	if err != nil {
		log.Println("failed to read length", err)
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
	defer stream.Close()
	go func() {
		_, err := io.Copy(conn, stream)
		if err != nil {
			log.Println("failed to copy: ", err)
		}
	}()
	_, err = io.Copy(stream, conn)
	if err != nil {
		log.Println("failed to copy: ", err)
		return
	}
	log.Println("finished handle conn: ", conn)
}

func entry(cmd *cobra.Command, args []string) {
	s := newServer(registerAddr, proxyAddr)
	s.serve()
}

var ServeCMD = &cobra.Command{
	Use:   "serve",
	Short: "run as a resplex server",
	Run:   entry,
}

func parseAddr(address string) (ip string, port int, err error) {
	sli := strings.Split(address, ":")
	if len(sli) != 2 {
		return "", 0, errors.New("Invalid listen address")
	}
	port, err = strconv.Atoi(sli[1])
	if err != nil {
		return "", 0, errors.New("Invalid port number")
	}
	return sli[0], port, nil
}

var proxyAddr string
var registerAddr string

func init() {
	ServeCMD.PersistentFlags().StringVarP(&proxyAddr,
		"listenProxy", "l", "0.0.0.0:9898",
		"address for visit proxied service")
	ServeCMD.PersistentFlags().StringVarP(&registerAddr, "register",
		"r", "0.0.0.0:6007", "address for devices register")
}
