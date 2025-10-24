package main

import (
	"fmt"
	"time"

	"github.com/ipoluianov/mic/mic"
)

func main() {
	fmt.Println("Started")
	filePath, version := mic.FindMicapDevice()
	fmt.Printf("Found MICAP device: %s, version: %s\n", filePath, version)
	return

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

		{
			fmt.Println("IsT0_Done:", mic.STATUS.SYSTEM.TIMING.IsT0_Done)
			fmt.Println("IsT1_Done:", mic.STATUS.SYSTEM.TIMING.IsT1_Done)
			fmt.Println("IsT2_Done:", mic.STATUS.SYSTEM.TIMING.IsT2_Done)
			fmt.Println("IsT3_Done:", mic.STATUS.SYSTEM.TIMING.IsT3_Done)
			fmt.Println("IsT4_Done:", mic.STATUS.SYSTEM.TIMING.IsT4_Done)
			fmt.Println("IsT5_Done:", mic.STATUS.SYSTEM.TIMING.IsT5_Done)
			fmt.Println("IsT6_Done:", mic.STATUS.SYSTEM.TIMING.IsT6_Done)
			fmt.Println("IsT7_Done:", mic.STATUS.SYSTEM.TIMING.IsT7_Done)
			fmt.Println("IsT8_Done:", mic.STATUS.SYSTEM.TIMING.IsT8_Done)
			fmt.Println("IsT9_Done:", mic.STATUS.SYSTEM.TIMING.IsT9_Done)

			fmt.Println("FLAGS:", fmt.Sprintf("0x%08X", mic.STATUS.SYSTEM.FLAGS))
			fmt.Println("Optical1:", mic.STATUS.SYSTEM.OPTICAL.Optical1)
			fmt.Println("Optical2:", mic.STATUS.SYSTEM.OPTICAL.Optical2)
			fmt.Println("Sensor1:", mic.STATUS.SYSTEM.TEMPERATURE.Sensor1)
			fmt.Println("Sensor2:", mic.STATUS.SYSTEM.TEMPERATURE.Sensor2)
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
