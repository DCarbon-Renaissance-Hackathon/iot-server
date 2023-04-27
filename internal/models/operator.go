package models

type OpIotStatus struct {
	Id      int64      `json:"id,omitempty"`      // Iot id
	Address EthAddress `json:"address,omitempty"` // Iot address
	Status  OpStatus   `json:"status,omitempty"`  // Operator status
	Latest  int64      `json:"latest,omitempty"`  // Last update
}

type OpSensorMetric struct {
	Id     int64      `json:"id,omitempty"`
	Type   SensorType `json:"type,omitempty"`
	Metric *AllMetric `json:"metric,omitempty"`
	Latest int64      `json:"latest,omitempty"`
}
