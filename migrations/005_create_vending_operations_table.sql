-- Создание таблицы операций
CREATE TABLE IF NOT EXISTS vending_operations (
    id SERIAL PRIMARY KEY,
    vending_machine_id INTEGER NOT NULL REFERENCES vending_machines(id),
    operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('restock', 'collection', 'maintenance')),
    performed_by INTEGER NOT NULL REFERENCES users(id),
    operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    toys_before INTEGER,
    toys_after INTEGER,
    toys_added INTEGER DEFAULT 0,
    cash_before DECIMAL(10,2),
    cash_after DECIMAL(10,2),
    cash_collected DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_vending_operations_date ON vending_operations(operation_date);
CREATE INDEX IF NOT EXISTS idx_vending_operations_machine ON vending_operations(vending_machine_id);
CREATE INDEX IF NOT EXISTS idx_vending_operations_type ON vending_operations(operation_type);