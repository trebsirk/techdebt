package entropy

type FileEntropy struct {
	Filename string  `json:"filename"`
	Entropy  float64 `json:"entropy"`
}
