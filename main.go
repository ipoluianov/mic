package main

import (
	"fmt"
	"time"

	"github.com/ipoluianov/mic/mic"
)

func main() {
	fmt.Println("Started")

	go mic.ThReadContinuous("/dev/uhid1")

	for {
		_, err := mic.WriteToDevice("/dev/uhid1", mic.MakeRequestADCFrame())
		if err != nil {
			fmt.Println("WriteToDevice error:", err)
			return
		}

		//fmt.Println("Waiting")
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
