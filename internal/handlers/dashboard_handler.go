package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
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

    // Получаем реальные данные для графиков из БД
    revenueChart, _ := h.getRevenueChartData()
    operationsTrend, _ := h.getOperationsTrendData()
    inventoryChart, _ := h.getInventoryChartData()
    machinesActivity, _ := h.getMachinesActivityData()
    weeklyOperations, _ := h.getWeeklyOperationsData()

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
        
        // Данные для графиков из реальной БД
        "RevenueChart":       revenueChart,
        "OperationsTrend":    operationsTrend,
        "InventoryChart":     inventoryChart,
        "MachinesActivity":   machinesActivity,
        "WeeklyOperations":   weeklyOperations,
    }
    
    h.renderer.Render(w, "dashboard_page.html", data)
}