package resplex

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"

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
	f, _ := os.Create("perfn")
	pprof.StartCPUProfile(f)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGALRM)
	go func() {
		<-sigs
		pprof.StopCPUProfile()
	}()
	log.Println("Waiting on ", visitAddr, " and ", proxyAddr)
	ip, port, err := parseAddr(visitAddr)
	if err != nil {
		return
	}
	lnTCP, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	defer lnTCP.Close()

	cl := NewConnLake()
	go cl.Serve(proxyAddr)

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

var visitAddr string
var proxyAddr string

func init() {
	ServeCMD.PersistentFlags().StringVarP(&visitAddr,
		"listenVisitor", "l", "0.0.0.0:9898",
		"address for visit proxied service")
	ServeCMD.PersistentFlags().StringVarP(&proxyAddr, "proxy",
		"p", "0.0.0.0:6007", "address for devices register")
}
