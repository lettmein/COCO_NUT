package dto

import "time"

type Customer struct {
	CompanyName string `json:"company_name"`
	INN         string `json:"inn"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}

type Cargo struct {
	Name              string  `json:"name"`
	Quantity          int     `json:"quantity"`
	Weight            float64 `json:"weight"`
	Volume            float64 `json:"volume"`
	SpecialRequirements string  `json:"special_requirements"`
}

type Recipient struct {
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
}

type CreateRequestDTO struct {
	LogisticPointID int       `json:"logistic_point_id"`
	Customer        Customer  `json:"customer"`
	Cargo           Cargo     `json:"cargo"`
	Recipient       Recipient `json:"recipient"`
}

type UpdateRequestDTO struct {
	Status string `json:"status"`
}

type RequestResponse struct {
	ID              string    `json:"id"`
	LogisticPointID int       `json:"logistic_point_id"`
	Customer        Customer  `json:"customer"`
	Cargo           Cargo     `json:"cargo"`
	Recipient       Recipient `json:"recipient"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
