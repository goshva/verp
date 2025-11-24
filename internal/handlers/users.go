package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    // УДАЛИТЬ: "time" - не используется
    "vend_erp/internal/models"
)

type UserHandler struct {
    db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: UserHandler.ListUsers called for URL: %s\n", r.URL.Path)
    fmt.Printf("DEBUG: Method: %s, HTMX: %s\n", r.Method, r.Header.Get("HX-Request"))
    
    rows, err := h.db.Query(`
        SELECT 
            id, username, email, userrole, status, lastipaddr,
            fullusername, companyname, companyrole, phone, 
            created_at, updated_at
        FROM users 
        ORDER BY created_at DESC
    `)
    if err != nil {
        fmt.Printf("DEBUG: User query error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        var createdAt, updatedAt sql.NullTime
        var lastIPAddr, fullUserName, companyName, companyRole, phone sql.NullString
        
        err := rows.Scan(
            &user.ID, &user.Username, &user.Email, &user.UserRole, 
            &user.Status, &lastIPAddr, &fullUserName, 
            &companyName, &companyRole, &phone, &createdAt, &updatedAt,
        )
        if err != nil {
            fmt.Printf("DEBUG: User scan error: %v\n", err)
            continue
        }
        
        // Handle nullable fields
        if lastIPAddr.Valid {
            user.LastIPAddr = lastIPAddr.String
        }
        if fullUserName.Valid {
            user.FullUserName = fullUserName.String
        }
        if companyName.Valid {
            user.CompanyName = companyName.String
        }
        if companyRole.Valid {
            user.CompanyRole = companyRole.String
        }
        if phone.Valid {
            user.Phone = phone.String
        }
        if createdAt.Valid {
            user.CreatedAt = createdAt.Time
        }
        if updatedAt.Valid {
            user.UpdatedAt = updatedAt.Time
        }
        
        users = append(users, user)
    }

    fmt.Printf("DEBUG: Loaded %d users\n", len(users))
    
    data := map[string]interface{}{
        "Users":  users,
        "Active": "users",
        "Title":  "Пользователи",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("DEBUG: Rendering users_list.html for HTMX\n")
        renderTemplate(w, "users_list.html", data)
        return
    }
    
    fmt.Printf("DEBUG: Rendering users.html for full page\n")
    renderTemplate(w, "users.html", data)
}

func (h *UserHandler) GetUserForm(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: UserHandler.GetUserForm called\n")
    idStr := r.URL.Query().Get("id")
    var user models.User
    
    if idStr != "" {
        id, _ := strconv.ParseInt(idStr, 10, 64)
        var fullUserName, companyName, companyRole, phone sql.NullString
        
        err := h.db.QueryRow(`
            SELECT id, username, email, userrole, status, 
                   fullusername, companyname, companyrole, phone
            FROM users WHERE id = $1
        `, id).Scan(
            &user.ID, &user.Username, &user.Email, &user.UserRole, 
            &user.Status, &fullUserName, &companyName, &companyRole, &phone,
        )
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Handle nullable fields
        if fullUserName.Valid {
            user.FullUserName = fullUserName.String
        }
        if companyName.Valid {
            user.CompanyName = companyName.String
        }
        if companyRole.Valid {
            user.CompanyRole = companyRole.String
        }
        if phone.Valid {
            user.Phone = phone.String
        }
    }
    
    data := map[string]interface{}{
        "User": user,
        "Edit": idStr != "",
    }
    renderTemplate(w, "user_form.html", data)
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: UserHandler.SaveUser called\n")
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    idStr := r.FormValue("id")
    status, _ := strconv.Atoi(r.FormValue("status"))
    
    user := models.User{
        Username:     r.FormValue("username"),
        Email:        r.FormValue("email"),
        UserRole:     r.FormValue("user_role"),
        Status:       status,
        FullUserName: r.FormValue("full_user_name"),
        CompanyName:  r.FormValue("company_name"),
        CompanyRole:  r.FormValue("company_role"),
        Phone:        r.FormValue("phone"),
    }
    
    var err error
    if idStr == "" || idStr == "0" {
        // Create new user
        password := r.FormValue("password")
        passwordConfirm := r.FormValue("password_confirm")
        
        if password != passwordConfirm {
            http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
            return
        }
        
        if password == "" {
            http.Error(w, "Пароль обязателен", http.StatusBadRequest)
            return
        }
        
        // For now, store plain text password (NOT FOR PRODUCTION)
        // In production, use proper hashing like bcrypt
        _, err = h.db.Exec(`
            INSERT INTO users (username, email, userrole, status, 
                             fullusername, companyname, companyrole, phone, password)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `, user.Username, user.Email, user.UserRole, user.Status,
           nullIfEmpty(user.FullUserName), nullIfEmpty(user.CompanyName), 
           nullIfEmpty(user.CompanyRole), nullIfEmpty(user.Phone), password)
    } else {
        // Update existing user
        id, _ := strconv.ParseInt(idStr, 10, 64)
        user.ID = id
        
        // Check if password is being updated
        password := r.FormValue("password")
        if password != "" {
            passwordConfirm := r.FormValue("password_confirm")
            if password != passwordConfirm {
                http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
                return
            }
            
            _, err = h.db.Exec(`
                UPDATE users 
                SET username=$1, email=$2, userrole=$3, status=$4, 
                    fullusername=$5, companyname=$6, companyrole=$7, phone=$8,
                    password=$9, updated_at=CURRENT_TIMESTAMP
                WHERE id=$10
            `, user.Username, user.Email, user.UserRole, user.Status,
               nullIfEmpty(user.FullUserName), nullIfEmpty(user.CompanyName), 
               nullIfEmpty(user.CompanyRole), nullIfEmpty(user.Phone), 
               password, user.ID)
        } else {
            _, err = h.db.Exec(`
                UPDATE users 
                SET username=$1, email=$2, userrole=$3, status=$4, 
                    fullusername=$5, companyname=$6, companyrole=$7, phone=$8,
                    updated_at=CURRENT_TIMESTAMP
                WHERE id=$9
            `, user.Username, user.Email, user.UserRole, user.Status,
               nullIfEmpty(user.FullUserName), nullIfEmpty(user.CompanyName), 
               nullIfEmpty(user.CompanyRole), nullIfEmpty(user.Phone), user.ID)
        }
    }
    
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "userSaved")
    h.ListUsers(w, r)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("DEBUG: UserHandler.DeleteUser called\n")
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    _, err = h.db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("HX-Trigger", "userDeleted")
    w.WriteHeader(http.StatusOK)
}

// Helper function for empty strings
func nullIfEmpty(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}

// УДАЛИТЕ весь дублированный код ниже этой линии!
// Не должно быть второго объявления ListUsers или других методов