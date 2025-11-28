package handlers

import (
    "vend_erp/internal/models"
)

type WarehouseStats struct {
    TotalValue      float64
    LowStockCount   int
    OutOfStockCount int
    TotalWarehouses int
}

func (h *DashboardHandler) getWarehouseStats() (WarehouseStats, error) {
    var stats WarehouseStats
    
    // Общая стоимость инвентаря
    err := h.db.QueryRow(`
        SELECT COALESCE(SUM(wi.quantity * wi.unit_price), 0) as total_value
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true
    `).Scan(&stats.TotalValue)
    if err != nil {
        return stats, err
    }
    
    // Позиции с низким запасом
    err = h.db.QueryRow(`
        SELECT COUNT(*) as low_stock_count
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true 
        AND wi.quantity < wi.min_stock_level 
        AND wi.quantity > 0
    `).Scan(&stats.LowStockCount)
    if err != nil {
        return stats, err
    }
    
    // Отсутствующие позиции
    err = h.db.QueryRow(`
        SELECT COUNT(*) as out_of_stock_count
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true AND wi.quantity = 0
    `).Scan(&stats.OutOfStockCount)
    if err != nil {
        return stats, err
    }
    
    // Активные склады
    err = h.db.QueryRow(`
        SELECT COUNT(*) as total_warehouses
        FROM warehouse 
        WHERE is_active = true
    `).Scan(&stats.TotalWarehouses)
    if err != nil {
        return stats, err
    }
    
    return stats, nil
}

type InventoryType struct {
    TypeName string
    Count    int
    Value    float64
}

func (h *DashboardHandler) getInventoryByType() ([]InventoryType, error) {
    rows, err := h.db.Query(`
        SELECT 
            CASE 
                WHEN item_type = 'vending_machine' THEN 'Автоматы'
                WHEN item_type = 'toy' THEN 'Игрушки' 
                WHEN item_type = 'capsule' THEN 'Капсулы'
                ELSE 'Другие'
            END as type_name,
            COUNT(*) as count,
            COALESCE(SUM(quantity * unit_price), 0) as total_value
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true
        GROUP BY item_type
        ORDER BY total_value DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var types []InventoryType
    for rows.Next() {
        var it InventoryType
        err := rows.Scan(&it.TypeName, &it.Count, &it.Value)
        if err != nil {
            continue
        }
        types = append(types, it)
    }
    
    return types, nil
}

func (h *DashboardHandler) getLowStockItems() ([]models.WarehouseInventory, error) {
    rows, err := h.db.Query(`
        SELECT 
            wi.id,
            wi.item_name,
            wi.quantity,
            wi.min_stock_level,
            wi.max_stock_level,
            wi.unit_price,
            w.name as warehouse_name
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true 
          AND wi.quantity < wi.min_stock_level 
          AND wi.quantity > 0
        ORDER BY (wi.quantity::float / wi.min_stock_level::float) ASC
        LIMIT 10
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var items []models.WarehouseInventory
    for rows.Next() {
        var item models.WarehouseInventory
        err := rows.Scan(
            &item.ID,
            &item.ItemName,
            &item.Quantity,
            &item.MinStockLevel,
            &item.MaxStockLevel,
            &item.UnitPrice,
            &item.WarehouseName,
        )
        if err != nil {
            continue
        }
        items = append(items, item)
    }
    
    return items, nil
}