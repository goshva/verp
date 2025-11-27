-- Таблица для логирования корректировок инвентаря
CREATE TABLE IF NOT EXISTS inventory_adjustments (
    id BIGSERIAL PRIMARY KEY,
    inventory_item_id BIGINT NOT NULL,
    adjustment_type VARCHAR(50) NOT NULL CHECK (adjustment_type IN ('add', 'subtract', 'set')),
    quantity INTEGER NOT NULL,
    new_quantity INTEGER NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_item_id) REFERENCES warehouse_inventory(id) ON DELETE CASCADE
);

-- Таблица для логирования перемещений между складами
CREATE TABLE IF NOT EXISTS inventory_transfers (
    id BIGSERIAL PRIMARY KEY,
    source_item_id BIGINT NOT NULL,
    target_item_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_item_id) REFERENCES warehouse_inventory(id) ON DELETE CASCADE,
    FOREIGN KEY (target_item_id) REFERENCES warehouse_inventory(id) ON DELETE CASCADE
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_item ON inventory_adjustments(inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_date ON inventory_adjustments(created_at);
CREATE INDEX IF NOT EXISTS idx_inventory_transfers_source ON inventory_transfers(source_item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_transfers_target ON inventory_transfers(target_item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_transfers_date ON inventory_transfers(created_at);