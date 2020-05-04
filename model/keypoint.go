package model

import "math"

type Keypoint struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

// Gap 求两个坐标之间的坐标差
func (v *Keypoint) Gap(input *Keypoint) float64 {
	return math.Abs(v.Y-input.Y) + math.Abs(v.X-input.X)
}
