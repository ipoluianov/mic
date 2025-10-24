package mic

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

type MicapController struct {
	mtx               sync.Mutex
	Status            Status
	activeDevicePath  string
	activeDeviceIndex int

	devicePaths []string

	version string
}

func NewMicapController() *MicapController {
	var c MicapController
	c.devicePaths = []string{"/dev/uhid0", "/dev/uhid1", "/dev/uhid2", "/dev/uhid3"}
	return &c
}

func (c *MicapController) Start() {
	for _, path := range c.devicePaths {
		go c.ThReadContinuous(path)
	}
	go c.ThRequestStatus()
}

func (c *MicapController) ThRequestStatus() {
	for {
		fmt.Println("ThRequestStatus", time.Now(), "ActiveDevicePath:", c.GetActiveDevicePath(), "Version:", c.version)
		_, err := c.WriteToDevice(MakeRequestVersionFrame())
		if err != nil {
			fmt.Println("WriteToDevice error:", err)
		}

		time.Sleep(500 * time.Millisecond)

		_, err = c.WriteToDevice(MakeRequestSystemStatusFrame())
		if err != nil {
			fmt.Println("WriteToDevice error:", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (c *MicapController) GetActiveDevicePath() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.activeDevicePath != "" {
		return c.activeDevicePath
	}
	c.activeDeviceIndex++
	if c.activeDeviceIndex >= len(c.devicePaths) {
		c.activeDeviceIndex = 0
	}
	return c.devicePaths[c.activeDeviceIndex]
}

func (c *MicapController) ResetActiveDevice() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.activeDevicePath = ""
}

func (c *MicapController) WriteToDevice(data []byte) (int, error) {
	devPath := c.GetActiveDevicePath()

	f, err := os.OpenFile(devPath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		c.ResetActiveDevice()
		return 0, err
	}
	defer f.Close()

	fd := int(f.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Println("SetNonblock error:", err)
		c.ResetActiveDevice()
		return 0, err
	}

	packetSize := len(data)

	out := make([]byte, packetSize)
	copy(out, data)
	//fmt.Println("SND:", hex.EncodeToString(out))
	n, err := f.Write(out)
	if err != nil {
		fmt.Println("Write error:", err)
		c.ResetActiveDevice()
		return 0, err
	}

	return n, nil
}

func (c *MicapController) ThReadContinuous(devPath string) {
	in := make([]byte, 64)
	var err error
	var f *os.File

	for {
		if f == nil {
			f, err = os.OpenFile(devPath, os.O_RDONLY, 0)
			if err != nil {
				//fmt.Println("OpenFile error:", err)
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Println("Open file success - READ", devPath)
		}

		n, err := f.Read(in)
		if err == nil && n > 0 {
			//fmt.Println("RCV:", hex.EncodeToString(in))
			c.ParseFrame(in, devPath)
		} else {
			//fmt.Println("Read error:", err)
			f.Close()
			f = nil
		}
	}
}

func (c *MicapController) ParseFrame(data []byte, pathToDevice string) {
	if len(data) < 64 {
		fmt.Println("ParseFrame: data too short")
		return
	}

	cmd := binary.LittleEndian.Uint16(data[0:])
	processed := false

	if cmd == 1103 { // 0x044F
		for i := 0; i < 8; i++ {
			c.Status.ADC[i] = binary.LittleEndian.Uint16(data[20+i*2:])
		}
		processed = true
	}

	if cmd == 1102 { // 0x044E
		offset := 20
		c.Status.SYSTEM.TIMING.IsT0_Done = data[offset+0]
		c.Status.SYSTEM.TIMING.IsT1_Done = data[offset+1]
		c.Status.SYSTEM.TIMING.IsT2_Done = data[offset+2]
		c.Status.SYSTEM.TIMING.IsT3_Done = data[offset+3]
		c.Status.SYSTEM.TIMING.IsT4_Done = data[offset+4]
		c.Status.SYSTEM.TIMING.IsT5_Done = data[offset+5]
		c.Status.SYSTEM.TIMING.IsT6_Done = data[offset+6]
		c.Status.SYSTEM.TIMING.IsT7_Done = data[offset+7]
		c.Status.SYSTEM.TIMING.IsT8_Done = data[offset+8]
		c.Status.SYSTEM.TIMING.IsT9_Done = data[offset+9]
		c.Status.SYSTEM.FLAGS = binary.LittleEndian.Uint32(data[offset+12:])
		c.Status.SYSTEM.OPTICAL.Optical1 = binary.LittleEndian.Uint16(data[offset+16:])
		c.Status.SYSTEM.OPTICAL.Optical2 = binary.LittleEndian.Uint16(data[offset+18:])
		c.Status.SYSTEM.TEMPERATURE.Sensor1 = binary.LittleEndian.Uint16(data[offset+20:])
		c.Status.SYSTEM.TEMPERATURE.Sensor2 = binary.LittleEndian.Uint16(data[offset+22:])
		processed = true
	}

	if cmd == 1101 { // 0x044D
		// version response
		version := binary.LittleEndian.Uint16(data[20:])
		// fmt.Println("MICAP Version:", fmt.Sprintf("0x%04X", version))
		c.version = fmt.Sprintf("0x%04X", version)
		processed = true
	}

	if processed {
		c.activeDevicePath = pathToDevice
	}
}

/*func ReadFromDeviceWithTimeout(devPath string, packetSize int, timeout time.Duration) ([]byte, error) {
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
}*/

/*func FindMicapDevice() (filePath string, version string) {
	listOfDevices := []string{"/dev/uhid0", "/dev/uhid1", "/dev/uhid2", "/dev/uhid3", "/dev/uhid4", "/dev/uhid5"}
	for _, devPath := range listOfDevices {
		fmt.Println("try", devPath)
		_, err := os.Stat(devPath)
		if err == nil {
			fmt.Println("Device found:", devPath)
			WriteToDevice(devPath, MakeRequestVersionFrame())
			response, err := ReadFromDeviceWithTimeout(devPath, 64, 1*time.Second)
			if err == nil {
				v := binary.LittleEndian.Uint16(response[20:])
				fmt.Println("Recevied version:", fmt.Sprintf("0x%04X", v))
				return devPath, fmt.Sprintf("0x%04X", v)
			}
		}
	}
	return
}*/
