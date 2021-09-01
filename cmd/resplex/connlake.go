package main

import (
	"encoding/binary"
	"errors"
	"log"
	"sync"

	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
	"github.com/WeenyWorks/resplex/lib/regheader"
)

type machine struct {
	ports []uint16
	session smux.Session
}

type connLake struct {
	lock sync.RWMutex
	lake map[regheader.MachineID]machine
}

func NewConnLake() *connLake{
	return &connLake{
		lock: sync.RWMutex{},
		lake: map[regheader.MachineID]machine{},
	}
}

func (c *connLake) SetMachine(m regheader.MachineID, session smux.Session, ports []uint16) {
	newM := machine{
		ports: ports,
		session:  session,
	}
	if origM, exist := c.lake[m]; exist {
		log.Println("replacing ", m, "from ", origM, "to ", newM)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lake[m] = newM
}

func (c *connLake) GetStreamTo(machineid regheader.MachineID) (*smux.Stream, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if machine, exist := c.lake[machineid]; exist {
		return machine.session.OpenStream()
	} else {
		return nil, errors.New("no such machine registered")
	}
}

func (c *connLake) GetSessionFor(machineid regheader.MachineID) (*smux.Session, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if machine, exist := c.lake[machineid]; exist {
		return &machine.session, errors.New("no such machine registered.")
	} else {
		return nil, errors.New("no such machine regstered")
	}
}

func (c *connLake) Listen(addr string) (*kcp.Listener, error) {
	return  kcp.ListenWithOptions(addr, nil, 10, 3)
}

func (c *connLake) _handleConn(s *kcp.UDPSession) {
	headerLengthBinary := make([]byte, 2)
	n, err := s.Read(headerLengthBinary)
	headerLength := binary.LittleEndian.Uint16(headerLengthBinary)
	if err != nil  || n < 2 {
		return
	}
	headerBuf := make([]byte, headerLength)
	n, err = s.Read(headerBuf)
	if err != nil {
		log.Fatal("failed to read register header: ", err)
		return
	}
	if n != int(headerLength) {
		log.Fatal("short register header: ", err)
		return
	}

	regh := &regheader.RegHeader{}
	log.Println(headerBuf)
	err = regheader.Unmarshal(headerBuf, regh)
	if err != nil {
		log.Fatal("invalid header: ", err)
		return
	}
	smuxSession, err := smux.Client(s, nil)
	if err != nil {
		log.Fatal("failed to start smux stream: ", err)
		return
	}
	c.SetMachine(regh.ID, *smuxSession, regh.Ports)
}


func (c *connLake) Serve(addr string) {
	listener, err := c.Listen(addr)
	if err != nil {
		panic(err)
	}
	for {
		s, err := listener.AcceptKCP()
		if err != nil {
			log.Fatal("failed to connect: ", err)
			continue
		}

		go c._handleConn(s)
	}
}

