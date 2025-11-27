package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "vend_erp/internal/models"
)

type DashboardHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewDashboardHandler(db *sql.DB, renderer *TemplateRenderer) *DashboardHandler {
    return &DashboardHandler{db: db, renderer: renderer}
}

func (h *DashboardHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: DashboardHandler.ShowDashboard called for URL: %s\n", r.URL.Path)
    
    // Получаем статистику складов
    stats, err := h.getWarehouseStats()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Получаем распределение по типам
    inventoryByType, err := h.getInventoryByType()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Получаем критические позиции
    lowStockItems, err := h.getLowStockItems()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Статистика автоматов
    var totalMachines, activeMachines int
    var totalCash float64
    var totalToys int

    h.db.QueryRow("SELECT COUNT(*) FROM vending_machines").Scan(&totalMachines)
    h.db.QueryRow("SELECT COUNT(*) FROM vending_machines WHERE status = 'active'").Scan(&activeMachines)
    h.db.QueryRow("SELECT COALESCE(SUM(cash_amount), 0) FROM vending_machines").Scan(&totalCash)
    h.db.QueryRow("SELECT COALESCE(SUM(current_toys_count), 0) FROM vending_machines").Scan(&totalToys)

    // Статистика операций
    var totalOperations, restockOperations, collectionOperations, maintenanceOperations int

    h.db.QueryRow("SELECT COUNT(*) FROM vending_operations").Scan(&totalOperations)
    h.db.QueryRow("SELECT COUNT(*) FROM vending_operations WHERE operation_type = 'restock'").Scan(&restockOperations)
    h.db.QueryRow("SELECT COUNT(*) FROM vending_operations WHERE operation_type = 'collection'").Scan(&collectionOperations)
    h.db.QueryRow("SELECT COUNT(*) FROM vending_operations WHERE operation_type = 'maintenance'").Scan(&maintenanceOperations)

    data := map[string]interface{}{
        "TotalValue":          fmt.Sprintf("%.2f ₽", stats.TotalValue),
        "LowStockCount":       stats.LowStockCount,
        "OutOfStockCount":     stats.OutOfStockCount,
        "TotalWarehouses":     stats.TotalWarehouses,
        "InventoryByType":     inventoryByType,
        "LowStockItems":       lowStockItems,
        "TotalMachines":       totalMachines,
        "ActiveMachines":      activeMachines,
        "TotalCash":           fmt.Sprintf("%.2f ₽", totalCash),
        "TotalToys":           totalToys,
        "TotalOperations":     totalOperations,
        "RestockOperations":   restockOperations,
        "CollectionOperations": collectionOperations,
        "MaintenanceOperations": maintenanceOperations,
        "Active":              "dashboard",
        "Title":               "Дашборд",
    }
    
    h.renderer.Render(w, "dashboard_page.html", data)
}

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
        WHERE w.is_active = true AND wi.quantity < wi.min_stock_level AND wi.quantity > 0
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
            COUNT(*) as count
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true
        GROUP BY item_type
        ORDER BY count DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var types []InventoryType
    for rows.Next() {
        var it InventoryType
        err := rows.Scan(&it.TypeName, &it.Count)
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
            wi.item_name,
            wi.quantity,
            wi.min_stock_level,
            w.name as warehouse_name
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true 
          AND wi.quantity < wi.min_stock_level 
          AND wi.quantity > 0
        ORDER BY (wi.min_stock_level - wi.quantity) DESC
        LIMIT 10
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var items []models.WarehouseInventory
    for rows.Next() {
        var item models.WarehouseInventory
        err := rows.Scan(&item.ItemName, &item.Quantity, &item.MinStockLevel, &item.WarehouseName)
        if err != nil {
            continue
        }
        items = append(items, item)
    }
    
    return items, nil
}