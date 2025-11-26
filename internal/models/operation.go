package models

import "time"

type Operation struct {
    ID              int       `json:"id"`
    MachineID       int       `json:"vending_machine_id"`
    MachineName     string    `json:"machine_name,omitempty"`
    OperationType   string    `json:"operation_type"`
    PerformedBy     int       `json:"performed_by"`
    PerformerName   string    `json:"performer_name"`
    OperationDate   time.Time `json:"operation_date"`
    ToysBefore      int       `json:"toys_before"`
    ToysAfter       int       `json:"toys_after"`
    ToysAdded       int       `json:"toys_added"`
    CashBefore      float64   `json:"cash_before"`
    CashAfter       float64   `json:"cash_after"`
    CashCollected   float64   `json:"cash_collected"`
    RoutePointID    *int      `json:"route_point_id,omitempty"` // Сделать опциональным
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}