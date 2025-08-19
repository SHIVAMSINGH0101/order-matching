# Order Routing API (Go + MySQL)

This assignment implements two APIs to create food delivery orders and compute the optimal delivery route for two simultaneous orders. Built with Go, Gorilla Mux, and MySQL using a clean architecture (handlers → services → repository → database).

## What was built

-   An order creation API that stores restaurant and customer locations and creates an order.
-   A best-route API that, given a delivery partner's location and exactly two order IDs, returns the optimal visiting sequence to minimize total time, including prep time at restaurants.

## Project Structure

```
go-demo/
├── cmd/api/                    # Application entry point
│   └── main.go
├── internal/
│   ├── config/                 # Env config loader
│   ├── database/               # DB connection
│   ├── models/                 # Entities (Location, Order)
│   ├── repository/             # Data access
│   ├── services/               # Business logic
│   └── handlers/               # HTTP handlers (order APIs)
├── database/
│   └── schema.sql              # MySQL schema
├── go.mod
└── go.sum
```

## Prerequisites

-   Go 1.24+
-   MySQL 9.4.0

## Database setup

Apply the schema to your MySQL instance:

```bash
mysql -u <DB_USER> -p -h <DB_HOST> -P <DB_PORT> < database/schema.sql
```

Schema creates:

-   `locations(id, name, latitude, longitude)`
-   `orders(orderId, resLocationId, cusLocationId, prepTimeInMinutes, createdAt, updatedAt)` with FKs to `locations`
-   Or we can directly run the SQL statements mentioned in schema.sql

## Configuration

Environment variables (defaults in parentheses from `internal/config/config.go`):

```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=wifiname
export DB_NAME=ordersdb
```

## Run the server

```bash
go mod tidy
go run cmd/api/main.go
```

Base URL: `http://localhost:${SERVER_PORT}` (default `8080`)
API prefix: `/api/v1`

---

## APIs (in `internal/handlers/order_handler.go`)

-   API's to run from Postman

### 1) Create Order

-   **Method**: POST
-   **Path**: `/api/v1/order/create`
-   **Description**: Creates restaurant and customer locations, then an order linking them.

Request headers:

-   `Content-Type: application/json`

Request body:

```json
{
	"restaurant_name": "Pizza Palace",
	"restaurant_lat": 40.7128,
	"restaurant_lon": -74.006,
	"customer_name": "John Doe",
	"customer_lat": 40.7589,
	"customer_lon": -73.9851,
	"prep_time_minutes": 15.0
}
```

Responses:

-   201 Created

```json
{ "status": "created", "orderId": 123 }
```

Curl example:

-   Curl which can be run from terminal

```bash
curl -X POST http://localhost:8080/api/v1/order/create \
  -H 'Content-Type: application/json' \
  -d '{
    "restaurant_name": "Pizza Palace",
    "restaurant_lat": 40.7128,
    "restaurant_lon": -74.006,
    "customer_name": "John Doe",
    "customer_lat": 40.7589,
    "customer_lon": -73.9851,
    "prep_time_minutes": 15.0
  }'
```

### 2) Get Best Route

-   **Method**: GET
-   **Path**: `/api/v1/order/best_route`
-   **Description**: Computes the optimal visiting sequence for exactly two orders given the rider's current location.

Query parameters:

-   `lat` (float, required): Rider latitude
-   `lon` (float, required): Rider longitude
-   `orderIds` (string, required): Comma separated two order IDs, e.g. `1,2`

Example request:

```
GET /api/v1/order/best_route?lat=40.7505&lon=-73.9934&orderIds=1,2
```

Response (200 OK):

```json
{
	"total_time_minutes": 45.5,
	"route": [
		{
			"step": "Pizza Palace",
			"location_id": 12,
			"time_taken_minutes": 10.2
		},
		{ "step": "John Doe", "location_id": 13, "time_taken_minutes": 8.3 },
		{
			"step": "Burger Barn",
			"location_id": 18,
			"time_taken_minutes": 12.1
		},
		{ "step": "Jane Roe", "location_id": 19, "time_taken_minutes": 14.9 }
	]
}
```

Curl example:

```bash
curl "http://localhost:8080/api/v1/order/best_route?lat=40.7505&lon=-73.9934&orderIds=1,2"
```

Notes on algorithm:

-   Only two orders are supported for routing.
-   Enumerates all valid sequences where each restaurant is visited before its customer.
-   Travel time uses a haversine distance approximation at a constant speed of 20 km/h.
-   Adds restaurant prep time to the leg that arrives at each restaurant.

### Health Check

-   This is a helth check API
-   **GET** `/health`

---

## Postman quickstart

1. Create a POST request to `http://localhost:8080/api/v1/order/create` with the JSON body above.
2. Create two orders and capture the returned `orderId` from each 201 response (e.g., `id1`, `id2`).
3. Send a GET request to `http://localhost:8080/api/v1/order/best_route?lat=40.7505&lon=-73.9934&orderIds=<id1>,<id2>`.
