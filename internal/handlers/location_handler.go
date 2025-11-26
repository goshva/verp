package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "vend_erp/internal/models"
)

type LocationHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewLocationHandler(db *sql.DB, renderer *TemplateRenderer) *LocationHandler {
    return &LocationHandler{db: db, renderer: renderer}
}

func (h *LocationHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: LocationHandler.ListLocations called for URL: %s\n", r.URL.Path)
    
    rows, err := h.db.Query(`
        SELECT id, name, address, contact_person, contact_phone, 
               monthly_rent, rent_due_day, is_active
        FROM locations ORDER BY created_at DESC
    `)
    if err != nil {
        fmt.Printf("DEBUG: Locations query error: %v\n", err)
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
            fmt.Printf("DEBUG: Location scan error: %v\n", err)
            continue
        }
        locations = append(locations, location)
    }

    fmt.Printf("DEBUG: Loaded %d locations\n", len(locations))

    data := map[string]interface{}{
        "Locations": locations,
        "Active":    "locations",
        "Title":     "Локации",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("DEBUG: Rendering locations_list.html for HTMX\n")
        h.renderer.Render(w, "locations_list.html", data)
        return
    }
    
    fmt.Printf("DEBUG: Rendering locations.html for full page\n")
    h.renderer.Render(w, "locations_page.html", data)
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
    h.renderer.Render(w, "location_form.html", data)
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
                monthly_rent=$5, rent_due_day=$6, is_active=$7
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