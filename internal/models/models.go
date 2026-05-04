package models

type Case struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	VinNumber    string `json:"vin_number"`
	Approved     bool   `json:"approved"`
	Manufactured bool   `json:"manufactured"`
}

type VinDetails struct {
	Manufacturer string
	ModelYear    string
	Make         string
	Model        string
	VehicleType  string
	BodyClass    string
}
