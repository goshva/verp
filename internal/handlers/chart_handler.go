package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// ChartHandler обрабатывает данные для графиков и диаграмм
type ChartHandler struct {
	db *sql.DB
}

// NewChartHandler создает новый обработчик графиков
func NewChartHandler(db *sql.DB) *ChartHandler {
	return &ChartHandler{db: db}
}

// Структуры данных для графиков

// ChartDataPoint представляет точку данных на графике
type ChartDataPoint struct {
	Date       string  `json:"date"`
	Label      string  `json:"label"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
	Value      float64 `json:"value,omitempty"`
}

// ChartSeries представляет серию данных для графика
type ChartSeries struct {
	Name   string          `json:"name"`
	Color  string          `json:"color"`
	Data   []ChartDataPoint `json:"data"`
}

// ChartResponse содержит данные для отрисовки графика
type ChartResponse struct {
	Title         string       `json:"title"`
	Series        []ChartSeries `json:"series"`
	Labels        []string     `json:"labels"`
	Total         int          `json:"total"`
	Change        int          `json:"change"`
	ChangePercent float64      `json:"change_percent"`
	Trend         int          `json:"trend"` // -1 = down, 0 = stable, 1 = up
	Period        string       `json:"period"`
}

// MachineChartData возвращает данные для графика автоматов за последние 30 дней
func (h *ChartHandler) GetMachinesChartData() (*ChartResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Получаем количество автоматов по дням
	query := `
		WITH date_series AS (
			SELECT generate_series($1::date, $2::date, '1 day')::date as chart_date
		),
		daily_counts AS (
			SELECT 
				ds.chart_date,
				COUNT(DISTINCT vm.id) as machine_count
			FROM date_series ds
			LEFT JOIN vending_machines vm ON date(vm.created_at) <= ds.chart_date
			GROUP BY ds.chart_date
			ORDER BY ds.chart_date
		)
		SELECT 
			chart_date,
			machine_count
		FROM daily_counts
		ORDER BY chart_date
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []ChartDataPoint
	var counts []int
	var dates []time.Time

	for rows.Next() {
		var date time.Time
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			continue
		}

		counts = append(counts, count)
		dates = append(dates, date)

		label := h.formatDateLabel(date, startDate, endDate)

		dataPoints = append(dataPoints, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Label: label,
			Count: count,
		})
	}

	// Вычисляем проценты для высоты столбцов
	if len(counts) > 0 {
		maxCount := h.getMax(counts)
		if maxCount > 0 {
			for i := range dataPoints {
				dataPoints[i].Percentage = float64(counts[i]) / float64(maxCount) * 100
				// Минимальная высота 5% для видимости
				if dataPoints[i].Percentage < 5 {
					dataPoints[i].Percentage = 5
				}
			}
		}
	}

	// Получаем текущее общее количество автоматов
	var totalMachines int
	h.db.QueryRow("SELECT COUNT(*) FROM vending_machines").Scan(&totalMachines)

	// Рассчитываем изменения и тренд
	change, changePercent, trend := h.calculateMetrics(counts)

	// Создаем метки для оси X (каждые 5 дней или важные даты)
	labels := h.generateChartLabels(dates)

	response := &ChartResponse{
		Title:         "Динамика автоматов",
		Series: []ChartSeries{
			{
				Name:  "Автоматы",
				Color: "#4F46E5",
				Data:  dataPoints,
			},
		},
		Labels:        labels,
		Total:         totalMachines,
		Change:        change,
		ChangePercent: changePercent,
		Trend:         trend,
		Period:        "30 дней",
	}

	return response, nil
}

