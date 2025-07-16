package smartevse

import (
	"fmt"
)

type (
	ChargeModeId int16
	LBMode       int16
	StateId      uint8
	ErrorId      uint8

	Wifi struct {
		Status        string `json:"status"`
		SSID          string `json:"ssid"`
		RSSI          int    `json:"rssi" chargeflux:"wifi_rssi"`
		BSSID         string `json:"bssid"`
		AutoConnect   bool   `json:"auto_connect"`
		AutoReconnect bool   `json:"auto_reconnect"`
	}

	EVSE struct {
		Temp      int  `json:"temp" chargeflux:"evse_temp"`
		TempMax   int  `json:"temp_max"`
		Connected bool `json:"connected"`
		Access    int  `json:"access"`
		// Mode is ignored, it's in SmartEVSESettings
		LBMode         LBMode `json:"loadbl"`
		PWM            uint16 `json:"pwm" chargeflux:"evse_pwm"`
		SolarStopTimer uint16 `json:"solar_stop_timer"`
		// JSON field "state" is ignored in favor of "state_id"
		State StateId `json:"state_id" chargeflux:"evse_state"`
		// JSON field "error" is ignored in favor of "error_id"
		Error ErrorId `json:"error_id" chargeflux:"evse_error"`
	}

	Settings struct {
		ChargeCurrent int `json:"charge_current" chargflux:"charge_current"`
	}

	PhaseCurrents struct {
		Total      int16 `json:"TOTAL"`
		L1         int16 `json:"L1" chargeflux:"phase_current_l1"`
		L2         int16 `json:"L2" chargeflux:"phase_current_l2"`
		L3         int16 `json:"L3" chargeflux:"phase_current_l3"`
		ChargingL1 bool  `json:"charging_L1" chargeflux:"phase_charging_l1"`
		ChargingL2 bool  `json:"charging_L2" chargeflux:"phase_charging_l2"`
		ChargingL3 bool  `json:"charging_L3" chargeflux:"phase_charging_l3"`
	}

	SmartEVSESettings struct {
		Version string `json:"version"`
		// JSON field "mode" is ignored in favor of "mode_id"
		Mode          ChargeModeId  `json:"mode_id" chargeflux:"charge_mode"`
		CarConnected  bool          `json:"car_connected" chargeflux:"car_connected"`
		Wifi          Wifi          `json:"wifi"`
		EVSE          EVSE          `json:"evse"`
		PhaseCurrents PhaseCurrents `json:"phase_currents"`
		Settings      Settings      `json:"settings"`
	}
)

const (
	ChargeModeNormal ChargeModeId = iota + 1
	ChargeModeSolar
	ChargeModeSmart
)

func (c ChargeModeId) String() string {
	switch c {
	case ChargeModeNormal:
		return "normal"
	case ChargeModeSolar:
		return "solar"
	case ChargeModeSmart:
		return "smart"
	default:
		return "unknwon"
	}
}

const (
	LBModeDisabled LBMode = iota
	LBModeMaster
	LBModeNode1
	LBModeNode2
	LBModeNode3
	LBModeNode4
	LBModeNode5
	LBModeNode6
	LBModeNode7
)

func (c LBMode) String() string {
	if c == LBModeDisabled {
		return "disabled"
	}
	if c == LBModeMaster {
		return "master"
	}
	if c >= LBModeNode1 && c <= LBModeNode7 {
		return fmt.Sprintf("node%d", c-1)
	}
	return "unknown"
}

const (
	StateA StateId = iota
	StateB
	StateC
	StateD
	StateCommB
	StateCommBOK
	StateCommC
	StateCommCOK
	StateActStart
	StateB1
	StateC1
	StateModemRequest
	StateModemWait
	StateModemDone
	StateModemDenied
	StateNone = 255
)

func (s StateId) String() string {
	switch s {
	case StateA:
		return "vehicle not connected"
	case StateB:
		return "vehicle connected, not ready to accept energy"
	case StateC:
		return "vehicle connected, ready to accept energy, ventilation not required"
	case StateD:
		return "vehicle connected, ready to accept energy, ventilation required"
	case StateCommB:
		return "state change request (A->B)"
	case StateCommBOK:
		return "state change (A->B) OK"
	case StateCommC:
		return "state change request (B->C)"
	case StateCommCOK:
		return "state change (B->C) OK"
	case StateActStart:
		return "state change (B->C) OK"
	case StateB1:
		return "vehicle connected, no PWM signal"
	case StateC1:
		return "vehicle charging, no PWM signal"
	case StateModemRequest:
		return "vehicle connected, requesting ISO15118 communication at 0% duty cycle"
	case StateModemWait:
		return "vehicle connected, requesting ISO15118 communication at 5% duty cycle"
	case StateModemDone:
		return "vehicle connected, ISO15118 communication done"
	case StateModemDenied:
		return "vehicle connected, ISO15118 communication denied"
	default:
		return "unknown"
	}
}

const (
	ErrorNone ErrorId = iota
	ErrorNoPower
	ErrorComm
	ErrorTempHigh
	ErrorUnused
	ErrorRCMTripped
	ErrorWaitingForSolar
	ErrorTestIO
	ErrorFlash
)

func (e ErrorId) String() string {
	switch e {
	case ErrorNone:
		return "no error"
	case ErrorNoPower:
		return "no power available"
	case ErrorComm:
		return "communication error"
	case ErrorTempHigh:
		return "temperature high"
	case ErrorRCMTripped:
		return "RCM tripped"
	case ErrorWaitingForSolar:
		return "waiting for solar"
	case ErrorFlash:
		return "flash error"
	default:
		return "unknown error"
	}
}
