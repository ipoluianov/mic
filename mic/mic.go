package mic

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"
	"time"
)

var STATUS Status

func WriteToDevice(devPath string, data []byte) (int, error) {
	f, err := os.OpenFile(devPath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return 0, err
	}
	defer f.Close()

	fd := int(f.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Println("SetNonblock error:", err)
		return 0, err
	}

	packetSize := len(data)

	out := make([]byte, packetSize)
	copy(out, data)
	fmt.Println("SND:", hex.EncodeToString(out))
	n, err := f.Write(out)
	if err != nil {
		fmt.Println("Write error:", err)
		return 0, err
	}
	//fmt.Printf("Sent %d bytes\n", n)

	return n, nil
}

func ThReadContinuous(devPath string) {
	in := make([]byte, 64)
	var err error
	var f *os.File
	// Читаем ответ с таймаутом
	//timeout := 1 * time.Second

	for {
		if f == nil {
			f, err = os.OpenFile(devPath, os.O_RDONLY, 0)
			if err != nil {
				fmt.Println("OpenFile error:", err)
				continue
			}
			fmt.Println("Open file success - READ")
		}

		//fmt.Println("Reading ...")
		n, err := f.Read(in)
		if err == nil && n > 0 {
			fmt.Println("RCV:", hex.EncodeToString(in))
			ParseFrame(in)
		} else {
			fmt.Println("Read error:", err)
			f.Close()
			f = nil
		}
	}
}

func ParseFrame(data []byte) {
	if len(data) < 64 {
		fmt.Println("ParseFrame: data too short")
		return
	}

	cmd := binary.LittleEndian.Uint16(data[0:])
	if cmd == 1103 { // 0x044F
		for i := 0; i < 8; i++ {
			STATUS.ADC[i] = binary.LittleEndian.Uint16(data[20+i*2:])
		}
	}

	if cmd == 1102 { // 0x044E
		offset := 20
		STATUS.SYSTEM.TIMING.IsT0_Done = data[offset+0]
		STATUS.SYSTEM.TIMING.IsT1_Done = data[offset+1]
		STATUS.SYSTEM.TIMING.IsT2_Done = data[offset+2]
		STATUS.SYSTEM.TIMING.IsT3_Done = data[offset+3]
		STATUS.SYSTEM.TIMING.IsT4_Done = data[offset+4]
		STATUS.SYSTEM.TIMING.IsT5_Done = data[offset+5]
		STATUS.SYSTEM.TIMING.IsT6_Done = data[offset+6]
		STATUS.SYSTEM.TIMING.IsT7_Done = data[offset+7]
		STATUS.SYSTEM.TIMING.IsT8_Done = data[offset+8]
		STATUS.SYSTEM.TIMING.IsT9_Done = data[offset+9]
		STATUS.SYSTEM.FLAGS = binary.LittleEndian.Uint32(data[offset+12:])
		STATUS.SYSTEM.OPTICAL.Optical1 = binary.LittleEndian.Uint16(data[offset+16:])
		STATUS.SYSTEM.OPTICAL.Optical2 = binary.LittleEndian.Uint16(data[offset+18:])
		STATUS.SYSTEM.TEMPERATURE.Sensor1 = binary.LittleEndian.Uint16(data[offset+20:])
		STATUS.SYSTEM.TEMPERATURE.Sensor2 = binary.LittleEndian.Uint16(data[offset+22:])
	}

}

func ReadFromDeviceWithTimeout(devPath string, packetSize int, timeout time.Duration) ([]byte, error) {
	f, err := os.OpenFile(devPath, os.O_RDWR, 0)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return nil, err
	}
	defer f.Close()
	fmt.Println("Open file success")
	// Читаем ответ с таймаутом
	in := make([]byte, packetSize)
	//timeout := 1 * time.Second

	deadline := time.Now().Add(timeout)
	for {
		n, err := f.Read(in)
		if err == nil && n > 0 {
			fmt.Printf("Read %d bytes:\n", n)
			for i := 0; i < n; i++ {
				fmt.Printf("%02X ", in[i])
			}
			fmt.Println()
			break
		}

		if time.Now().After(deadline) {
			fmt.Println("Read timeout")
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return in, nil
}

func FindMicapDevice() (filePath string, version string) {
	listOfDevices := []string{"/dev/uhid0", "/dev/uhid1", "/dev/uhid2", "/dev/uhid3", "/dev/uhid4", "/dev/uhid5"}

	for _, devPath := range listOfDevices {
		fmt.Println("try", devPath)
		_, err := os.Stat(devPath)
		if err == nil {
			fmt.Println("Device found:", devPath)
			WriteToDevice(devPath, MakeRequestVersionFrame())
			response, err := ReadFromDeviceWithTimeout(devPath, 64, 1*time.Second)
			if err == nil {
				v := binary.LittleEndian.Uint16(response[40:])
				fmt.Println("Recevied version:", fmt.Sprintf("0x%04X", v))
			}
		}
	}
	return
}
