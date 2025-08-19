package repository

import (
	"database/sql"
	"errors"
	"strings"

	routeModels "github.com/SHIVAMSINGH0101/go-demo/internal/models"
)

// Order repository interacts with orders and locations table
type OrderRepository interface {
	InsertLocation(loc *routeModels.Location) (int64, error)
	GetLocationByID(id int64) (*routeModels.Location, error)
	GetLocationsByIDs(ids []int64) ([]routeModels.Location, error)

	InsertOrder(order *routeModels.Order) (int64, error)
	GetOrderByID(id int64) (*routeModels.Order, error)
	GetOrdersByIDs(ids []int64) ([]routeModels.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository{
	return &orderRepository{
		db: db,
	}
}

// CRUD operations on Location and Order 
func (r *orderRepository) InsertLocation(loc *routeModels.Location) (int64, error) {
	query := `INSERT INTO locations 
			(name, latitude, longitude)
			VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, loc.Name, loc.Latitude, loc.Longitude)
	if err != nil {
		return 0, errors.New("failed to insert location")
	}

	return result.LastInsertId()
}

func (r *orderRepository) GetLocationByID(id int64) (*routeModels.Location, error) {
	query := `SELECT id, name, latitude, longitude 
			  FROM locations
			  WHERE id = ?`

	var loc routeModels.Location
	err := r.db.QueryRow(query, id).Scan(&loc.ID, &loc.Name, &loc.Latitude, &loc.Longitude)
	if err != nil {
		return nil, errors.New("location not found")
	}

	return &loc, nil
}

func (r *orderRepository) GetLocationsByIDs(locationIds []int64) ([]routeModels.Location, error) {
	if len(locationIds) == 0 {
        return nil, nil
    }

	placeHolder := strings.Repeat("?,", len(locationIds) - 1) + "?"
	query := `SELECT id, name, latitude, longitude 
				FROM locations
				WHERE id IN (` + placeHolder + `)`
	
	args := make([]interface{}, len(locationIds))
    for i, id := range locationIds {
        args[i] = id
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var locs []routeModels.Location
    for rows.Next() {
        var loc routeModels.Location
        if err := rows.Scan(&loc.ID, &loc.Name, &loc.Latitude, &loc.Longitude); err != nil {
            return nil, err
        }
        locs = append(locs, loc)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return locs, nil
}

func (r *orderRepository) InsertOrder(order *routeModels.Order) (int64, error) {
	query := `INSERT INTO orders 
		(resLocationId, cusLocationId, prepTimeInMinutes)
		VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, order.ResLocationID, order.CusLocationID, order.PrepTimeInMinutes)
	if err != nil {
		return 0, errors.New("failed to insert order")
	}

	return result.LastInsertId()
}

func (r *orderRepository) GetOrderByID(orderId int64) (*routeModels.Order, error) {
	query := `SELECT orderId, resLocationId, cusLocationId, prepTimeInMinutes 
		FROM orders
		WHERE orderId = ?`

	var order routeModels.Order

	err := r.db.QueryRow(query, orderId).Scan(
		&order.CreatedAt,
		&order.ResLocationID,
		&order.CusLocationID,
		&order.PrepTimeInMinutes,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetOrdersByIDs(orderIds []int64) ([]routeModels.Order, error) {
	if len(orderIds) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(orderIds)-1) + "?"
	query := `SELECT orderId, resLocationId, cusLocationId, prepTimeInMinutes 
				FROM orders
				WHERE orderId IN (` + placeholders + `)`

	args := make([]interface{}, len(orderIds))
    for i, id := range orderIds {
        args[i] = id
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []routeModels.Order
    for rows.Next() {
        var order routeModels.Order
        if err := rows.Scan(&order.OrderID, &order.ResLocationID, &order.CusLocationID, &order.PrepTimeInMinutes); err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return orders, nil
}