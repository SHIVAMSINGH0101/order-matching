package utils

import (
	"math"

	"github.com/SHIVAMSINGH0101/go-demo/internal/models"
)

// RouteStep - Each step taken in the optimal approach
type RouteStep struct {
	Step       string  `json:"step"`
	LocationID int     `json:"location_id"`
	TimeTaken  float64 `json:"time_taken_minutes"`
}

// BestRouteResponse - Optimal steps for delivery partner to take
type BestRouteResponse struct {
	TotalTime float64     `json:"total_time_minutes"`
	Route     []RouteStep `json:"route"`
}

func GetBestRoute(
	userLocation models.Location,
	orders []models.Order, 
	locations []models.Location,
) BestRouteResponse {
	// We want to find the best route starting from userLocation, visiting both restaurants and customers in some order.
	// The route must visit each restaurant before its customer.

	// Possible sequences of indices representing the order of visits:
	// We have two orders: O1 and O2.
	// Constraints:
	// - O1.Restaurant before O1.Customer
	// - O2.Restaurant before O2.Customer
	// So the valid permutations of [R1, R2, C1, C2] with constraints are:
	// R1 R2 C1 C2
	// R1 R2 C2 C1
	// R2 R1 C1 C2 
	// R2 R1 C2 C1
	// R1 C1 R2 C2
	// R2 C2 R1 C1

	bestTime := math.MaxFloat64
	// Stores the best route sequence from all sequences
	bestRoute := BestRouteResponse{}
	// Get all possible location's sequences that we can try
	allPosibleSequences := getAllLocationSequences(orders)

	locationMap := make(map[int]models.Location)
	for _, location := range locations {
		locationMap[location.ID] = location
	}

	// Loop all possible route sequences
	for _, currSeq := range allPosibleSequences {
		currLocation := userLocation
		totalTime := 0.0
		routeTaken := make([]RouteStep, 0)

		for _, nextLocationId := range currSeq {
			nextloc := locationMap[int(nextLocationId)]
			timeTaken := getTravelTimeInMinutes(currLocation, nextloc)

			if nextLocationId == orders[0].ResLocationID {
				timeTaken += orders[0].PrepTimeInMinutes
			}
			if nextLocationId == orders[1].ResLocationID {
				timeTaken += orders[1].PrepTimeInMinutes
			}

			totalTime += timeTaken
			currLocation = nextloc
			
			routeTaken = append(routeTaken, RouteStep{
				Step: nextloc.Name,
				LocationID: nextloc.ID,
				TimeTaken: timeTaken,
			})
		}
		
		// Update best route sequence
		if totalTime < bestTime {
			bestTime = totalTime
			bestRoute = BestRouteResponse{
				TotalTime: bestTime,
				Route: routeTaken,
			} 
		}
	}

	return bestRoute
}

func getAllLocationSequences(orders []models.Order) [][]int64 {
	r1 := orders[0].ResLocationID
	r2 := orders[1].ResLocationID
	c1 := orders[0].CusLocationID
	c2 := orders[1].CusLocationID

	return [][]int64{
		{r1, r2, c1, c2},
		{r1, r2, c2, c1},
		{r2, r1, c1, c2},
		{r2, r1, c2, c1},
		{r1, c1, r2, c2},
		{r2, c2, r1, c1},
	}
}

// Returns approax time taken to reach from -> to location
func getTravelTimeInMinutes(from, to models.Location) float64 {
	const speed = 20.0 // Speed is in km/h
	dist := haversine(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
	return (dist / speed) * 60
}

// Users haversine formula to get approax distance between from -> to lat and lon
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusInKM = 6371
	dLat := (lat2 - lat1) * (3.14159 / 180)
	dLon := (lon2 - lon1) * (3.14159 / 180)
	lat1 *= (3.14159 / 180)
	lat2 *= (3.14159 / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadiusInKM * c
}