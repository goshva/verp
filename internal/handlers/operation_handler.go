package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "time"
    "vend_erp/internal/models"
)

type OperationHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewOperationHandler(db *sql.DB, renderer *TemplateRenderer) *OperationHandler {
    return &OperationHandler{db: db, renderer: renderer}
}

func (h *OperationHandler) ListOperations(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: OperationHandler.ListOperations called for URL: %s\n", r.URL.Path)
    
    rows, err := h.db.Query(`
        SELECT 
            o.id, o.vending_machine_id, o.operation_type, o.performed_by,
            o.operation_date, o.toys_before, o.toys_after, o.toys_added,
            o.cash_before, o.cash_after, o.cash_collected,
            o.created_at, o.updated_at,
            vm.serial_number as machine_serial,
            u.username as performer_name
        FROM vending_operations o
        LEFT JOIN vending_machines vm ON o.vending_machine_id = vm.id
        LEFT JOIN users u ON o.performed_by = u.id
        ORDER BY o.operation_date DESC
    `)
    if err != nil {
        fmt.Printf("DEBUG: Operations query error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var operations []models.VendingOperation
    for rows.Next() {
        var operation models.VendingOperation
        var operationDate, createdAt, updatedAt sql.NullTime
        
        err := rows.Scan(
            &operation.ID, &operation.VendingMachineID, &operation.OperationType, 
            &operation.PerformedBy, &operationDate, &operation.ToysBefore, 
            &operation.ToysAfter, &operation.ToysAdded, &operation.CashBefore,
            &operation.CashAfter, &operation.CashCollected, &createdAt, &updatedAt,
            &operation.MachineSerial, &operation.PerformerName,
        )
        if err != nil {
            fmt.Printf("Error scanning operation: %v\n", err)
            continue
        }
        
        // Handle nullable dates
        if operationDate.Valid {
            operation.OperationDate = operationDate.Time
        }
        if createdAt.Valid {
            operation.CreatedAt = createdAt.Time
        }
        if updatedAt.Valid {
            operation.UpdatedAt = updatedAt.Time
        }
        
        operations = append(operations, operation)
    }

    fmt.Printf("DEBUG: Loaded %d operations for template\n", len(operations))

    data := map[string]interface{}{
        "Operations": operations,
        "Active":     "operations",
        "Title":      "Операции",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("DEBUG: Rendering operations_list.html for HTMX\n")
        h.renderer.RenderTemplate(w, "operations_list.html", data)
        return
    }
    
    fmt.Printf("DEBUG: Rendering operations.html for full page with %d operations\n", len(operations))
    h.renderer.RenderTemplate(w, "operations_page.html", data)
}

func (h *OperationHandler) GetOperationForm(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: OperationHandler.GetOperationForm called\n")
    idStr := r.URL.Query().Get("id")
    var operation models.VendingOperation
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        err := h.db.QueryRow(`
            SELECT id, vending_machine_id, operation_type, performed_by,
                   operation_date, toys_before, toys_after, toys_added,
                   cash_before, cash_after, cash_collected
            FROM vending_operations WHERE id = $1
        `, id).Scan(
            &operation.ID, &operation.VendingMachineID, &operation.OperationType,
            &operation.PerformedBy, &operation.OperationDate, &operation.ToysBefore,
            &operation.ToysAfter, &operation.ToysAdded, &operation.CashBefore,
            &operation.CashAfter, &operation.CashCollected,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    // Fetch machines and users for dropdowns
    machines, err := h.getActiveMachines()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    users, err := h.getActiveUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    data := map[string]interface{}{
        "Operation": operation,
        "Machines":  machines,
        "Users":     users,
        "Edit":      idStr != "",
    }
    h.renderer.RenderTemplate(w, "operation_form.html", data)
}

// Helper function to get active machines
func (h *OperationHandler) getActiveMachines() ([]models.VendingMachine, error) {
    rows, err := h.db.Query(`
        SELECT vm.id, vm.serial_number, COALESCE(l.name, 'Не назначена') as location_name 
        FROM vending_machines vm
        LEFT JOIN locations l ON vm.location_id = l.id
        WHERE vm.status = 'active' 
        ORDER BY vm.serial_number
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var machines []models.VendingMachine
    for rows.Next() {
        var machine models.VendingMachine
        err := rows.Scan(&machine.ID, &machine.SerialNumber, &machine.LocationName)
        if err != nil {
            continue
        }
        machines = append(machines, machine)
    }
    return machines, nil
}

// Helper function to get active users
func (h *OperationHandler) getActiveUsers() ([]models.User, error) {
    rows, err := h.db.Query(`
        SELECT id, username, fullusername 
        FROM users 
        WHERE status = 1 
        ORDER BY username
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Username, &user.FullUserName)
        if err != nil {
            continue
        }
        users = append(users, user)
    }
    return users, nil
}

func (h *OperationHandler) SaveOperation(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    vendingMachineID, _ := strconv.ParseInt(r.FormValue("vending_machine_id"), 10, 64)
    performedBy, _ := strconv.ParseInt(r.FormValue("performed_by"), 10, 64)
    toysBefore, _ := strconv.Atoi(r.FormValue("toys_before"))
    toysAfter, _ := strconv.Atoi(r.FormValue("toys_after"))
    toysAdded, _ := strconv.Atoi(r.FormValue("toys_added"))
    cashBefore, _ := strconv.ParseFloat(r.FormValue("cash_before"), 64)
    cashAfter, _ := strconv.ParseFloat(r.FormValue("cash_after"), 64)
    cashCollected, _ := strconv.ParseFloat(r.FormValue("cash_collected"), 64)
    
    // Parse operation date
    var operationDate time.Time
    if date := r.FormValue("operation_date"); date != "" {
        operationDate, _ = time.Parse("2006-01-02T15:04", date) // For datetime-local
    } else {
        operationDate = time.Now()
    }
    
    operation := models.VendingOperation{
        VendingMachineID: vendingMachineID,
        OperationType:    r.FormValue("operation_type"),
        PerformedBy:      performedBy,
        OperationDate:    operationDate,
        ToysBefore:       toysBefore,
        ToysAfter:        toysAfter,
        ToysAdded:        toysAdded,
        CashBefore:       cashBefore,
        CashAfter:        cashAfter,
        CashCollected:    cashCollected,
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        _, err = h.db.Exec(`
            INSERT INTO vending_operations 
            (vending_machine_id, operation_type, performed_by, operation_date,
             toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, operation.VendingMachineID, operation.OperationType, operation.PerformedBy,
           operation.OperationDate, operation.ToysBefore, operation.ToysAfter,
           operation.ToysAdded, operation.CashBefore, operation.CashAfter,
           operation.CashCollected)
    } else {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        operation.ID = id
        
        _, err = h.db.Exec(`
            UPDATE vending_operations 
            SET vending_machine_id=$1, operation_type=$2, performed_by=$3, 
                operation_date=$4, toys_before=$5, toys_after=$6, toys_added=$7,
                cash_before=$8, cash_after=$9, cash_collected=$10, 
                updated_at=CURRENT_TIMESTAMP
            WHERE id=$11
        `, operation.VendingMachineID, operation.OperationType, operation.PerformedBy,
           operation.OperationDate, operation.ToysBefore, operation.ToysAfter,
           operation.ToysAdded, operation.CashBefore, operation.CashAfter,
           operation.CashCollected, operation.ID)
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "operationSaved")
    h.ListOperations(w, r)
}

func (h *OperationHandler) DeleteOperation(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    _, err = h.db.Exec("DELETE FROM vending_operations WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "operationDeleted")
    w.WriteHeader(http.StatusOK)
}