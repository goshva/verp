package handlers

import (
    "database/sql"
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    db       *sql.DB
    renderer *TemplateRenderer
}

func NewAuthHandler(db *sql.DB, renderer *TemplateRenderer) *AuthHandler {
    return &AuthHandler{db: db, renderer: renderer}
}

// Session represents a user session
type Session struct {
    ID        string
    UserID    int64
    ExpiresAt time.Time
}

// AuthData represents authentication form data
type AuthData struct {
    SignUp   bool
    Email    string
    Username string
    Error    string
    Title    string
    Active   string
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        data := AuthData{
            SignUp: false,
            Title:  "Вход в систему",
            Active: "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    email := r.FormValue("email")
    password := r.FormValue("password")
    
    var userID int64
    var hashedPassword, username string
    var status int
    
    err := h.db.QueryRow(`
        SELECT id, username, password, status 
        FROM users WHERE email = $1
    `, email).Scan(&userID, &username, &hashedPassword, &status)
    
    if err != nil {
        data := AuthData{
            SignUp: false,
            Email:  email,
            Error:  "Неверный email или пароль",
            Title:  "Вход в систему",
            Active: "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Check if user is active
    if status != 1 {
        data := AuthData{
            SignUp: false,
            Email:  email,
            Error:  "Аккаунт неактивен",
            Title:  "Вход в систему",
            Active: "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Verify password using bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        data := AuthData{
            SignUp: false,
            Email:  email,
            Error:  "Неверный email или пароль",
            Title:  "Вход в систему",
            Active: "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Create session
    sessionID, err := generateSessionID()
    if err != nil {
        http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
        return
    }
    
    expiresAt := time.Now().Add(24 * time.Hour)
    
    _, err = h.db.Exec(`
        INSERT INTO sessions (id, user_id, expires_at) 
        VALUES ($1, $2, $3)
    `, sessionID, userID, expiresAt)
    
    if err != nil {
        http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
        return
    }
    
    // Set session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Expires:  expiresAt,
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
    })
    
    http.Redirect(w, r, "/operations", http.StatusSeeOther)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        data := AuthData{
            SignUp: true,
            Title:  "Регистрация",
            Active: "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    email := r.FormValue("email")
    username := r.FormValue("username")
    password := r.FormValue("password")
    passwordConfirm := r.FormValue("password_confirm")
    
    // Validate input
    if password != passwordConfirm {
        data := AuthData{
            SignUp:   true,
            Email:    email,
            Username: username,
            Error:    "Пароли не совпадают",
            Title:    "Регистрация",
            Active:   "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    if len(password) < 6 {
        data := AuthData{
            SignUp:   true,
            Email:    email,
            Username: username,
            Error:    "Пароль должен содержать минимум 6 символов",
            Title:    "Регистрация",
            Active:   "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Check if user already exists
    var exists bool
    h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", 
        email, username).Scan(&exists)
    
    if exists {
        data := AuthData{
            SignUp:   true,
            Email:    email,
            Username: username,
            Error:    "Пользователь с таким email или именем уже существует",
            Title:    "Регистрация",
            Active:   "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Hash password using bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        data := AuthData{
            SignUp:   true,
            Email:    email,
            Username: username,
            Error:    "Ошибка создания аккаунта",
            Title:    "Регистрация",
            Active:   "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    // Create user
    _, err = h.db.Exec(`
        INSERT INTO users (username, email, password, userrole, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `, username, email, string(hashedPassword), "user", 1)
    
    if err != nil {
        data := AuthData{
            SignUp:   true,
            Email:    email,
            Username: username,
            Error:    "Ошибка создания аккаунта",
            Title:    "Регистрация",
            Active:   "auth",
        }
        h.renderer.RenderTemplate(w, "auth.html", data)
        return
    }
    
    http.Redirect(w, r, "/auth/signin?message=registered", http.StatusSeeOther)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        // Delete session from database
        h.db.Exec("DELETE FROM sessions WHERE id = $1", cookie.Value)
        
        // Clear cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    "",
            Expires:  time.Now().Add(-time.Hour),
            Path:     "/",
            HttpOnly: true,
        })
    }
    
    http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
}

func (h *AuthHandler) GetUserFromSession(r *http.Request) (*User, error) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return nil, err
    }
    
    var userID int64
    var expiresAt time.Time
    
    err = h.db.QueryRow(`
        SELECT user_id, expires_at 
        FROM sessions 
        WHERE id = $1 AND expires_at > CURRENT_TIMESTAMP
    `, cookie.Value).Scan(&userID, &expiresAt)
    
    if err != nil {
        return nil, err
    }
    
    var user User
    var fullUserName, companyName, companyRole, phone sql.NullString
    
    err = h.db.QueryRow(`
        SELECT id, username, email, userrole, status, 
               fullusername, companyname, companyrole, phone
        FROM users 
        WHERE id = $1 AND status = 1
    `, userID).Scan(
        &user.ID, &user.Username, &user.Email, &user.UserRole, 
        &user.Status, &fullUserName, &companyName, &companyRole, &phone,
    )
    
    if err != nil {
        return nil, err
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
    
    return &user, nil
}

func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := h.GetUserFromSession(r)
        if err != nil || user == nil {
            http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}

// User struct for authentication
type User struct {
    ID           int64
    Username     string
    Email        string
    UserRole     string
    Status       int
    FullUserName string
    CompanyName  string
    CompanyRole  string
    Phone        string
}