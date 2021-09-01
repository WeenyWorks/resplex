package visheader
import (
	"github.com/WeenyWorks/resplex/lib/regheader"
	"github.com/kelindar/binary"
)

type VisHeader struct {
	MachineID regheader.MachineID
	Port uint16
}

func (v VisHeader) Marshal() ([]byte ,error) {
	res, err := binary.Marshal(v)
	if err != nil {
		return nil, err
	}
	length := make([]byte, 2)
	binary.LittleEndian.PutUint16(length, uint16(len(res)))
	return append(length, res...), nil
}

func Unmarshal(buf []byte, vh *VisHeader) error {
	err := binary.Unmarshal(buf[2:], vh)
	return err
}
