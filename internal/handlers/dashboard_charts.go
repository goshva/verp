package handlers

import (
    "time"
)

// ChartData представляет данные для графика
type ChartData struct {
    Labels   []string    `json:"labels"`
    Datasets []Dataset   `json:"datasets"`
}

type Dataset struct {
    Label           string    `json:"label"`
    Data           []float64 `json:"data"`
    BorderColor    string    `json:"borderColor"`
    BackgroundColor string   `json:"backgroundColor"`
    Fill           bool      `json:"fill"`
}

// Получение данных по доходам за последние 30 дней
func (h *DashboardHandler) getRevenueChartData() (ChartData, error) {
    var chart ChartData
    
    // Генерируем метки для последних 30 дней
    chart.Labels = h.generateLast30DaysLabels()
    
    query := `
        SELECT 
            DATE(operation_date) as date,
            SUM(cash_collected) as daily_revenue
        FROM vending_operations 
        WHERE operation_type = 'collection' 
        AND operation_date >= $1
        GROUP BY DATE(operation_date)
        ORDER BY date
    `
    
    rows, err := h.db.Query(query, time.Now().AddDate(0, 0, -30))
    if err != nil {
        return chart, err
    }
    defer rows.Close()
    
    // Создаем мапу для быстрого доступа к данным по датам
    revenueData := make(map[string]float64)
    for rows.Next() {
        var date string
        var revenue float64
        if err := rows.Scan(&date, &revenue); err == nil {
            revenueData[date] = revenue
        }
    }
    
    // Заполняем данные для графика
    var dataset Dataset
    dataset.Label = "Доход (₽)"
    dataset.BorderColor = "#10b981"
    dataset.BackgroundColor = "rgba(16, 185, 129, 0.1)"
    dataset.Fill = true
    
    // Сопоставляем данные с метками
    for _, label := range chart.Labels {
        // Преобразуем формат даты для сравнения
        dateObj, _ := time.Parse("02.01", label)
        dbDate := time.Now().AddDate(0, 0, -29).AddDate(0, 0, dateObj.Day()-1).Format("2006-01-02")
        
        if revenue, exists := revenueData[dbDate]; exists {
            dataset.Data = append(dataset.Data, revenue)
        } else {
            dataset.Data = append(dataset.Data, 0)
        }
    }
    
    chart.Datasets = []Dataset{dataset}
    return chart, nil
}

// Получение тренда операций за последние 30 дней
func (h *DashboardHandler) getOperationsTrendData() (ChartData, error) {
    var chart ChartData
    chart.Labels = h.generateLast30DaysLabels()
    
    query := `
        SELECT 
            DATE(operation_date) as date,
            COUNT(*) as operations_count
        FROM vending_operations 
        WHERE operation_date >= $1
        GROUP BY DATE(operation_date)
        ORDER BY date
    `
    
    rows, err := h.db.Query(query, time.Now().AddDate(0, 0, -30))
    if err != nil {
        return chart, err
    }
    defer rows.Close()
    
    operationsData := make(map[string]float64)
    for rows.Next() {
        var date string
        var count float64
        if err := rows.Scan(&date, &count); err == nil {
            operationsData[date] = count
        }
    }
    
    var dataset Dataset
    dataset.Label = "Количество операций"
    dataset.BorderColor = "#3b82f6"
    dataset.BackgroundColor = "rgba(59, 130, 246, 0.1)"
    dataset.Fill = true
    
    for _, label := range chart.Labels {
        dateObj, _ := time.Parse("02.01", label)
        dbDate := time.Now().AddDate(0, 0, -29).AddDate(0, 0, dateObj.Day()-1).Format("2006-01-02")
        
        if count, exists := operationsData[dbDate]; exists {
            dataset.Data = append(dataset.Data, count)
        } else {
            dataset.Data = append(dataset.Data, 0)
        }
    }
    
    chart.Datasets = []Dataset{dataset}
    return chart, nil
}

// Получение данных по стоимости инвентаря
func (h *DashboardHandler) getInventoryChartData() (ChartData, error) {
    var chart ChartData
    chart.Labels = h.generateLast30DaysLabels()
    
    query := `
        SELECT 
            DATE(wi.created_at) as date,
            SUM(wi.quantity * wi.unit_price) as inventory_value
        FROM warehouse_inventory wi
        JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE w.is_active = true 
        AND wi.created_at >= $1
        GROUP BY DATE(wi.created_at)
        ORDER BY date
    `
    
    rows, err := h.db.Query(query, time.Now().AddDate(0, 0, -30))
    if err != nil {
        return chart, err
    }
    defer rows.Close()
    
    inventoryData := make(map[string]float64)
    for rows.Next() {
        var date string
        var value float64
        if err := rows.Scan(&date, &value); err == nil {
            inventoryData[date] = value
        }
    }
    
    var dataset Dataset
    dataset.Label = "Стоимость инвентаря (₽)"
    dataset.BorderColor = "#f59e0b"
    dataset.BackgroundColor = "rgba(245, 158, 11, 0.1)"
    dataset.Fill = true
    
    for _, label := range chart.Labels {
        dateObj, _ := time.Parse("02.01", label)
        dbDate := time.Now().AddDate(0, 0, -29).AddDate(0, 0, dateObj.Day()-1).Format("2006-01-02")
        
        if value, exists := inventoryData[dbDate]; exists {
            dataset.Data = append(dataset.Data, value)
        } else {
            // Если данных нет, используем последнее известное значение или 0
            dataset.Data = append(dataset.Data, 0)
        }
    }
    
    chart.Datasets = []Dataset{dataset}
    return chart, nil
}

