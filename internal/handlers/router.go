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
    templates := []string{"machines.html", "users.html", "locations.html", "auth.html"}
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
    authHandler := NewAuthHandler(db)
    
    fmt.Println("DEBUG: Setting up routes...")
    
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    
    // Auth routes
    mux.HandleFunc("/auth/signin", authHandler.SignIn)
    mux.HandleFunc("/auth/signup", authHandler.SignUp)
    mux.HandleFunc("/auth/signout", authHandler.SignOut)
    
    // Dashboard
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("DEBUG: Root request to: %s\n", r.URL.Path)
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    })
    
    // Add authentication middleware to protected routes
    mux.HandleFunc("/dashboard", authHandler.RequireAuth(dashboardHandler.Dashboard))
    mux.HandleFunc("/users", authHandler.RequireAuth(userHandler.ListUsers))
    mux.HandleFunc("/users/form", authHandler.RequireAuth(userHandler.GetUserForm))
    mux.HandleFunc("/users/save", authHandler.RequireAuth(userHandler.SaveUser))
    mux.HandleFunc("/users/delete", authHandler.RequireAuth(userHandler.DeleteUser))
    mux.HandleFunc("/locations", authHandler.RequireAuth(locationHandler.ListLocations))
    mux.HandleFunc("/locations/form", authHandler.RequireAuth(locationHandler.GetLocationForm))
    mux.HandleFunc("/locations/save", authHandler.RequireAuth(locationHandler.SaveLocation))
    mux.HandleFunc("/locations/delete", authHandler.RequireAuth(locationHandler.DeleteLocation))
    mux.HandleFunc("/machines", authHandler.RequireAuth(machineHandler.ListMachines))
    mux.HandleFunc("/machines/form", authHandler.RequireAuth(machineHandler.GetMachineForm))
    mux.HandleFunc("/machines/save", authHandler.RequireAuth(machineHandler.SaveMachine))
    mux.HandleFunc("/machines/delete", authHandler.RequireAuth(machineHandler.DeleteMachine))
    
    // API routes
    mux.HandleFunc("/api/stats", authHandler.RequireAuth(dashboardHandler.GetStats))
    
    fmt.Println("DEBUG: Routes setup complete")
    fmt.Println("DEBUG: Available routes:")
    fmt.Println("  GET  / -> /dashboard")
    fmt.Println("  GET  /auth/signin")
    fmt.Println("  GET  /auth/signup")
    fmt.Println("  POST /auth/signin")
    fmt.Println("  POST /auth/signup")
    fmt.Println("  POST /auth/signout")
    fmt.Println("  GET  /dashboard")
    fmt.Println("  GET  /users")
    fmt.Println("  GET  /machines")
    fmt.Println("  GET  /locations")
    fmt.Println("  GET  /api/stats")
    fmt.Println("  POST /locations/save")
}