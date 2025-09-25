package types

type Routes struct {
	Legs       []Legs   `json:"legs"`
	WeightName string   `json:"weight_name"`
	Geometry   Geometry `json:"geometry"`
	Weight     float64  `json:"weight"`
	Duration   float64  `json:"duration"`
	Distance   float64  `json:"distance"`
}
