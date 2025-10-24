package mic

type Status struct {
	ADC    [8]uint16
	SYSTEM SystemStatus
}

type SystemStatus struct {
	TIMING struct {
		IsT0_Done uint8
		IsT1_Done uint8
		IsT2_Done uint8
		IsT3_Done uint8
		IsT4_Done uint8
		IsT5_Done uint8
		IsT6_Done uint8
		IsT7_Done uint8
		IsT8_Done uint8
		IsT9_Done uint8
	}

	FLAGS uint32

	OPTICAL struct {
		Optical1 uint16
		Optical2 uint16
	}
	TEMPERATURE struct {
		Sensor1 uint16
		Sensor2 uint16
	}
}

/*
#pragma pack(push, 1)
typedef struct
{
	struct
	{
		UINT8 isT0_Done : 1; // Flag which indicates whether T0 action has been executed
		UINT8 isT1_Done : 1; // Flag which indicates whether T1 action has been executed
		UINT8 isT2_Done : 1; // Flag which indicates whether T2 action has been executed
		UINT8 isT3_Done : 1; // Flag which indicates whether T3 action has been executed
		UINT8 isT4_Done : 1; // Flag which indicates whether T4 action has been executed
		UINT8 isT5_Done : 1; // Flag which indicates whether T5 action has been executed
		UINT8 isT6_Done : 1; // Flag which indicates whether T6 action has been executed
		UINT8 isT7_Done : 1; // Flag which indicates whether T7 action has been executed
		UINT8 isT8_Done : 1; // Flag which indicates whether T8 action has been executed
		UINT8 isT9_Done : 1; // Flag which indicates whether T8 action has been executed
		UINT8 : 6;			 // 6 empty fields
		UINT8 : 8;			 // 8 empty fields
		UINT8 : 8;			 // 8 empty fields
	} TIMING;
	union
	{
		UINT32 Flags;
		struct
		{
			UINT8 isMFCOn : 1;			// This variable determines if mass flow controllers are powered
			UINT8 isMFCPurge : 1;		// This variable indicates the MFC is in purge mode
			UINT8 isMagFanON : 1;		// This variable determines if magnetron fan is currently powered
			UINT8 isTeslaCoilON : 1;	// This variable determines if the tesla coil is enabled
			UINT8 isArgSolenoidON : 1;	// This variable determines if argon solenoid is currently powered
			UINT8 isPumpON : 1;			// This variable determines if peristaltic pump is spinning
			UINT8 isPumpCW : 1;			// This variable determines if peristaltic pump is spinning clock-wise
			UINT8 isPumpMAXRPM : 1;		// This variable determines if peristaltic pump is spinning at max speed
			UINT8 isPlasmaOn : 1;		// This variable determines whether plasma is present
			UINT8 isMagArcing : 1;		// This variable determines whether arcing happened
			UINT8 isSens1OverTemp : 1;	// This variable determines if temperature sensor 1 has overheated
			UINT8 isSens2OverTemp : 1;	// This variable determines if temperature sensor 2 has overheated
			UINT8 isMWAvailable : 1;	// This variable indicates that MW energy may be available ******************
			UINT8 isIgnitionFail : 1;	// This variable  indicates that auto-ignition attempt has failed *******
			UINT8 isTTLInterlock1 : 1;	// This variable determines if there is a TTL interlock
			UINT8 isTTLInterlock2 : 1;	// This variable determines if there is a TTL interlock
			UINT8 isSMPSOn : 1;			// This variable determines if HV power supply is currently powered
			UINT8 isSMPSError : 1;		// This variable indicates that some sort of error occurred while attempting to communicate with SMPS
			UINT8 isCoolFlowError : 1;	// This flag indicates that actual coolant flow rate (feedback from MFC), is 2LPM below desired flow rate
			UINT8 isAuxFlowError : 1;	// This flag indicates that actual auxiliary flow rate (feedback from MFC), is 200mLPM below desired flow rate
			UINT8 isNebFlowError : 1;	// This flag indicates that actual nebulizer flow rate (feedback from MFC), is 200mLPM below desired flow rate
			UINT8 isCoolMonitoring : 1; // This flag indicates whether the system is monitoring coolant feedback system (some delay is introduced because of MFC lag)
			UINT8 isAuxMonitoring : 1;	// This flag indicates whether the system is monitoring coolant feedback system (some delay is introduced because of MFC lag)
			UINT8 isNebMonitoring : 1;	// This flag indicates whether the system is monitoring coolant feedback system (some delay is introduced because of MFC lag)
			UINT8 : 8;					// Empty fields
		} FLAGS;
	};
	struct
	{
		UINT16 Optical1; // Arcing detect sensor
		UINT16 Optical2; // Plasma detect sensor
	} OPTICAL;
	struct
	{
		UINT16 Sensor1; // Temperature sensor 1
		UINT16 Sensor2; // Temperature sensor 2
	} TEMPERATURE;
} SYSTEM_STATUS;

#pragma pack(pop)

*/
