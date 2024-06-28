package domain

type FigaroData struct {
	Price              string `json:"price"`
	City               string `json:"city"`
	EnergeticPerf      string `json:"dpe"`
	GreenhouseGasIndex string `json:"ges"`
	Bedrooms           string `json:"bedrooms"`
	Bathrooms          string `json:"bathrooms"`
	Rooms              string `json:"rooms"`
	Area               string `json:"area"`
	Url                string `json:"url"`
	Source             string `json:"source"`
	Ref                string `json:"ref"`
	Region             string `json:"region"`
	Department         string `json:"department"`
	Type               int    `json:"type"`
}