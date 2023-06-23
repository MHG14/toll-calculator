package types

type Invoice struct {
	OBUID         int     `json:"obuId"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuId"`
	Unix  int64   `json:"unix"`
}

type OBUData struct {
	OBUID int     `json:"obuId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
