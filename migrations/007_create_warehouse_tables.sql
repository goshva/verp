-- Migration: 007_create_warehouse_tables.sql

-- Таблица склада для хранения общей информации
CREATE TABLE IF NOT EXISTS warehouse (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    contact_person VARCHAR(255),
    contact_phone VARCHAR(50),
    total_capacity INTEGER NOT NULL, -- общая вместимость в единицах
    current_usage INTEGER DEFAULT 0, -- текущее использование
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица категорий товаров на складе
CREATE TABLE IF NOT EXISTS warehouse_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица инвентаря на складе
CREATE TABLE IF NOT EXISTS warehouse_inventory (
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    item_type VARCHAR(50) NOT NULL CHECK (item_type IN ('vending_machine', 'toy', 'capsule')),
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    min_stock_level INTEGER DEFAULT 10, -- минимальный уровень запаса для пополнения
    max_stock_level INTEGER DEFAULT 100, -- максимальный уровень запаса
    unit_price DECIMAL(10,2) DEFAULT 0,
    sku VARCHAR(100) UNIQUE, -- артикул
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES warehouse_categories(id) ON DELETE CASCADE
);

-- Таблица поставок на склад
CREATE TABLE IF NOT EXISTS warehouse_supplies (
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    supply_date DATE NOT NULL,
    expected_date DATE,
    status VARCHAR(50) DEFAULT 'ordered' CHECK (status IN ('ordered', 'in_transit', 'delivered', 'cancelled')),
    total_amount DECIMAL(10,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id) ON DELETE CASCADE
);

-- Таблица позиций в поставке
CREATE TABLE IF NOT EXISTS supply_items (
    id BIGSERIAL PRIMARY KEY,
    supply_id BIGINT NOT NULL,
    inventory_item_id BIGINT NOT NULL,
    quantity_ordered INTEGER NOT NULL,
    quantity_received INTEGER DEFAULT 0,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) GENERATED ALWAYS AS (quantity_ordered * unit_price) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supply_id) REFERENCES warehouse_supplies(id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_item_id) REFERENCES warehouse_inventory(id) ON DELETE CASCADE
);

-- Таблица отгрузок со склада
CREATE TABLE IF NOT EXISTS warehouse_shipments (
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL,
    shipment_type VARCHAR(50) NOT NULL CHECK (shipment_type IN ('to_location', 'to_courier', 'return', 'other')),
    target_location_id BIGINT, -- для отгрузки автоматов в локации
    courier_info TEXT, -- информация о курьере для отгрузки игрушек/капсул
    shipment_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'preparing' CHECK (status IN ('preparing', 'shipped', 'delivered', 'cancelled')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id) ON DELETE CASCADE,
    FOREIGN KEY (target_location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- Таблица позиций в отгрузке
CREATE TABLE IF NOT EXISTS shipment_items (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    inventory_item_id BIGINT NOT NULL,
    vending_machine_id BIGINT, -- для отгрузки конкретного автомата
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shipment_id) REFERENCES warehouse_shipments(id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_item_id) REFERENCES warehouse_inventory(id) ON DELETE CASCADE,
    FOREIGN KEY (vending_machine_id) REFERENCES vending_machines(id) ON DELETE SET NULL
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_warehouse_inventory_type ON warehouse_inventory(item_type);
CREATE INDEX IF NOT EXISTS idx_warehouse_inventory_quantity ON warehouse_inventory(quantity);
CREATE INDEX IF NOT EXISTS idx_warehouse_supplies_status ON warehouse_supplies(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_supplies_date ON warehouse_supplies(supply_date);
CREATE INDEX IF NOT EXISTS idx_warehouse_shipments_type ON warehouse_shipments(shipment_type);
CREATE INDEX IF NOT EXISTS idx_warehouse_shipments_status ON warehouse_shipments(status);
CREATE INDEX IF NOT EXISTS idx_warehouse_shipments_date ON warehouse_shipments(shipment_date);

-- Вставляем триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_warehouse_updated_at BEFORE UPDATE ON warehouse FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_warehouse_inventory_updated_at BEFORE UPDATE ON warehouse_inventory FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_warehouse_supplies_updated_at BEFORE UPDATE ON warehouse_supplies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_warehouse_shipments_updated_at BEFORE UPDATE ON warehouse_shipments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();