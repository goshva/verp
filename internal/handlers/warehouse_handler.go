package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "vend_erp/internal/models"
)

type WarehouseHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewWarehouseHandler(db *sql.DB, renderer *TemplateRenderer) *WarehouseHandler {
    return &WarehouseHandler{db: db, renderer: renderer}
}

func (h *WarehouseHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: WarehouseHandler.ListWarehouses called for URL: %s\n", r.URL.Path)
    
    // Получаем все склады
    warehouses, err := h.getActiveWarehouses()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Получаем инвентарь со всеми складами
    inventory, err := h.getInventoryWithFilters(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Рассчитываем статистику
    stats := h.calculateInventoryStats(inventory)
    
    data := map[string]interface{}{
        "Warehouses":        warehouses,
        "Inventory":         inventory,
        "TotalItems":        len(inventory),
        "TotalWarehouses":   len(warehouses),
        "TotalValue":        fmt.Sprintf("%.2f ₽", stats.TotalValue),
        "LowStockCount":     stats.LowStockCount,
        "OutOfStockCount":   stats.OutOfStockCount,
        "Active":            "warehouses",
        "Title":             "Склады и инвентарь",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        h.renderer.Render(w, "warehouses_list.html", data)
        return
    }
    
    h.renderer.Render(w, "warehouses_page.html", data)
}

func (h *WarehouseHandler) getInventoryWithFilters(r *http.Request) ([]models.WarehouseInventory, error) {
    warehouseID := r.URL.Query().Get("warehouse_id")
    category := r.URL.Query().Get("category")
    stockFilter := r.URL.Query().Get("stock")
    
    query := `
        SELECT 
            wi.id, wi.warehouse_id, wi.category_id, wi.item_type, 
            wi.item_name, wi.description, wi.quantity, wi.min_stock_level,
            wi.max_stock_level, wi.unit_price, wi.sku, wi.created_at, wi.updated_at,
            w.name as warehouse_name, w.address as warehouse_address,
            c.name as category_name
        FROM warehouse_inventory wi
        LEFT JOIN warehouse w ON wi.warehouse_id = w.id
        LEFT JOIN warehouse_categories c ON wi.category_id = c.id
        WHERE w.is_active = true
    `
    
    args := []interface{}{}
    argCount := 0
    
    if warehouseID != "" {
        argCount++
        query += fmt.Sprintf(" AND wi.warehouse_id = $%d", argCount)
        args = append(args, warehouseID)
    }
    
    if category != "" {
        argCount++
        query += fmt.Sprintf(" AND wi.item_type = $%d", argCount)
        args = append(args, category)
    }
    
    query += " ORDER BY w.name, wi.item_type, wi.item_name"
    
    rows, err := h.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var inventory []models.WarehouseInventory
    for rows.Next() {
        var item models.WarehouseInventory
        var createdAt, updatedAt sql.NullTime
        
        err := rows.Scan(
            &item.ID, &item.WarehouseID, &item.CategoryID, &item.ItemType,
            &item.ItemName, &item.Description, &item.Quantity, &item.MinStockLevel,
            &item.MaxStockLevel, &item.UnitPrice, &item.SKU, &createdAt, &updatedAt,
            &item.WarehouseName, &item.WarehouseAddress, &item.CategoryName,
        )
        if err != nil {
            continue
        }
        
        // Применяем фильтр по запасу
        if stockFilter != "" {
            if stockFilter == "low" && item.Quantity >= item.MinStockLevel {
                continue
            }
            if stockFilter == "out" && item.Quantity > 0 {
                continue
            }
            if stockFilter == "normal" && (item.Quantity < item.MinStockLevel || item.Quantity == 0) {
                continue
            }
        }
        
        if createdAt.Valid {
            item.CreatedAt = createdAt.Time
        }
        if updatedAt.Valid {
            item.UpdatedAt = updatedAt.Time
        }
        
        inventory = append(inventory, item)
    }
    
    return inventory, nil
}

func (h *WarehouseHandler) getActiveWarehouses() ([]models.Warehouse, error) {
    rows, err := h.db.Query(`
        SELECT id, name, address, contact_person, contact_phone,
               total_capacity, current_usage, is_active
        FROM warehouse 
        WHERE is_active = true 
        ORDER BY name
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var warehouses []models.Warehouse
    for rows.Next() {
        var warehouse models.Warehouse
        err := rows.Scan(
            &warehouse.ID, &warehouse.Name, &warehouse.Address,
            &warehouse.ContactPerson, &warehouse.ContactPhone,
            &warehouse.TotalCapacity, &warehouse.CurrentUsage, &warehouse.IsActive,
        )
        if err != nil {
            continue
        }
        warehouses = append(warehouses, warehouse)
    }
    return warehouses, nil
}

type InventoryStats struct {
    TotalValue      float64
    LowStockCount   int
    OutOfStockCount int
}

func (h *WarehouseHandler) calculateInventoryStats(inventory []models.WarehouseInventory) InventoryStats {
    var stats InventoryStats
    
    for _, item := range inventory {
        stats.TotalValue += float64(item.Quantity) * item.UnitPrice
        
        if item.Quantity == 0 {
            stats.OutOfStockCount++
        } else if item.Quantity < item.MinStockLevel {
            stats.LowStockCount++
        }
    }
    
    return stats
}

func (h *WarehouseHandler) GetWarehouseForm(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    var warehouse models.Warehouse
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        err := h.db.QueryRow(`
            SELECT id, name, address, contact_person, contact_phone,
                   total_capacity, current_usage, is_active
            FROM warehouse WHERE id = $1
        `, id).Scan(
            &warehouse.ID, &warehouse.Name, &warehouse.Address,
            &warehouse.ContactPerson, &warehouse.ContactPhone,
            &warehouse.TotalCapacity, &warehouse.CurrentUsage, &warehouse.IsActive,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    data := map[string]interface{}{
        "Warehouse": warehouse,
        "Edit":      idStr != "",
    }
    h.renderer.Render(w, "warehouse_form.html", data)
}

func (h *WarehouseHandler) SaveWarehouse(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    totalCapacity, _ := strconv.Atoi(r.FormValue("total_capacity"))
    isActive := r.FormValue("is_active") == "true"
    
    warehouse := models.Warehouse{
        Name:          r.FormValue("name"),
        Address:       r.FormValue("address"),
        ContactPerson: r.FormValue("contact_person"),
        ContactPhone:  r.FormValue("contact_phone"),
        TotalCapacity: totalCapacity,
        IsActive:      isActive,
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        _, err = h.db.Exec(`
            INSERT INTO warehouse (name, address, contact_person, contact_phone, 
                                 total_capacity, is_active)
            VALUES ($1, $2, $3, $4, $5, $6)
        `, warehouse.Name, warehouse.Address, warehouse.ContactPerson,
           warehouse.ContactPhone, warehouse.TotalCapacity, warehouse.IsActive)
    } else {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        warehouse.ID = id
        
        _, err = h.db.Exec(`
            UPDATE warehouse 
            SET name=$1, address=$2, contact_person=$3, contact_phone=$4,
                total_capacity=$5, is_active=$6, updated_at=CURRENT_TIMESTAMP
            WHERE id=$7
        `, warehouse.Name, warehouse.Address, warehouse.ContactPerson,
           warehouse.ContactPhone, warehouse.TotalCapacity, warehouse.IsActive, warehouse.ID)
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "warehouseSaved")
    h.ListWarehouses(w, r)
}

func (h *WarehouseHandler) GetInventoryForm(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    var inventoryItem models.WarehouseInventory
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        err := h.db.QueryRow(`
            SELECT id, warehouse_id, category_id, item_type, item_name,
                   description, quantity, min_stock_level, max_stock_level,
                   unit_price, sku
            FROM warehouse_inventory WHERE id = $1
        `, id).Scan(
            &inventoryItem.ID, &inventoryItem.WarehouseID, &inventoryItem.CategoryID,
            &inventoryItem.ItemType, &inventoryItem.ItemName, &inventoryItem.Description,
            &inventoryItem.Quantity, &inventoryItem.MinStockLevel, &inventoryItem.MaxStockLevel,
            &inventoryItem.UnitPrice, &inventoryItem.SKU,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    warehouses, _ := h.getActiveWarehouses()
    categories, _ := h.getCategories()
    
    data := map[string]interface{}{
        "InventoryItem": inventoryItem,
        "Warehouses":    warehouses,
        "Categories":    categories,
        "Edit":          idStr != "",
    }
    h.renderer.Render(w, "inventory_form.html", data)
}

func (h *WarehouseHandler) SaveInventory(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    warehouseID, _ := strconv.ParseInt(r.FormValue("warehouse_id"), 10, 64)
    categoryID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
    quantity, _ := strconv.Atoi(r.FormValue("quantity"))
    minStockLevel, _ := strconv.Atoi(r.FormValue("min_stock_level"))
    maxStockLevel, _ := strconv.Atoi(r.FormValue("max_stock_level"))
    unitPrice, _ := strconv.ParseFloat(r.FormValue("unit_price"), 64)
    
    inventoryItem := models.WarehouseInventory{
        WarehouseID:   warehouseID,
        CategoryID:    categoryID,
        ItemType:      r.FormValue("item_type"),
        ItemName:      r.FormValue("item_name"),
        Description:   r.FormValue("description"),
        Quantity:      quantity,
        MinStockLevel: minStockLevel,
        MaxStockLevel: maxStockLevel,
        UnitPrice:     unitPrice,
        SKU:           r.FormValue("sku"),
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        _, err = h.db.Exec(`
            INSERT INTO warehouse_inventory 
            (warehouse_id, category_id, item_type, item_name, description,
             quantity, min_stock_level, max_stock_level, unit_price, sku)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, inventoryItem.WarehouseID, inventoryItem.CategoryID, inventoryItem.ItemType,
           inventoryItem.ItemName, inventoryItem.Description, inventoryItem.Quantity,
           inventoryItem.MinStockLevel, inventoryItem.MaxStockLevel, inventoryItem.UnitPrice,
           inventoryItem.SKU)
    } else {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        inventoryItem.ID = id
        
        _, err = h.db.Exec(`
            UPDATE warehouse_inventory 
            SET warehouse_id=$1, category_id=$2, item_type=$3, item_name=$4,
                description=$5, quantity=$6, min_stock_level=$7, max_stock_level=$8,
                unit_price=$9, sku=$10, updated_at=CURRENT_TIMESTAMP
            WHERE id=$11
        `, inventoryItem.WarehouseID, inventoryItem.CategoryID, inventoryItem.ItemType,
           inventoryItem.ItemName, inventoryItem.Description, inventoryItem.Quantity,
           inventoryItem.MinStockLevel, inventoryItem.MaxStockLevel, inventoryItem.UnitPrice,
           inventoryItem.SKU, inventoryItem.ID)
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Обновляем использование склада
    h.updateWarehouseUsage(warehouseID)
    
    w.Header().Set("HX-Trigger", "inventorySaved")
    h.ListWarehouses(w, r)
}

func (h *WarehouseHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    // Получаем warehouse_id перед удалением для обновления использования
    var warehouseID int64
    err = h.db.QueryRow("SELECT warehouse_id FROM warehouse_inventory WHERE id = $1", id).Scan(&warehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    _, err = h.db.Exec("DELETE FROM warehouse_inventory WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Обновляем использование склада
    h.updateWarehouseUsage(warehouseID)
    
    w.Header().Set("HX-Trigger", "inventoryDeleted")
    h.ListWarehouses(w, r)
}

func (h *WarehouseHandler) updateWarehouseUsage(warehouseID int64) {
    h.db.Exec(`
        UPDATE warehouse 
        SET current_usage = (
            SELECT COALESCE(SUM(quantity), 0) 
            FROM warehouse_inventory 
            WHERE warehouse_id = $1
        )
        WHERE id = $1
    `, warehouseID)
}

func (h *WarehouseHandler) getCategories() ([]models.WarehouseCategory, error) {
    rows, err := h.db.Query("SELECT id, name FROM warehouse_categories ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var categories []models.WarehouseCategory
    for rows.Next() {
        var category models.WarehouseCategory
        err := rows.Scan(&category.ID, &category.Name)
        if err != nil {
            continue
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func (h *WarehouseHandler) GetQuickActionForm(w http.ResponseWriter, r *http.Request) {
    itemID := r.URL.Query().Get("item_id")
    actionType := r.URL.Query().Get("action")
    
    var item models.WarehouseInventory
    err := h.db.QueryRow(`
        SELECT wi.id, wi.quantity, wi.item_name, w.name as warehouse_name, w.id as warehouse_id
        FROM warehouse_inventory wi
        LEFT JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE wi.id = $1
    `, itemID).Scan(&item.ID, &item.Quantity, &item.ItemName, &item.WarehouseName, &item.WarehouseID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    warehouses, _ := h.getActiveWarehouses()
    
    data := map[string]interface{}{
        "ItemID":           item.ID,
        "ActionType":       actionType,
        "CurrentQuantity":  item.Quantity,
        "ItemName":         item.ItemName,
        "SourceWarehouse":  item.WarehouseName,
        "SourceWarehouseID": item.WarehouseID,
        "Warehouses":       warehouses,
        "Title":            getActionTitle(actionType),
    }
    
    h.renderer.Render(w, "quick_action_form.html", data)
}

func (h *WarehouseHandler) ExecuteQuickAction(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    itemID, _ := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
    actionType := r.FormValue("action_type")
    
    switch actionType {
    case "adjust":
        h.handleQuantityAdjustment(w, r, itemID)
    case "transfer":
        h.handleInventoryTransfer(w, r, itemID)
    default:
        http.Error(w, "Unknown action type", http.StatusBadRequest)
    }
}

func (h *WarehouseHandler) handleQuantityAdjustment(w http.ResponseWriter, r *http.Request, itemID int64) {
    adjustmentType := r.FormValue("adjustment_type")
    quantity, _ := strconv.Atoi(r.FormValue("quantity"))
    reason := r.FormValue("reason")
    
    var currentQuantity int
    var warehouseID int64
    err := h.db.QueryRow("SELECT quantity, warehouse_id FROM warehouse_inventory WHERE id = $1", itemID).
        Scan(&currentQuantity, &warehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    var newQuantity int
    switch adjustmentType {
    case "add":
        newQuantity = currentQuantity + quantity
    case "subtract":
        newQuantity = currentQuantity - quantity
        if newQuantity < 0 {
            newQuantity = 0
        }
    case "set":
        newQuantity = quantity
    }
    
    _, err = h.db.Exec(`
        UPDATE warehouse_inventory 
        SET quantity = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $2
    `, newQuantity, itemID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Логируем операцию
    h.db.Exec(`
        INSERT INTO inventory_adjustments 
        (inventory_item_id, adjustment_type, quantity, new_quantity, reason)
        VALUES ($1, $2, $3, $4, $5)
    `, itemID, adjustmentType, quantity, newQuantity, reason)
    
    // Обновляем использование склада
    h.updateWarehouseUsage(warehouseID)
    
    w.Header().Set("HX-Trigger", "inventoryAdjusted")
    h.ListWarehouses(w, r)
}

func (h *WarehouseHandler) handleInventoryTransfer(w http.ResponseWriter, r *http.Request, itemID int64) {
    quantity, _ := strconv.Atoi(r.FormValue("quantity"))
    targetWarehouseID, _ := strconv.ParseInt(r.FormValue("target_warehouse_id"), 10, 64)
    notes := r.FormValue("notes")
    
    // Получаем информацию об исходном товаре
    var sourceItem models.WarehouseInventory
    err := h.db.QueryRow(`
        SELECT wi.*, w.name as warehouse_name 
        FROM warehouse_inventory wi
        LEFT JOIN warehouse w ON wi.warehouse_id = w.id
        WHERE wi.id = $1
    `, itemID).Scan(
        &sourceItem.ID, &sourceItem.WarehouseID, &sourceItem.CategoryID, &sourceItem.ItemType,
        &sourceItem.ItemName, &sourceItem.Description, &sourceItem.Quantity, &sourceItem.MinStockLevel,
        &sourceItem.MaxStockLevel, &sourceItem.UnitPrice, &sourceItem.SKU, &sourceItem.CreatedAt, &sourceItem.UpdatedAt,
        &sourceItem.WarehouseName,
    )
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if sourceItem.Quantity < quantity {
        http.Error(w, "Недостаточно товара для перемещения", http.StatusBadRequest)
        return
    }
    
    // Находим или создаем запись в целевом складе
    var targetItemID int64
    err = h.db.QueryRow(`
        SELECT id FROM warehouse_inventory 
        WHERE warehouse_id = $1 AND sku = $2
    `, targetWarehouseID, sourceItem.SKU).Scan(&targetItemID)
    
    if err == sql.ErrNoRows {
        // Создаем новую запись в целевом складе
        err = h.db.QueryRow(`
            INSERT INTO warehouse_inventory 
            (warehouse_id, category_id, item_type, item_name, description,
             quantity, min_stock_level, max_stock_level, unit_price, sku)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id
        `, targetWarehouseID, sourceItem.CategoryID, sourceItem.ItemType,
           sourceItem.ItemName, sourceItem.Description, quantity,
           sourceItem.MinStockLevel, sourceItem.MaxStockLevel, sourceItem.UnitPrice,
           sourceItem.SKU).Scan(&targetItemID)
    } else if err == nil {
        // Обновляем существующую запись
        _, err = h.db.Exec(`
            UPDATE warehouse_inventory 
            SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP
            WHERE id = $2
        `, quantity, targetItemID)
    }
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Уменьшаем количество в исходном складе
    _, err = h.db.Exec(`
        UPDATE warehouse_inventory 
        SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `, quantity, itemID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Логируем перемещение
    h.db.Exec(`
        INSERT INTO inventory_transfers 
        (source_item_id, target_item_id, quantity, notes)
        VALUES ($1, $2, $3, $4)
    `, itemID, targetItemID, quantity, notes)
    
    // Обновляем использование обоих складов
    h.updateWarehouseUsage(sourceItem.WarehouseID)
    h.updateWarehouseUsage(targetWarehouseID)
    
    w.Header().Set("HX-Trigger", "inventoryTransferred")
    h.ListWarehouses(w, r)
}

func getActionTitle(actionType string) string {
    switch actionType {
    case "adjust":
        return "Корректировка количества"
    case "transfer":
        return "Перемещение между складами"
    default:
        return "Действие"
    }
}