// GetOperationsChartData возвращает данные для графика операций
func (h *ChartHandler) GetOperationsChartData(days int) (*ChartResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT 
			date(created_at) as op_date,
			operation_type,
			COUNT(*) as operation_count
		FROM vending_operations
		WHERE created_at >= $1 AND created_at <= $2
		GROUP BY date(created_at), operation_type
		ORDER BY op_date, operation_type
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Группируем данные по типам операций
	operationsByType := make(map[string][]ChartDataPoint)
	dateSet := make(map[string]bool)

	for rows.Next() {
		var opDate time.Time
		var opType string
		var count int
		if err := rows.Scan(&opDate, &opType, &count); err != nil {
			continue
		}

		dateStr := opDate.Format("2006-01-02")
		dateSet[dateStr] = true

		if _, exists := operationsByType[opType]; !exists {
			operationsByType[opType] = []ChartDataPoint{}
		}

		label := h.formatDateLabel(opDate, startDate, endDate)
		operationsByType[opType] = append(operationsByType[opType], ChartDataPoint{
			Date:  dateStr,
			Label: label,
			Count: count,
		})
	}

	// Заполняем пропущенные даты нулевыми значениями
	allDates := h.generateDateRange(startDate, endDate)
	series := []ChartSeries{}

	colorMap := map[string]string{
		"restock":    "#10B981", // зеленый
		"collection": "#F59E0B", // желтый
		"maintenance": "#EF4444", // красный
	}

	for opType, dataPoints := range operationsByType {
		// Создаем map для быстрого доступа
		dataMap := make(map[string]int)
		for _, point := range dataPoints {
			dataMap[point.Date] = point.Count
		}

		// Создаем полный набор данных
		fullData := []ChartDataPoint{}
		for _, date := range allDates {
			dateStr := date.Format("2006-01-02")
			count := 0
			if val, exists := dataMap[dateStr]; exists {
				count = val
			}

			label := h.formatDateLabel(date, startDate, endDate)
			fullData = append(fullData, ChartDataPoint{
				Date:  dateStr,
				Label: label,
				Count: count,
			})
		}

		// Нормализуем проценты
		maxCount := h.getMaxFromPoints(fullData)
		if maxCount > 0 {
			for i := range fullData {
				fullData[i].Percentage = float64(fullData[i].Count) / float64(maxCount) * 100
			}
		}

		color, exists := colorMap[opType]
		if !exists {
			color = "#6B7280" // серый по умолчанию
		}

		series = append(series, ChartSeries{
			Name:  h.translateOperationType(opType),
			Color: color,
			Data:  fullData,
		})
	}

	// Суммарная статистика операций
	var totalOps int
	h.db.QueryRow("SELECT COUNT(*) FROM vending_operations WHERE created_at >= $1", startDate).Scan(&totalOps)

	// Создаем метки
	labels := h.generateChartLabels(allDates)

	response := &ChartResponse{
		Title:  "Операции с автоматами",
		Series: series,
		Labels: labels,
		Total:  totalOps,
		Period: h.formatPeriod(days),
	}

	return response, nil
}

// GetRevenueChartData возвращает данные для графика выручки
func (h *ChartHandler) GetRevenueChartData(days int) (*ChartResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT 
			date(operation_date) as revenue_date,
			SUM(cash_collected) as daily_revenue
		FROM vending_operations
		WHERE operation_type = 'collection' 
		  AND operation_date >= $1 
		  AND operation_date <= $2
		GROUP BY date(operation_date)
		ORDER BY revenue_date
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []ChartDataPoint
	var dates []time.Time
	var revenues []float64

	for rows.Next() {
		var date time.Time
		var revenue sql.NullFloat64
		if err := rows.Scan(&date, &revenue); err != nil {
			continue
		}

		rev := 0.0
		if revenue.Valid {
			rev = revenue.Float64
		}

		dates = append(dates, date)
		revenues = append(revenues, rev)

		label := h.formatDateLabel(date, startDate, endDate)
		dataPoints = append(dataPoints, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Label: label,
			Value: rev,
		})
	}

	// Вычисляем проценты
	if len(revenues) > 0 {
		maxRev := h.getMaxFloat(revenues)
		if maxRev > 0 {
			for i := range dataPoints {
				dataPoints[i].Percentage = revenues[i] / maxRev * 100
				if dataPoints[i].Percentage < 5 {
					dataPoints[i].Percentage = 5
				}
			}
		}
	}

	// Общая выручка за период
	var totalRevenue float64
	h.db.QueryRow("SELECT COALESCE(SUM(cash_collected), 0) FROM vending_operations WHERE operation_type = 'collection' AND operation_date >= $1", startDate).Scan(&totalRevenue)

	// Рассчитываем изменения
	change, changePercent, trend := h.calculateFloatMetrics(revenues)

	labels := h.generateChartLabels(dates)

	response := &ChartResponse{
		Title:  "Выручка",
		Series: []ChartSeries{
			{
				Name:  "Выручка (руб.)",
				Color: "#10B981",
				Data:  dataPoints,
			},
		},
		Labels:        labels,
		Total:         int(totalRevenue),
		Change:        int(change),
		ChangePercent: changePercent,
		Trend:         trend,
		Period:        h.formatPeriod(days),
	}

	return response, nil
}

