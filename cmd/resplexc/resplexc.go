package resplexc

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/WeenyWorks/resplex/lib/regheader"
	"github.com/spf13/cobra"
	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
)

func (c client)handler(stream *smux.Stream) {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", "127.0.0.1", c.sshPort))
	if err != nil {
		log.Println("failed to connect to sshd: ", err)
		return
	}
	go func() {
		_, err := io.Copy(stream, conn)
		if err != nil {
			// FIXME: need to reconnect
			return
		}
	}()
	func() {
		_, err := io.Copy(conn, stream)
		if err != nil {
			// FIXME: need to reconnect
			return
		}
	}()
}

type client struct {
	machineid regheader.MachineID
	hubAddress string
	sshPort uint16
	kcpSession *kcp.UDPSession
	smuxSession *smux.Session
	_headerBuf []byte
}

func newClient(machineid string, hubAddr string, sshPort uint16) (*client, error) {
	c := &client{
		machineid:  regheader.MachineID(machineid),
		hubAddress: hubAddr,
		sshPort: sshPort,
	}
	err := c.genHeader()
	return c, err
}

func (c *client)genHeader()error{
	header := &regheader.RegHeader{
		ID:    c.machineid,
		Ports: sshPort,
	}
	buf, err := header.Marshal()
	if err != nil {
		return err
	}
	c._headerBuf = buf
	return nil
}

func (c *client)connectHub() error {
	kcpSession, err := kcp.DialWithOptions(c.hubAddress, nil, 10, 3)
	if err != nil {
		return err
	}
	c.kcpSession = kcpSession
	n, err := kcpSession.Write(c._headerBuf)
	if err != nil {
		kcpSession.Close()
		return err
	}
	if n != len(c._headerBuf) {
		kcpSession.Close()
		return err
	}
	smuxSession, err := smux.Server(kcpSession, nil)
	if err != nil {
		kcpSession.Close()
		return err
	}
	c.smuxSession = smuxSession
	return nil
}

func (c *client)serve() error {
	for {
		stream, err := c.smuxSession.AcceptStream()
		if err != nil {
			return err
		}
		go c.handler(stream)
	}
}

func entry(cmd *cobra.Command, args []string) {
	c, err := newClient(machineID, hubAddr, sshPort)
	if err != nil {
		log.Panic(err)
	}
	for {
		err = c.connectHub()
		if err != nil {
			continue
		}

		err = c.serve()
		if err != nil {
			continue
		}
		c.smuxSession.Close()
		c.kcpSession.Close()
		time.Sleep(5 * time.Second)
	}
}

var (
	machineID string
	hubAddr string
	sshPort uint16
)

var RegisterCMD = &cobra.Command{
	Use:        "register",
	Short:      "register local service(port) to hub.",
	Run: entry,
}

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	RegisterCMD.PersistentFlags().StringVarP(&machineID, "machineid", "m", hostname, "unique machine id")
	RegisterCMD.PersistentFlags().StringVarP(&hubAddr, "hubaddr", "a", "127.0.0.1:6007", "unique machine id")
	RegisterCMD.PersistentFlags().Uint16VarP(&sshPort, "sshport", "s", 22, "the ssh port that you want to expose")
}
