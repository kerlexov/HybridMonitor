package power

import "encoding/json"

func UnmarshalChassisPower(data []byte) (ChassisPower, error) {
	var r ChassisPower
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ChassisPower) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ChassisPower struct {
	OdataContext       string            `json:"@odata.context"`
	OdataID            string            `json:"@odata.id"`
	OdataType          string            `json:"@odata.type"`
	ID                 string            `json:"Id"`
	Name               string            `json:"Name"`
	OEM                ChassisPowerOEM   `json:"Oem"`
	PowerCapacityWatts int64             `json:"PowerCapacityWatts"`
	PowerConsumedWatts int64             `json:"PowerConsumedWatts"`
	PowerControl       []PowerControl    `json:"PowerControl"`
	PowerLimit         PowerLimit        `json:"PowerLimit"`
	PowerMetrics       PowerMetrics      `json:"PowerMetrics"`
	PowerSupplies      []PowerSupply     `json:"PowerSupplies"`
	Redundancy         []Redundancy      `json:"Redundancy"`
	Type               string            `json:"Type"`
	Links              ChassisPowerLinks `json:"links"`
}

type ChassisPowerLinks struct {
	Self Self `json:"self"`
}

type Self struct {
	Href string `json:"href"`
}

type ChassisPowerOEM struct {
	HP PurpleHP `json:"Hp"`
}

type PurpleHP struct {
	OdataType               string                  `json:"@odata.type"`
	SNMPPowerThresholdAlert SNMPPowerThresholdAlert `json:"SNMPPowerThresholdAlert"`
	Type                    string                  `json:"Type"`
	Links                   HPLinks                 `json:"links"`
}

type HPLinks struct {
	FastPowerMeter        Self `json:"FastPowerMeter"`
	FederatedGroupCapping Self `json:"FederatedGroupCapping"`
	PowerMeter            Self `json:"PowerMeter"`
}

type SNMPPowerThresholdAlert struct {
	DurationInMin  int64  `json:"DurationInMin"`
	ThresholdWatts int64  `json:"ThresholdWatts"`
	Trigger        string `json:"Trigger"`
}

type PowerControl struct {
	PowerCapacityWatts int64        `json:"PowerCapacityWatts"`
	PowerConsumedWatts int64        `json:"PowerConsumedWatts"`
	PowerLimit         PowerLimit   `json:"PowerLimit"`
	PowerMetrics       PowerMetrics `json:"PowerMetrics"`
}

type PowerLimit struct {
	LimitInWatts interface{} `json:"LimitInWatts"`
}

type PowerMetrics struct {
	AverageConsumedWatts int64 `json:"AverageConsumedWatts"`
	IntervalInMin        int64 `json:"IntervalInMin"`
	MaxConsumedWatts     int64 `json:"MaxConsumedWatts"`
	MinConsumedWatts     int64 `json:"MinConsumedWatts"`
}

type PowerSupply struct {
	OEM                  PowerSupplyOEM `json:"Oem"`
	Status               StatusClass    `json:"Status"`
	FirmwareVersion      *string        `json:"FirmwareVersion,omitempty"`
	LastPowerOutputWatts *int64         `json:"LastPowerOutputWatts,omitempty"`
	LineInputVoltage     *int64         `json:"LineInputVoltage,omitempty"`
	LineInputVoltageType *string        `json:"LineInputVoltageType,omitempty"`
	Model                *string        `json:"Model,omitempty"`
	Name                 *string        `json:"Name,omitempty"`
	PowerCapacityWatts   *int64         `json:"PowerCapacityWatts,omitempty"`
	PowerSupplyType      *string        `json:"PowerSupplyType,omitempty"`
	SerialNumber         *string        `json:"SerialNumber,omitempty"`
	SparePartNumber      *string        `json:"SparePartNumber,omitempty"`
}

type PowerSupplyOEM struct {
	HP FluffyHP `json:"Hp"`
}

type FluffyHP struct {
	OdataType               string  `json:"@odata.type"`
	BayNumber               int64   `json:"BayNumber"`
	Type                    string  `json:"Type"`
	AveragePowerOutputWatts *int64  `json:"AveragePowerOutputWatts,omitempty"`
	HotplugCapable          *bool   `json:"HotplugCapable,omitempty"`
	MaxPowerOutputWatts     *int64  `json:"MaxPowerOutputWatts,omitempty"`
	Mismatched              *bool   `json:"Mismatched,omitempty"`
	PowerSupplyStatus       *Status `json:"PowerSupplyStatus,omitempty"`
	IPDU                    *IPDU   `json:"iPDU,omitempty"`
	IPDUCapable             *bool   `json:"iPDUCapable,omitempty"`
}

type IPDU struct {
	ID           string `json:"Id"`
	Model        string `json:"Model"`
	SerialNumber string `json:"SerialNumber"`
	IPDUStatus   Status `json:"iPDUStatus"`
}

type Status struct {
	State string `json:"State"`
}

type StatusClass struct {
	Health string `json:"Health"`
	State  string `json:"State"`
}

type Redundancy struct {
	MaxNumSupported int64           `json:"MaxNumSupported"`
	MemberID        string          `json:"MemberId"`
	MinNumNeeded    int64           `json:"MinNumNeeded"`
	Name            string          `json:"Name"`
	RedundancySet   []RedundancySet `json:"RedundancySet"`
}

type RedundancySet struct {
	OdataID string `json:"@odata.id"`
}