// GetInventoryValueChartData возвращает данные для графика стоимости инвентаря
func (h *ChartHandler) GetInventoryValueChartData() (*ChartResponse, error) {
	query := `
		SELECT 
			date(wi.created_at) as inv_date,
			SUM(wi.quantity * wi.unit_price) as daily_value
		FROM warehouse_inventory wi
		GROUP BY date(wi.created_at)
		ORDER BY inv_date DESC
		LIMIT 30
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []ChartDataPoint
	var dates []time.Time
	var values []float64

	for rows.Next() {
		var date time.Time
		var value sql.NullFloat64
		if err := rows.Scan(&date, &value); err != nil {
			continue
		}

		val := 0.0
		if value.Valid {
			val = value.Float64
		}

		dates = append(dates, date)
		values = append(values, val)

		dataPoints = append(dataPoints, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Label: date.Format("02"),
			Value: val,
		})
	}

	// Реверсируем данные для правильного порядка
	for i, j := 0, len(dataPoints)-1; i < j; i, j = i+1, j-1 {
		dataPoints[i], dataPoints[j] = dataPoints[j], dataPoints[i]
		dates[i], dates[j] = dates[j], dates[i]
		values[i], values[j] = values[j], values[i]
	}

	// Вычисляем проценты
	if len(values) > 0 {
		maxVal := h.getMaxFloat(values)
		if maxVal > 0 {
			for i := range dataPoints {
				dataPoints[i].Percentage = values[i] / maxVal * 100
			}
		}
	}

	// Текущая стоимость инвентаря
	var totalValue float64
	h.db.QueryRow("SELECT COALESCE(SUM(quantity * unit_price), 0) FROM warehouse_inventory").Scan(&totalValue)

	// Рассчитываем изменения
	change, changePercent, trend := h.calculateFloatMetrics(values)

	labels := h.generateChartLabels(dates)

	response := &ChartResponse{
		Title:  "Стоимость инвентаря",
		Series: []ChartSeries{
			{
				Name:  "Стоимость (руб.)",
				Color: "#8B5CF6",
				Data:  dataPoints,
			},
		},
		Labels:        labels,
		Total:         int(totalValue),
		Change:        int(change),
		ChangePercent: changePercent,
		Trend:         trend,
		Period:        "30 дней",
	}

	return response, nil
}

// Вспомогательные методы

// formatDateLabel форматирует метку даты в зависимости от периода
func (h *ChartHandler) formatDateLabel(date, startDate, endDate time.Time) string {
	daysDiff := int(endDate.Sub(startDate).Hours() / 24)

	if daysDiff <= 7 {
		// Для недели: день недели
		return date.Format("Mon")[:2]
	} else if daysDiff <= 31 {
		// Для месяца: число
		if date.Day() == 1 || date.Day() == 15 || date.Day() == date.AddDate(0, 1, -1).Day() {
			return date.Format("02")
		} else if date.Day()%5 == 0 {
			return date.Format("02")
		} else {
			return "•"
		}
	} else {
		// Для длинных периодов: месяц.число
		if date.Day() == 1 {
			return date.Format("02 Jan")
		}
		return date.Format("02")
	}
}

// generateDateRange генерирует диапазон дат
func (h *ChartHandler) generateDateRange(start, end time.Time) []time.Time {
	var dates []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

// generateChartLabels генерирует метки для оси X
func (h *ChartHandler) generateChartLabels(dates []time.Time) []string {
	if len(dates) == 0 {
		return []string{}
	}

	var labels []string
	startDate := dates[0]
	endDate := dates[len(dates)-1]

	for _, date := range dates {
		labels = append(labels, h.formatDateLabel(date, startDate, endDate))
	}

	return labels
}

// formatPeriod форматирует текст периода
func (h *ChartHandler) formatPeriod(days int) string {
	switch {
	case days == 7:
		return "7 дней"
	case days == 30:
		return "30 дней"
	case days == 90:
		return "90 дней"
	default:
		return "30 дней"
	}
}

// translateOperationType переводит тип операции
func (h *ChartHandler) translateOperationType(opType string) string {
	switch opType {
	case "restock":
		return "Пополнение"
	case "collection":
		return "Инкассация"
	case "maintenance":
		return "Обслуживание"
	default:
		return opType
	}
}

// Метрики и вычисления

// calculateMetrics рассчитывает изменения и тренд для целых чисел
func (h *ChartHandler) calculateMetrics(values []int) (change int, changePercent float64, trend int) {
	if len(values) < 2 {
		return 0, 0, 0
	}

	startValue := values[0]
	endValue := values[len(values)-1]
	change = endValue - startValue

	if startValue != 0 {
		changePercent = float64(change) / float64(startValue) * 100
	}

	// Определяем тренд (анализ последних 7 дней)
	if len(values) >= 7 {
		last7 := values[len(values)-7:]
		first7 := values[:7]

		avgLast7 := h.averageInt(last7)
		avgFirst7 := h.averageInt(first7)

		if avgLast7 > avgFirst7*1.05 {
			trend = 1 // Рост
		} else if avgLast7 < avgFirst7*0.95 {
			trend = -1 // Спад
		}
	}

	return change, changePercent, trend
}

// calculateFloatMetrics рассчитывает изменения и тренд для дробных чисел
func (h *ChartHandler) calculateFloatMetrics(values []float64) (change float64, changePercent float64, trend int) {
	if len(values) < 2 {
		return 0, 0, 0
	}

	startValue := values[0]
	endValue := values[len(values)-1]
	change = endValue - startValue

	if startValue != 0 {
		changePercent = change / startValue * 100
	}

	// Определяем тренд
	if len(values) >= 7 {
		last7 := values[len(values)-7:]
		first7 := values[:7]

		avgLast7 := h.averageFloat(last7)
		avgFirst7 := h.averageFloat(first7)

		if avgLast7 > avgFirst7*1.05 {
			trend = 1
		} else if avgLast7 < avgFirst7*0.95 {
			trend = -1
		}
	}

	return change, changePercent, trend
}

// Вспомогательные математические функции
func (h *ChartHandler) getMax(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func (h *ChartHandler) getMaxFloat(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func (h *ChartHandler) getMaxFromPoints(points []ChartDataPoint) int {
	max := 0
	for _, point := range points {
		if point.Count > max {
			max = point.Count
		}
	}
	return max
}

func (h *ChartHandler) averageInt(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return float64(sum) / float64(len(nums))
}

func (h *ChartHandler) averageFloat(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

// GetTrendInfo возвращает информацию о тренде для использования в шаблонах
func (h *ChartHandler) GetTrendInfo(trend int) (string, string, string) {
	switch trend {
	case 1:
		return "up", "Рост", "📈"
	case -1:
		return "down", "Спад", "📉"
	default:
		return "stable", "Стабильно", "➡️"
	}
}

// GetMachinesChartJSON возвращает данные для графика автоматов в формате JSON
func (h *ChartHandler) GetMachinesChartJSON() ([]byte, error) {
	data, err := h.GetMachinesChartData()
	if err != nil {
		return nil, err
	}
	
	return json.Marshal(data)
}

// HandleMachinesChart обрабатывает HTTP запрос для данных графика автоматов
func (h *ChartHandler) HandleMachinesChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	data, err := h.GetMachinesChartJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// HandleOperationsChart обрабатывает HTTP запрос для данных графика операций
func (h *ChartHandler) HandleOperationsChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// По умолчанию 30 дней
	days := 30
	
	// Здесь вы можете извлечь параметры из URL если нужно
	// Например: /api/charts/operations?days=7
	
	data, err := h.GetOperationsChartData(days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleRevenueChart обрабатывает HTTP запрос для данных графика выручки
func (h *ChartHandler) HandleRevenueChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	days := 30
	
	data, err := h.GetRevenueChartData(days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleInventoryChart обрабатывает HTTP запрос для данных графика инвентаря
func (h *ChartHandler) HandleInventoryChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	data, err := h.GetInventoryValueChartData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
// GetCashChartData возвращает данные для графика денег в автоматах
func (h *ChartHandler) GetCashChartData() (*ChartResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Получаем сумму денег по дням
	query := `
		WITH date_series AS (
			SELECT generate_series($1::date, $2::date, '1 day')::date as chart_date
		),
		daily_cash AS (
			SELECT 
				ds.chart_date,
				COALESCE(SUM(vm.cash_amount), 0) as daily_cash_amount
			FROM date_series ds
			LEFT JOIN vending_machines vm ON date(vm.created_at) <= ds.chart_date
			GROUP BY ds.chart_date
			ORDER BY ds.chart_date
		)
		SELECT 
			chart_date,
			daily_cash_amount
		FROM daily_cash
		ORDER BY chart_date
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []ChartDataPoint
	var amounts []float64
	var dates []time.Time

	for rows.Next() {
		var date time.Time
		var amount float64
		if err := rows.Scan(&date, &amount); err != nil {
			continue
		}

		amounts = append(amounts, amount)
		dates = append(dates, date)

		label := h.formatDateLabel(date, startDate, endDate)

		dataPoints = append(dataPoints, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Label: label,
			Value: amount,
		})
	}

	// Получаем текущую общую сумму денег
	var totalCash float64
	h.db.QueryRow("SELECT COALESCE(SUM(cash_amount), 0) FROM vending_machines").Scan(&totalCash)

	// Рассчитываем изменения и тренд
	change, changePercent, trend := h.calculateFloatMetrics(amounts)

	// Создаем метки для оси X
	labels := h.generateChartLabels(dates)

	response := &ChartResponse{
		Title:         "Деньги в автоматах",
		Series: []ChartSeries{
			{
				Name:  "Деньги (руб.)",
				Color: "#10B981", // Зеленый цвет для денег
				Data:  dataPoints,
			},
		},
		Labels:        labels,
		Total:         int(totalCash),
		Change:        int(change),
		ChangePercent: changePercent,
		Trend:         trend,
		Period:        "30 дней",
	}

	return response, nil
}

// HandleCashChart обрабатывает HTTP запрос для данных графика денег
func (h *ChartHandler) HandleCashChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.GetCashChartData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}