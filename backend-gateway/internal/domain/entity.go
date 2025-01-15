package domain

type Task struct {
	ID                string  `json:"id"`
	Distance          float64 `json:"distance"`
	Weather           string  `json:"weather"`
	TrafficLevel      string  `json:"traffic_level"`
	TimeOfDay         string  `json:"time_of_day"`
	VehicleType       string  `json:"vehicle_type"`
	PreparationTime   int     `json:"preparation_time"`
	CourierExperience float64 `json:"courier_experience"`
}
