package visit

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/WeenyWorks/resplex/lib/regheader"
	"github.com/WeenyWorks/resplex/lib/visheader"
	"github.com/spf13/cobra"
	"github.com/xtaci/smux"
)

type vistor struct{
	proxyAddr string
	proxyPort string
	targetMachine string
	targetPort int
	conn net.Conn
	stream *smux.Stream
}

func newVistor(pa string, pp string, tm string, tp int) *vistor {
	return &vistor{
		proxyAddr:     pa,
		proxyPort:     pp,
		targetMachine: tm,
		targetPort:    tp,
	}
}

// visitor is a one-shot command thus we shouldn't retry any connection.
func (v *vistor)connectProxy() error {
	conn, err := net.Dial("tcp", fmt.Sprint(v.proxyAddr, ":", v.proxyPort))
	v.conn = conn
	return err
}

func(v *vistor)connectMachine() error {
	vh := &visheader.VisHeader{
		MachineID: regheader.MachineID(v.targetMachine),
		Port:      uint16(v.targetPort),
	}
	vhs, err := vh.Marshal()
	if err != nil {
		panic(err)
	}
	lenBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBuf, uint16(len(vhs)))
	n, err := v.conn.Write(lenBuf)
	if n != 2 {
		panic("failed to write len")
	}
	log.Println(vhs)
	n, err = v.conn.Write(vhs)
	if err != nil {
		return err
	}
	if n != len(vhs) {
		return errors.New("Error: short write when writing visheader")
	}
	return nil
}

func(v vistor)serve() error {
	stdio := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	go func(){
		_, err := io.Copy(v.conn, stdio)
		if err != nil {
			return
		}
	}()
	_, err := io.Copy(stdio, v.conn)
	return err

}

var targetMachine string
var targetPort int
var proxyAddr string
var proxyPort string

func entry(cmd *cobra.Command, args []string) {
	v := newVistor(proxyAddr, proxyPort, targetMachine, targetPort)
	v.connectProxy()
	v.connectMachine()
	v.serve()
}

var VisitCMD = &cobra.Command{
	Use: "visit",
	Short: "visit machine",
	Run: entry,
}

func init() {
	VisitCMD.PersistentFlags().StringVarP(&targetMachine, "machine", "m", "", "The machine ID that you want to visit")
	VisitCMD.PersistentFlags().IntVarP(&targetPort, "service", "s", 22, "The target service that you want to visit.")
	VisitCMD.PersistentFlags().StringVarP(&proxyAddr, "proxyaddr", "x", "127.0.0.1", "The proxy service address")
	VisitCMD.PersistentFlags().StringVarP(&proxyPort, "proxyport", "p", "9898", "The proxy service port")
}
