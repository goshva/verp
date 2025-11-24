-- Migration: 003_create_locations_table.sql
CREATE TABLE IF NOT EXISTS locations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    contact_person VARCHAR(255),
    contact_phone VARCHAR(50),
    monthly_rent DECIMAL(10,2) DEFAULT 0,
    rent_due_day INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_locations_active ON locations(is_active);
CREATE INDEX IF NOT EXISTS idx_locations_name ON locations(name);

-- Insert sample locations
INSERT INTO locations (name, address, contact_person, contact_phone, monthly_rent, rent_due_day) VALUES
('ТЦ "Москва"', 'ул. Ленина, 1', 'Иванов Иван', '+7-999-123-45-67', 15000.00, 15),
('ТРК "Европа"', 'пр. Мира, 25', 'Петрова Мария', '+7-999-765-43-21', 20000.00, 10),
('Аэропорт', 'ш. Аэропортовское, 10', 'Сидоров Алексей', '+7-999-555-44-33', 30000.00, 5);