package models

import (
    "time"
)

type Location struct {
    ID            int64     `json:"id" db:"id"`
    Name          string    `json:"name" db:"name"`
    Address       string    `json:"address" db:"address"`
    ContactPerson string    `json:"contact_person" db:"contact_person"`
    ContactPhone  string    `json:"contact_phone" db:"contact_phone"`
    MonthlyRent   float64   `json:"monthly_rent" db:"monthly_rent"`
    RentDueDay    int       `json:"rent_due_day" db:"rent_due_day"`
    IsActive      bool      `json:"is_active" db:"is_active"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}