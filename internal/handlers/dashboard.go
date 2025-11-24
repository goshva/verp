package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "vend_erp/internal/models"
)

type DashboardHandler struct {
    db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
    return &DashboardHandler{db: db}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
    stats, err := h.getDashboardStats()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    data := map[string]interface{}{
        "Stats":  stats,
        "Active": "dashboard",
    }
    
    renderTemplate(w, "dashboard.html", data)
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.getDashboardStats()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func (h *DashboardHandler) getDashboardStats() (*models.DashboardStats, error) {
    var stats models.DashboardStats
    
    // Total machines - используем правильное имя таблицы из вашей БД
    err := h.db.QueryRow("SELECT COUNT(*) FROM vending_machines").Scan(&stats.TotalMachines)
    if err != nil {
        return nil, err
    }
    
    // Active machines
    err = h.db.QueryRow("SELECT COUNT(*) FROM vending_machines WHERE status = 'active'").Scan(&stats.ActiveMachines)
    if err != nil {
        return nil, err
    }
    
    // Total users
    err = h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
    if err != nil {
        return nil, err
    }
    
    // Total videos
    err = h.db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&stats.TotalVideos)
    if err != nil {
        return nil, err
    }
    
    // Total revenue - используем правильные колонки из таблицы finance
    // В вашей БД finance имеет колонки: amount_views, amount_clicks, Amount (integer)
    err = h.db.QueryRow("SELECT COALESCE(SUM(amount_views), 0) + COALESCE(SUM(amount_clicks), 0) FROM finance WHERE status = 1").Scan(&stats.TotalRevenue)
    if err != nil {
        // Если все еще ошибка, используем альтернативный запрос
        err = h.db.QueryRow("SELECT COALESCE(SUM(amount_views), 0) FROM finance").Scan(&stats.TotalRevenue)
        if err != nil {
            stats.TotalRevenue = 0 // Устанавливаем 0 если не можем получить данные
        }
    }
    
    // Pending tasks - используем routes таблицу
    err = h.db.QueryRow("SELECT COUNT(*) FROM routes WHERE status = 'pending' OR status = 'planned'").Scan(&stats.PendingTasks)
    if err != nil {
        // Если таблицы routes нет, устанавливаем 0
        stats.PendingTasks = 0
    }
    
    return &stats, nil
}