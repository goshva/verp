package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "time"
    "vend_erp/internal/models"
)

type MachineHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewMachineHandler(db *sql.DB, renderer *TemplateRenderer) *MachineHandler {
    return &MachineHandler{db: db, renderer: renderer}
}

func (h *MachineHandler) ListMachines(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: MachineHandler.ListMachines called for URL: %s\n", r.URL.Path)
    
    rows, err := h.db.Query(`
        SELECT 
            m.id, m.serial_number, m.model, m.status, 
            m.current_toys_count, m.capacity_toys, m.cash_amount, 
            m.last_maintenance_date, m.next_maintenance_date, m.installation_date,
            m.created_at, m.updated_at,
            m.location_id,
            COALESCE(l.name, 'Не назначена') as location_name
        FROM vending_machines m
        LEFT JOIN locations l ON m.location_id = l.id
        ORDER BY m.created_at DESC
    `)
    if err != nil {
        fmt.Printf("DEBUG: Machine query error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var machines []models.VendingMachine
    for rows.Next() {
        var machine models.VendingMachine
        var lastMaintenanceDate, nextMaintenanceDate, installationDate, createdAt, updatedAt sql.NullTime
        
        err := rows.Scan(
            &machine.ID, &machine.SerialNumber, &machine.Model, 
            &machine.Status, &machine.CurrentToysCount, &machine.CapacityToys,
            &machine.CashAmount, &lastMaintenanceDate, &nextMaintenanceDate, 
            &installationDate, &createdAt, &updatedAt,
            &machine.LocationID, &machine.LocationName,
        )
        if err != nil {
            fmt.Printf("Error scanning machine: %v\n", err)
            continue
        }
        
        // Обработка nullable дат
        if lastMaintenanceDate.Valid {
            machine.LastMaintenanceDate = lastMaintenanceDate.Time
        }
        if nextMaintenanceDate.Valid {
            machine.NextMaintenanceDate = nextMaintenanceDate.Time
        }
        if installationDate.Valid {
            machine.InstallationDate = installationDate.Time
        }
        if createdAt.Valid {
            machine.CreatedAt = createdAt.Time
        }
        if updatedAt.Valid {
            machine.UpdatedAt = updatedAt.Time
        }
        
        machines = append(machines, machine)
    }

    fmt.Printf("DEBUG: Loaded %d machines for template\n", len(machines))

    data := map[string]interface{}{
        "Machines": machines,
        "Active":   "machines",
        "Title":    "Автоматы",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("DEBUG: Rendering machines_list.html for HTMX\n")
        h.renderer.RenderTemplate(w, "machines_list.html", data)
        return
    }
    
    fmt.Printf("DEBUG: Rendering machines.html for full page with %d machines\n", len(machines))
    h.renderer.RenderTemplate(w, "machines.html", data)
}

func (h *MachineHandler) GetMachineForm(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: MachineHandler.GetMachineForm called\n")
    idStr := r.URL.Query().Get("id")
    var machine models.VendingMachine
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        err := h.db.QueryRow(`
            SELECT id, serial_number, model, status, location_id, 
                   capacity_toys, current_toys_count, cash_amount,
                   last_maintenance_date, next_maintenance_date, installation_date
            FROM vending_machines WHERE id = $1
        `, id).Scan(
            &machine.ID, &machine.SerialNumber, &machine.Model, 
            &machine.Status, &machine.LocationID, &machine.CapacityToys,
            &machine.CurrentToysCount, &machine.CashAmount,
            &machine.LastMaintenanceDate, &machine.NextMaintenanceDate, &machine.InstallationDate,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    // Fetch all active locations for dropdown
    locations, err := h.getActiveLocations()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    data := map[string]interface{}{
        "Machine":   machine,
        "Locations": locations,
        "Edit":      idStr != "",
    }
    h.renderer.RenderTemplate(w, "machine_form.html", data)
}

// Helper function to get active locations
func (h *MachineHandler) getActiveLocations() ([]models.Location, error) {
    rows, err := h.db.Query(`
        SELECT id, name, address 
        FROM locations 
        WHERE is_active = true 
        ORDER BY name
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var locations []models.Location
    for rows.Next() {
        var location models.Location
        err := rows.Scan(&location.ID, &location.Name, &location.Address)
        if err != nil {
            continue
        }
        locations = append(locations, location)
    }
    return locations, nil
}

func (h *MachineHandler) SaveMachine(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    locationID, _ := strconv.ParseInt(r.FormValue("location_id"), 10, 64)
    capacityToys, _ := strconv.Atoi(r.FormValue("capacity_toys"))
    currentToysCount, _ := strconv.Atoi(r.FormValue("current_toys_count"))
    cashAmount, _ := strconv.ParseFloat(r.FormValue("cash_amount"), 64)
    
    // Парсинг дат
    var lastMaintenanceDate, nextMaintenanceDate, installationDate time.Time
    if date := r.FormValue("last_maintenance_date"); date != "" {
        lastMaintenanceDate, _ = time.Parse("2006-01-02", date)
    }
    if date := r.FormValue("next_maintenance_date"); date != "" {
        nextMaintenanceDate, _ = time.Parse("2006-01-02", date)
    }
    if date := r.FormValue("installation_date"); date != "" {
        installationDate, _ = time.Parse("2006-01-02", date)
    }
    
    machine := models.VendingMachine{
        SerialNumber:        r.FormValue("serial_number"),
        Model:               r.FormValue("model"),
        LocationID:          locationID,
        CapacityToys:        capacityToys,
        CurrentToysCount:    currentToysCount,
        CashAmount:          cashAmount,
        Status:              r.FormValue("status"),
        LastMaintenanceDate: lastMaintenanceDate,
        NextMaintenanceDate: nextMaintenanceDate,
        InstallationDate:    installationDate,
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        _, err = h.db.Exec(`
            INSERT INTO vending_machines 
            (serial_number, model, location_id, status, capacity_toys, 
             current_toys_count, cash_amount, last_maintenance_date, 
             next_maintenance_date, installation_date)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, machine.SerialNumber, machine.Model, machine.LocationID, machine.Status,
           machine.CapacityToys, machine.CurrentToysCount, machine.CashAmount,
           nullIfZeroTime(machine.LastMaintenanceDate),
           nullIfZeroTime(machine.NextMaintenanceDate),
           nullIfZeroTime(machine.InstallationDate))
    } else {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        machine.ID = id
        
        _, err = h.db.Exec(`
            UPDATE vending_machines 
            SET serial_number=$1, model=$2, location_id=$3, status=$4,
                capacity_toys=$5, current_toys_count=$6, cash_amount=$7,
                last_maintenance_date=$8, next_maintenance_date=$9, 
                installation_date=$10, updated_at=CURRENT_TIMESTAMP
            WHERE id=$11
        `, machine.SerialNumber, machine.Model, machine.LocationID, machine.Status,
           machine.CapacityToys, machine.CurrentToysCount, machine.CashAmount,
           nullIfZeroTime(machine.LastMaintenanceDate),
           nullIfZeroTime(machine.NextMaintenanceDate),
           nullIfZeroTime(machine.InstallationDate), machine.ID)
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "machineSaved")
    h.ListMachines(w, r)
}

func (h *MachineHandler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    _, err = h.db.Exec("DELETE FROM vending_machines WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "machineDeleted")
    w.WriteHeader(http.StatusOK)
}

// Вспомогательная функция для обработки нулевых дат
func nullIfZeroTime(t time.Time) interface{} {
    if t.IsZero() {
        return nil
    }
    return t
}