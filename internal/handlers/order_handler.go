package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	orderModel "github.com/SHIVAMSINGH0101/go-demo/internal/models"
	orderService "github.com/SHIVAMSINGH0101/go-demo/internal/services"
	"github.com/SHIVAMSINGH0101/go-demo/internal/utils"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	Service orderService.OrderService
}

func NewOrderHandler(service orderService.OrderService) *OrderHandler {
	return &OrderHandler{
		Service: service,
	}
}

func (h *OrderHandler) RegisterOrderHandlers(r *mux.Router) {
	r.HandleFunc("/order/create", h.CreateOrder).Methods("POST")
	r.HandleFunc("/order/best_route", h.GetBestRoute).Methods("GET")
}

/*
* CreateOrderRequest - This object stores Restaurant and Customer
* location information. The unique row id is taken as orderId.
*/
type CreateOrderRequest struct {
	RestaurantName string  `json:"restaurant_name"`
	RestaurantLat  float64 `json:"restaurant_lat"`
	RestaurantLon  float64 `json:"restaurant_lon"`
	CustomerName   string  `json:"customer_name"`
	CustomerLat    float64 `json:"customer_lat"`
	CustomerLon    float64 `json:"customer_lon"`
	PrepTimeMin    float64 `json:"prep_time_minutes"`
}

/*
* CreateOrder : This API creates Order in our DB
*/
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	restaurant := &orderModel.Location{
		Name:      req.RestaurantName,
		Latitude:  req.RestaurantLat,
		Longitude: req.RestaurantLon,
	}
	customer := &orderModel.Location{
		Name:      req.CustomerName,
		Latitude:  req.CustomerLat,
		Longitude: req.CustomerLon,
	}

	// Save restaurant location - In production restaurant location will be mapped to resId
	resID, err := h.Service.CreateLocation(restaurant)
	if err != nil {
		log.Printf("failed to save restaurant, err %+v", err)
		http.Error(w, "failed to save restaurant", http.StatusInternalServerError)
		return
	}

	// Save customer location - In production customer location will be mapped to cusId
	cusID, err := h.Service.CreateLocation(customer)
	if err != nil {
		log.Printf("failed to save customer, err %+v", err)
		http.Error(w, "failed to save customer", http.StatusInternalServerError)
		return
	}

	order := &orderModel.Order{
		ResLocationID: resID,
		CusLocationID:  cusID,
		PrepTimeInMinutes: req.PrepTimeMin,
	}
	orderId, err := h.Service.CreateOrder(order)
	if err != nil {
		log.Printf("failed to create order, err %+v", err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "created",
		"orderId": orderId,
	})
}

// GetBestRoute - Returns optimal path for the delivery partner
func (h *OrderHandler) GetBestRoute(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)

	if err1 != nil || err2 != nil {
		http.Error(w, "invalid coordinates", http.StatusBadRequest)
		return
	}

	stringOrderIDs := strings.Split(r.URL.Query().Get("orderIds"), ",")
	orderIDs := make([]int64, 0)

	for _, id := range stringOrderIDs {
		orderId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Printf("failed to covert orderId, err %+v", err)
			http.Error(w, `failed to covert orderId `, http.StatusBadRequest)
			return
		}
		orderIDs = append(orderIDs, orderId)
	}

	// Get Orders data
	orders, err := h.Service.GetOrdersByIDs(orderIDs)
	if err != nil {
		log.Printf("failed to fetch orders, err %+v", err)
		http.Error(w, "failed to fetch orders", http.StatusInternalServerError)
		return
	}

	if len(orders) != 2 {
		http.Error(w, "expected exactly 2 orders", http.StatusBadRequest)
		return
	}

	locIDs := []int64{
		int64(orders[0].ResLocationID), int64(orders[0].CusLocationID),
		int64(orders[1].ResLocationID), int64(orders[1].CusLocationID),
	}
	// Get locations data for the locationIds in the orders
	locations, err := h.Service.GetLocationsByIDs(locIDs)
	if err != nil {
		http.Error(w, "failed to fetch locations", http.StatusInternalServerError)
		return
	}

	// User - Delivery partners current location
	userLocation := orderModel.Location{
		Latitude: lat, 
		Longitude: lon,
	}

	// Returns the best possible route to cover all orders
	bestRoute := utils.GetBestRoute(userLocation, orders, locations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestRoute)
}
