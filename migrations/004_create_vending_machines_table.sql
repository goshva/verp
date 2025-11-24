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

-- Insert sample machines
INSERT INTO vending_machines (serial_number, model, location_id, capacity_toys, current_toys_count, cash_amount, installation_date) VALUES
('VM001', 'KidyFun Pro 100', 1, 100, 45, 1250.50, '2024-01-15'),
('VM002', 'ToyMagic 200', 1, 150, 120, 890.75, '2024-02-20'),
('VM003', 'HappyKids Plus', 2, 100, 30, 567.25, '2024-03-10'),
('VM004', 'SuperToy Deluxe', 3, 200, 180, 2100.00, '2024-01-05');