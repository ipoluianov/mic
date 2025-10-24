package mic

import "encoding/binary"

type MicFrame struct {
	MICAPCommand     uint16
	MICAPStatus      uint16
	OutTransferCount uint16
	InTransferCount  uint16
	Address          uint16
	Var_1            uint16
	Var_2            uint16
	Var_3            uint16
	Var_4            uint16
	Var_5            uint16
}

func (c *MicFrame) ToBytes() []byte {
	frame := make([]byte, 64)
	binary.LittleEndian.PutUint16(frame[0:], c.MICAPCommand)
	binary.LittleEndian.PutUint16(frame[2:], c.MICAPStatus)
	binary.LittleEndian.PutUint16(frame[4:], c.OutTransferCount)
	binary.LittleEndian.PutUint16(frame[6:], c.InTransferCount)
	binary.LittleEndian.PutUint16(frame[8:], c.Address)
	binary.LittleEndian.PutUint16(frame[10:], c.Var_1)
	binary.LittleEndian.PutUint16(frame[12:], c.Var_2)
	binary.LittleEndian.PutUint16(frame[14:], c.Var_3)
	binary.LittleEndian.PutUint16(frame[16:], c.Var_4)
	binary.LittleEndian.PutUint16(frame[18:], c.Var_5)
	return frame
}

func MakeRequestADCFrame() []byte {
	var fr MicFrame
	fr.MICAPCommand = 1103
	fr.InTransferCount = 16
	return fr.ToBytes()
}
