package regheader

import(
	"github.com/kelindar/binary"
)

type MachineID string

type RegHeader struct {
	ID MachineID
	Token []byte
	Ports uint16
}


func (r RegHeader) Marshal() ([]byte, error) {
	res, err := binary.Marshal(r)
	if err != nil {
		return nil, err
	}
	length := make([]byte, 2)
	binary.LittleEndian.PutUint16(length, uint16(len(res)))
	return append(length, res...), nil
}

func Unmarshal(buf []byte, rh *RegHeader) error {
	err := binary.Unmarshal(buf, rh)
	return err
}
