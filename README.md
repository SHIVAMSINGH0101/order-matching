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
    "restaurant_name": "Truffles",
    "restaurant_lat": 12.9620,
    "restaurant_lon": 77.6386,
    "customer_name": "Ananya Mehta",
    "customer_lat": 12.9652,
    "customer_lon": 77.6101,
    "prep_time_minutes": 15.0
}
```
-  For second order
```json
{
    "restaurant_name": "Empire Restaurant",
    "restaurant_lat": 12.9351,
    "restaurant_lon": 77.6250,
    "customer_name": "Rohit Sharma",
    "customer_lat": 12.9279,
    "customer_lon": 77.6271,
    "prep_time_minutes": 10.0
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
    "restaurant_name": "Truffles",
    "restaurant_lat": 12.9620,
    "restaurant_lon": 77.6386,
    "customer_name": "Ananya Mehta",
    "customer_lat": 12.9652,
    "customer_lon": 77.6101,
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
    "total_time_minutes": 52.178401910206894,
    "route": [
        {
            "step": "Empire Restaurant",
            "location_id": 5,
            "time_taken_minutes": 13.381407234709194
        },
        {
            "step": "Rohit Sharma",
            "location_id": 6,
            "time_taken_minutes": 2.496967359325224
        },
        {
            "step": "Truffles",
            "location_id": 7,
            "time_taken_minutes": 26.973886989705857
        },
        {
            "step": "Ananya Mehta",
            "location_id": 8,
            "time_taken_minutes": 9.326140326466621
        }
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
