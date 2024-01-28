package models

type JWTToken struct {
	PrivateKey  string `json:"privateKey"`
	ClientEmail string `json:"clientEmail"`
}

type gcs struct {
}
