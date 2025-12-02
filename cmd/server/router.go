package main

import (
	"database/sql"
	"net/http"

	"vend_erp/internal/handlers"
)

func setupRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	renderer := handlers.NewTemplateRenderer()

	// Handlers
	auth := handlers.NewAuthHandler(db, renderer)
	users := handlers.NewUserHandler(db, renderer)
	machines := handlers.NewMachineHandler(db, renderer)
	locations := handlers.NewLocationHandler(db, renderer)
	operations := handlers.NewOperationHandler(db, renderer)

	// Chart handler - создаем первым
	chartHandler := handlers.NewChartHandler(db)

	// Dashboard handler - передаем chartHandler
	dashboard := handlers.NewDashboardHandler(db, renderer, chartHandler)

	warehouses := handlers.NewWarehouseHandler(db, renderer)

	// Auth middleware closure
	requireAuth := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUserFromSession(r)
			if err != nil || user == nil {
				http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
				return
			}
			next(w, r)
		}
	}

	// Routes
	mux.HandleFunc("/auth/signin", auth.SignIn)
	mux.HandleFunc("/auth/signup", auth.SignUp)
	mux.HandleFunc("/auth/signout", auth.SignOut)
	mux.HandleFunc("/dashboard", requireAuth(dashboard.ShowDashboard))

	mux.HandleFunc("/accounts", requireAuth(users.ListUsers))
	mux.HandleFunc("/accounts/form", requireAuth(users.GetUserForm))
	mux.HandleFunc("/accounts/save", requireAuth(users.SaveUser))
	mux.HandleFunc("/accounts/delete", requireAuth(users.DeleteUser))

	mux.HandleFunc("/machines", requireAuth(machines.ListMachines))
	mux.HandleFunc("/machines/form", requireAuth(machines.GetMachineForm))
	mux.HandleFunc("/machines/save", requireAuth(machines.SaveMachine))
	mux.HandleFunc("/machines/delete", requireAuth(machines.DeleteMachine))

	mux.HandleFunc("/locations", requireAuth(locations.ListLocations))
	mux.HandleFunc("/locations/form", requireAuth(locations.GetLocationForm))
	mux.HandleFunc("/locations/save", requireAuth(locations.SaveLocation))
	mux.HandleFunc("/locations/delete", requireAuth(locations.DeleteLocation))

	mux.HandleFunc("/operations", requireAuth(operations.ListOperations))
	mux.HandleFunc("/operations/form", requireAuth(operations.GetOperationForm))
	mux.HandleFunc("/operations/save", requireAuth(operations.SaveOperation))
	mux.HandleFunc("/operations/delete", requireAuth(operations.DeleteOperation))

	mux.HandleFunc("/warehouses", requireAuth(warehouses.ListWarehouses))
	mux.HandleFunc("/warehouses/filter", requireAuth(warehouses.ListWarehouses))
	mux.HandleFunc("/warehouses/form", requireAuth(warehouses.GetWarehouseForm))
	mux.HandleFunc("/warehouses/save", requireAuth(warehouses.SaveWarehouse))
	mux.HandleFunc("/warehouses/inventory-form", requireAuth(warehouses.GetInventoryForm))
	mux.HandleFunc("/warehouses/inventory-save", requireAuth(warehouses.SaveInventory))
	mux.HandleFunc("/warehouses/inventory-delete", requireAuth(warehouses.DeleteInventory))
	mux.HandleFunc("/warehouses/quick-action", requireAuth(warehouses.GetQuickActionForm))
	mux.HandleFunc("/warehouses/quick-action-execute", requireAuth(warehouses.ExecuteQuickAction))

	// API routes for charts
	mux.HandleFunc("/api/charts/machines", requireAuth(chartHandler.HandleMachinesChart))
	mux.HandleFunc("/api/charts/operations", requireAuth(chartHandler.HandleOperationsChart))
	mux.HandleFunc("/api/charts/cash", chartHandler.HandleCashChart)
	mux.HandleFunc("/api/charts/revenue", requireAuth(chartHandler.HandleRevenueChart))
	mux.HandleFunc("/api/charts/inventory", requireAuth(chartHandler.HandleInventoryChart))
	mux.HandleFunc("/api/charts/toys", chartHandler.HandleToysChart)
	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
	})

	return mux
}
