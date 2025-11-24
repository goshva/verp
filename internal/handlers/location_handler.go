package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "vend_erp/internal/models"
)

type LocationHandler struct {
    db *sql.DB
}

func NewLocationHandler(db *sql.DB) *LocationHandler {
    return &LocationHandler{db: db}
}

func (h *LocationHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
    // Явная проверка URL
    if r.URL.Path != "/locations" {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }
    
    fmt.Println("DEBUG: ListLocations called for path:", r.URL.Path)
    
    rows, err := h.db.Query(`
        SELECT id, name, address, contact_person, contact_phone, 
               monthly_rent, rent_due_day, is_active
        FROM locations ORDER BY created_at DESC
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var locations []models.Location
    for rows.Next() {
        var location models.Location
        err := rows.Scan(
            &location.ID, &location.Name, &location.Address,
            &location.ContactPerson, &location.ContactPhone,
            &location.MonthlyRent, &location.RentDueDay, &location.IsActive,
        )
        if err != nil {
            continue
        }
        locations = append(locations, location)
    }

    data := map[string]interface{}{
        "Locations": locations,
        "Active":    "locations",
        "Title":     "Локации",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        renderTemplate(w, "locations_list.html", data)
        return
    }
    
    renderTemplate(w, "locations.html", data)
}

func (h *LocationHandler) GetLocationForm(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    var location models.Location
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        err := h.db.QueryRow(`
            SELECT id, name, address, contact_person, contact_phone, 
                   monthly_rent, rent_due_day, is_active
            FROM locations WHERE id = $1
        `, id).Scan(
            &location.ID, &location.Name, &location.Address,
            &location.ContactPerson, &location.ContactPhone,
            &location.MonthlyRent, &location.RentDueDay, &location.IsActive,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    data := map[string]interface{}{
        "Location": location,
        "Edit":     idStr != "",
    }
    renderTemplate(w, "location_form.html", data)
}

func (h *LocationHandler) SaveLocation(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    monthlyRent, _ := strconv.ParseFloat(r.FormValue("monthly_rent"), 64)
    rentDueDay, _ := strconv.Atoi(r.FormValue("rent_due_day"))
    isActive := r.FormValue("is_active") == "true"
    
    location := models.Location{
        Name:          r.FormValue("name"),
        Address:       r.FormValue("address"),
        ContactPerson: r.FormValue("contact_person"),
        ContactPhone:  r.FormValue("contact_phone"),
        MonthlyRent:   monthlyRent,
        RentDueDay:    rentDueDay,
        IsActive:      isActive,
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        _, err = h.db.Exec(`
            INSERT INTO locations (name, address, contact_person, contact_phone, 
                                 monthly_rent, rent_due_day, is_active)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, location.Name, location.Address, location.ContactPerson, 
           location.ContactPhone, location.MonthlyRent, location.RentDueDay, location.IsActive)
    } else {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        location.ID = id
        
        _, err = h.db.Exec(`
            UPDATE locations 
            SET name=$1, address=$2, contact_person=$3, contact_phone=$4,
                monthly_rent=$5, rent_due_day=$6, is_active=$7,
                updated_at=CURRENT_TIMESTAMP
            WHERE id=$8
        `, location.Name, location.Address, location.ContactPerson,
           location.ContactPhone, location.MonthlyRent, location.RentDueDay, 
           location.IsActive, location.ID)
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "locationSaved")
    h.ListLocations(w, r)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    _, err = h.db.Exec("DELETE FROM locations WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "locationDeleted")
    w.WriteHeader(http.StatusOK)
}