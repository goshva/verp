package main

import (
    "database/sql"
    "net/http"
    
    "vend_erp/internal/handlers"
)

// Router структура для управления всеми маршрутами
type Router struct {
    mux  *http.ServeMux
    db   *sql.DB
    auth *handlers.AuthHandler
}

// setupRoutes создает и настраивает все маршруты приложения
func setupRoutes(db *sql.DB) *Router {
    // Create template renderer first
    renderer := handlers.NewTemplateRenderer()
    
    r := &Router{
        mux:  http.NewServeMux(),
        db:   db,
        auth: handlers.NewAuthHandler(db, renderer),
    }
    
    // Setup all routes
    r.setupAllRoutes(renderer)
    
    return r
}

// setupAllRoutes настраивает все маршруты приложения
func (r *Router) setupAllRoutes(renderer *handlers.TemplateRenderer) {
    // Initialize all handlers from handlers package with renderer
    userHandler := handlers.NewUserHandler(r.db, renderer)
    machineHandler := handlers.NewMachineHandler(r.db, renderer)
    locationHandler := handlers.NewLocationHandler(r.db, renderer)
    operationHandler := handlers.NewOperationHandler(r.db, renderer)

    // Setup routes
    r.setupAuthRoutes()
    r.setupUserRoutes(userHandler)
    r.setupMachineRoutes(machineHandler)
    r.setupLocationRoutes(locationHandler)
    r.setupOperationRoutes(operationHandler)
    
    // Static files
    r.setupStaticRoutes()
    
    // Root route
    //r.mux.HandleFunc("/", r.rootHandler)
}

// setupAuthRoutes настраивает маршруты аутентификации
func (r *Router) setupAuthRoutes() {
    r.mux.HandleFunc("/auth/signin", r.auth.SignIn)
    r.mux.HandleFunc("/auth/signup", r.auth.SignUp)
    r.mux.HandleFunc("/auth/signout", r.auth.SignOut)
}

// setupUserRoutes настраивает маршруты пользователей
func (r *Router) setupUserRoutes(h *handlers.UserHandler) {
    r.mux.HandleFunc("/accounts", r.requireAuth(h.ListUsers))
    r.mux.HandleFunc("/accounts/form", r.requireAuth(h.GetUserForm))
    r.mux.HandleFunc("/accounts/save", r.requireAuth(h.SaveUser))
    r.mux.HandleFunc("/accounts/delete", r.requireAuth(h.DeleteUser))
}

// setupMachineRoutes настраивает маршруты автоматов
func (r *Router) setupMachineRoutes(h *handlers.MachineHandler) {
    r.mux.HandleFunc("/machines", r.requireAuth(h.ListMachines))
    r.mux.HandleFunc("/machines/form", r.requireAuth(h.GetMachineForm))
    r.mux.HandleFunc("/machines/save", r.requireAuth(h.SaveMachine))
    r.mux.HandleFunc("/machines/delete", r.requireAuth(h.DeleteMachine))
}

// setupLocationRoutes настраивает маршруты локаций
func (r *Router) setupLocationRoutes(h *handlers.LocationHandler) {
    r.mux.HandleFunc("/locations", r.requireAuth(h.ListLocations))
    r.mux.HandleFunc("/locations/form", r.requireAuth(h.GetLocationForm))
    r.mux.HandleFunc("/locations/save", r.requireAuth(h.SaveLocation))
    r.mux.HandleFunc("/locations/delete", r.requireAuth(h.DeleteLocation))
}

// setupOperationRoutes настраивает маршруты операций
func (r *Router) setupOperationRoutes(h *handlers.OperationHandler) {
    r.mux.HandleFunc("/operations", r.requireAuth(h.ListOperations))
    r.mux.HandleFunc("/operations/form", r.requireAuth(h.GetOperationForm))
    r.mux.HandleFunc("/operations/save", r.requireAuth(h.SaveOperation))
    r.mux.HandleFunc("/operations/delete", r.requireAuth(h.DeleteOperation))
}

// setupStaticRoutes настраивает маршруты для статических файлов
func (r *Router) setupStaticRoutes() {
    // Serve static files (CSS, JS, images)
    fs := http.FileServer(http.Dir("static"))
    r.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

// ServeHTTP реализует интерфейс http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    r.mux.ServeHTTP(w, req)
}

// requireAuth middleware для проверки аутентификации
func (r *Router) requireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        user, err := r.auth.GetUserFromSession(req)
        if err != nil || user == nil {
            http.Redirect(w, req, "/auth/signin", http.StatusSeeOther)
            return
        }
        next(w, req)
    }
}

// rootHandler обрабатывает корневой маршрут
func (r *Router) rootHandler(w http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/" {
        http.NotFound(w, req)
        return
    }
    
    // Проверяем аутентификацию для корневого маршрута
    user, err := r.auth.GetUserFromSession(req)
    if err != nil || user == nil {
        http.Redirect(w, req, "/auth/signin", http.StatusSeeOther)
        return
    }
    
    // Перенаправляем на страницу пользователей по умолчанию
    //http.Redirect(w, req, "/operations", http.StatusSeeOther)
}