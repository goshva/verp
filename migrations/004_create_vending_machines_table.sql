-- Migration: 004_create_vending_machines_table.sql
CREATE TABLE IF NOT EXISTS vending_machines (
    id BIGSERIAL PRIMARY KEY,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    model VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    location_id BIGINT,
    capacity_toys INTEGER DEFAULT 100,
    current_toys_count INTEGER DEFAULT 0,
    cash_amount DECIMAL(10,2) DEFAULT 0,
    last_maintenance_date DATE,
    next_maintenance_date DATE,
    installation_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_machines_serial ON vending_machines(serial_number);
CREATE INDEX IF NOT EXISTS idx_machines_location ON vending_machines(location_id);
CREATE INDEX IF NOT EXISTS idx_machines_status ON vending_machines(status);