package models

import "time"

type Finance struct {
    ID          string    `json:"id" db:"id"`
    UserID      int64     `json:"user_id" db:"user_id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Amount      int       `json:"amount" db:"amount"`
    Status      int16     `json:"status" db:"status"`
    VideoID     int64     `json:"video_id" db:"video_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PaymentSettings struct {
    ID                    int64     `json:"id" db:"id"`
    PaymentsEnabled       bool      `json:"payments_enabled" db:"payments_enabled"`
    VideoPaymentsEnabled  bool      `json:"video_payments_enabled" db:"video_payments_enabled"`
    AgentPaymentsEnabled  bool      `json:"agent_payments_enabled" db:"agent_payments_enabled"`
    PartnerBonusesEnabled bool      `json:"partner_bonuses_enabled" db:"partner_bonuses_enabled"`
    CreatedAt             time.Time `json:"created_at" db:"created_at"`
    UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}
