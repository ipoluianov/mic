package mic

import (
	"encoding/hex"
	"fmt"
	"os"
	"syscall"
	"time"
)

func WriteToDevice(devPath string, data []byte) (int, error) {
	f, err := os.OpenFile(devPath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return 0, err
	}
	defer f.Close()
	fmt.Println("Open file success")

	fd := int(f.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Println("SetNonblock error:", err)
		return 0, err
	}

	packetSize := len(data)

	out := make([]byte, packetSize)
	copy(out, data)
	fmt.Println("Write:", hex.EncodeToString(out))
	n, err := f.Write(out)
	if err != nil {
		fmt.Println("Write error:", err)
		return 0, err
	}
	fmt.Printf("Sent %d bytes\n", n)

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
				return
			}
			fmt.Println("Open file success")
		}

		fmt.Println("Reading ...")
		n, err := f.Read(in)
		if err == nil && n > 0 {
			fmt.Println("Read1:", hex.EncodeToString(in))
		} else {
			fmt.Println("Read error:", err)
			f.Close()
			f = nil
		}
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
