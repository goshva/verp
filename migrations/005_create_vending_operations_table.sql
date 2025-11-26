-- Migration: 005_create_vending_operations_table.sql
CREATE TABLE IF NOT EXISTS vending_operations (
    id BIGSERIAL PRIMARY KEY,
    vending_machine_id BIGINT NOT NULL,
    operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('restock', 'collection', 'maintenance')),
    performed_by BIGINT NOT NULL,
    operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    toys_before INTEGER,
    toys_after INTEGER,
    toys_added INTEGER DEFAULT 0,
    cash_before DECIMAL(10,2),
    cash_after DECIMAL(10,2),
    cash_collected DECIMAL(10,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vending_machine_id) REFERENCES vending_machines(id) ON DELETE CASCADE,
    FOREIGN KEY (performed_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_vending_operations_date ON vending_operations(operation_date);
CREATE INDEX IF NOT EXISTS idx_vending_operations_machine ON vending_operations(vending_machine_id);
CREATE INDEX IF NOT EXISTS idx_vending_operations_type ON vending_operations(operation_type);
CREATE INDEX IF NOT EXISTS idx_vending_operations_performed_by ON vending_operations(performed_by);