// Получение данных по активности автоматов
func (h *DashboardHandler) getMachinesActivityData() (ChartData, error) {
    var chart ChartData
    chart.Labels = h.generateLast30DaysLabels()
    
    query := `
        SELECT 
            DATE(vo.operation_date) as date,
            COUNT(DISTINCT vo.vending_machine_id) as active_machines
        FROM vending_operations vo
        WHERE vo.operation_date >= $1
        GROUP BY DATE(vo.operation_date)
        ORDER BY date
    `
    
    rows, err := h.db.Query(query, time.Now().AddDate(0, 0, -30))
    if err != nil {
        return chart, err
    }
    defer rows.Close()
    
    activityData := make(map[string]float64)
    for rows.Next() {
        var date string
        var activeCount float64
        if err := rows.Scan(&date, &activeCount); err == nil {
            activityData[date] = activeCount
        }
    }
    
    // Получаем общее количество автоматов для расчета процента
    var totalMachines int
    h.db.QueryRow("SELECT COUNT(*) FROM vending_machines WHERE status = 'active'").Scan(&totalMachines)
    
    var dataset Dataset
    dataset.Label = "Активность автоматов (%)"
    dataset.BorderColor = "#8b5cf6"
    dataset.BackgroundColor = "rgba(139, 92, 246, 0.1)"
    dataset.Fill = true
    
    for _, label := range chart.Labels {
        dateObj, _ := time.Parse("02.01", label)
        dbDate := time.Now().AddDate(0, 0, -29).AddDate(0, 0, dateObj.Day()-1).Format("2006-01-02")
        
        if activeCount, exists := activityData[dbDate]; exists && totalMachines > 0 {
            percentage := (activeCount / float64(totalMachines)) * 100
            dataset.Data = append(dataset.Data, percentage)
        } else {
            dataset.Data = append(dataset.Data, 0)
        }
    }
    
    chart.Datasets = []Dataset{dataset}
    return chart, nil
}

// Получение данных по операциям по дням недели
func (h *DashboardHandler) getWeeklyOperationsData() (ChartData, error) {
    var chart ChartData
    chart.Labels = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
    
    query := `
        SELECT 
            EXTRACT(DOW FROM operation_date) as day_of_week,
            operation_type,
            COUNT(*) as count
        FROM vending_operations 
        WHERE operation_date >= $1
        GROUP BY EXTRACT(DOW FROM operation_date), operation_type
        ORDER BY day_of_week
    `
    
    rows, err := h.db.Query(query, time.Now().AddDate(0, 0, -90)) // Данные за 3 месяца для статистики
    if err != nil {
        return chart, err
    }
    defer rows.Close()
    
    // PostgreSQL: 0=воскресенье, 1=понедельник, ..., 6=суббота
    // Преобразуем в: 0=понедельник, 1=вторник, ..., 6=воскресенье
    dayData := make(map[int]map[string]int)
    for i := 0; i < 7; i++ {
        dayData[i] = map[string]int{
            "restock":    0,
            "collection": 0,
            "maintenance": 0,
        }
    }
    
    for rows.Next() {
        var dayOfWeek int
        var opType string
        var count int
        if err := rows.Scan(&dayOfWeek, &opType, &count); err == nil {
            // Преобразуем день недели: PostgreSQL -> наш формат
            adjustedDay := (dayOfWeek + 6) % 7 // Сдвигаем: 0(вс)->6, 1(пн)->0, 2(вт)->1, etc.
            if dayData[adjustedDay] == nil {
                dayData[adjustedDay] = make(map[string]int)
            }
            dayData[adjustedDay][opType] += count
        }
    }
    
    // Создаем datasets для каждого типа операций
    restockData := make([]float64, 7)
    collectionData := make([]float64, 7)
    maintenanceData := make([]float64, 7)
    
    for i := 0; i < 7; i++ {
        restockData[i] = float64(dayData[i]["restock"])
        collectionData[i] = float64(dayData[i]["collection"])
        maintenanceData[i] = float64(dayData[i]["maintenance"])
    }
    
    chart.Datasets = []Dataset{
        {
            Label:           "Пополнения",
            Data:           restockData,
            BorderColor:    "#4CAF50",
            BackgroundColor: "rgba(76, 175, 80, 0.1)",
            Fill:           true,
        },
        {
            Label:           "Инкассации",
            Data:           collectionData,
            BorderColor:    "#2196F3",
            BackgroundColor: "rgba(33, 150, 243, 0.1)",
            Fill:           true,
        },
        {
            Label:           "Обслуживания",
            Data:           maintenanceData,
            BorderColor:    "#FF9800",
            BackgroundColor: "rgba(255, 152, 0, 0.1)",
            Fill:           true,
        },
    }
    
    return chart, nil
}

// Генерирует метки для последних 30 дней в формате "01.01", "02.01", ...
func (h *DashboardHandler) generateLast30DaysLabels() []string {
    var labels []string
    for i := 29; i >= 0; i-- {
        date := time.Now().AddDate(0, 0, -i)
        labels = append(labels, date.Format("02.01"))
    }
    return labels
}