package models

import (
    "time"
)

type VendingOperation struct {
    ID               int64     `json:"id" db:"id"`
    VendingMachineID int64     `json:"vending_machine_id" db:"vending_machine_id"`
    MachineSerial    string    `json:"machine_serial" db:"machine_serial"` // Added for display
    OperationType    string    `json:"operation_type" db:"operation_type"`
    PerformedBy      int64     `json:"performed_by" db:"performed_by"`
    PerformerName    string    `json:"performer_name" db:"performer_name"` // Added for display
    OperationDate    time.Time `json:"operation_date" db:"operation_date"`
    ToysBefore       int       `json:"toys_before" db:"toys_before"`
    ToysAfter        int       `json:"toys_after" db:"toys_after"`
    ToysAdded        int       `json:"toys_added" db:"toys_added"`
    CashBefore       float64   `json:"cash_before" db:"cash_before"`
    CashAfter        float64   `json:"cash_after" db:"cash_after"`
    CashCollected    float64   `json:"cash_collected" db:"cash_collected"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}