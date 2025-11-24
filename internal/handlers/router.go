package handlers

import (
    "database/sql"
    "fmt"
    "html/template"
    "net/http"
)

var tmpl *template.Template

func init() {
    tmpl = template.Must(template.ParseGlob("internal/templates/*.html"))
    fmt.Printf("DEBUG: Loaded templates: %s\n", tmpl.DefinedTemplates())
    
    // Проверка на конфликты
    templates := []string{"machines.html", "users.html", "locations.html"}
    for _, tmplName := range templates {
        t := tmpl.Lookup(tmplName)
        if t == nil {
            fmt.Printf("WARNING: Template %s not found!\n", tmplName)
        } else {
            fmt.Printf("DEBUG: Template %s found\n", tmplName)
        }
    }
}
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
    fmt.Printf("DEBUG: Rendering template: %s with data: %+v\n", name, data)
    
    // Добавьте заголовки против кэширования
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    
    err := tmpl.ExecuteTemplate(w, name, data)
    if err != nil {
        fmt.Printf("DEBUG: Template error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func SetupRoutes(mux *http.ServeMux, db *sql.DB) {
    userHandler := NewUserHandler(db)
    machineHandler := NewMachineHandler(db)
    dashboardHandler := NewDashboardHandler(db)
    locationHandler := NewLocationHandler(db)
    
    fmt.Println("DEBUG: Setting up routes...")
    
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    
    // Dashboard
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Root request to: %s\n", r.URL.Path)
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    })
    
    mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Dashboard request: %s\n", r.URL.Path)
        dashboardHandler.Dashboard(w, r)
    })
    
    // Users routes
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Users request: %s\n", r.URL.Path)
        userHandler.ListUsers(w, r)
    })
    mux.HandleFunc("/users/form", userHandler.GetUserForm)
    mux.HandleFunc("/users/save", userHandler.SaveUser)
    mux.HandleFunc("/users/delete", userHandler.DeleteUser)
    
    // Location routes
    mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Locations request: %s\n", r.URL.Path)
        locationHandler.ListLocations(w, r)
    })
    mux.HandleFunc("/locations/form", locationHandler.GetLocationForm)
    mux.HandleFunc("/locations/save", locationHandler.SaveLocation)
    mux.HandleFunc("/locations/delete", locationHandler.DeleteLocation)
    
    // Machine routes
    mux.HandleFunc("/machines", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Machines request: %s\n", r.URL.Path)
        machineHandler.ListMachines(w, r)
    })
    mux.HandleFunc("/machines/form", machineHandler.GetMachineForm)
    mux.HandleFunc("/machines/save", machineHandler.SaveMachine)
    mux.HandleFunc("/machines/delete", machineHandler.DeleteMachine)
    
    // API routes
    mux.HandleFunc("/api/stats", dashboardHandler.GetStats)
    
    fmt.Println("DEBUG: Routes setup complete")
    fmt.Println("DEBUG: Available routes:")
    fmt.Println("  GET  / -> /dashboard")
    fmt.Println("  GET  /dashboard")
    fmt.Println("  GET  /users")
    fmt.Println("  GET  /machines")
    fmt.Println("  GET  /locations")
    fmt.Println("  GET  /api/stats")
    fmt.Println("  POST /locations/save") // Добавьте это

}
