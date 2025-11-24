package models

import (
    "time"
)

type VendingMachine struct {
    ID                   int64     `json:"id" db:"id"`
    SerialNumber         string    `json:"serial_number" db:"serial_number"`
    LocationID           int64     `json:"location_id" db:"location_id"`
    LocationName         string    `json:"location_name" db:"location_name"` // Add this
    Model                string    `json:"model" db:"model"`
    CapacityToys         int       `json:"capacity_toys" db:"capacity_toys"`
    CurrentToysCount     int       `json:"current_toys_count" db:"current_toys_count"`
    CashAmount           float64   `json:"cash_amount" db:"cash_amount"`
    LastMaintenanceDate  time.Time `json:"last_maintenance_date" db:"last_maintenance_date"`
    NextMaintenanceDate  time.Time `json:"next_maintenance_date" db:"next_maintenance_date"`
    InstallationDate     time.Time `json:"installation_date" db:"installation_date"`
    Status               string    `json:"status" db:"status"`
    CreatedAt            time.Time `json:"created_at" db:"created_at"`
    UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}
