package main

import (
	"fmt"
	"time"

	"github.com/ipoluianov/mic/mic"
)

func main() {
	fmt.Println("Started")

	go mic.ThReadContinuous("/dev/uhid1")

	for i := 0; i < 10000; i++ {
		fmt.Println("Iteration", i)
		_, err := mic.WriteToDevice("/dev/uhid1", mic.MakeRequestSystemStatusFrame())
		if err != nil {
			fmt.Println("WriteToDevice error:", err)
		}

		fmt.Println("Status:")
		for j := 0; j < 8; j++ {
			fmt.Printf(" ADC[%d] = %d\n", j, mic.STATUS.ADC[j])
		}
		time.Sleep(1 * time.Second)
	}

	/*resp, err := mic.ReadFromDeviceWithTimeout("/dev/uhid1", 64, 2*time.Second)
	if err != nil {
		fmt.Println("ReadFromDeviceWithTimeout error:", err)
		return
	}

	fmt.Println("Response", hex.EncodeToString(resp))

	fmt.Println("Finished")*/
}
