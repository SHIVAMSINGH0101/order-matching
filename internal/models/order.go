package models

import "time"

// Location - Stores location coordinates for actors - Restaurant and Customer
type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Order - Stores customer's order info - restaurant and customer locationId
type Order struct {
	OrderID int `json:"orderId"`
	ResLocationID int64 `json:"resLocationId"`
	CusLocationID int64 `json:"cusLocationId"`
	PrepTimeInMinutes float64 `json:"prepTimeInMinutes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VehicleModel struct {
	MakeId int64 `json:"Make_ID"`
	MakeName string `json:"Make_Name"`
	ModelId int64 `json:"Model_ID"`
	ModelName string `json:"Model_Name"`
}

// VeichlesSold 

type VPICResponse struct {
	Count int `json:"Count"`
	Message string `json:"Message"`
	SearchCriteria string `json:"SearchCriteria"`
	Results []VehicleModel `json:"Results"`
}
