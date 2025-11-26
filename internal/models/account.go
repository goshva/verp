// In models/user.go
package models

import "time"

type User struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    UserRole     string    `json:"userrole"`
    Status       int       `json:"status"`
    LastIPAddr   string    `json:"lastipaddr"`
    FullUserName string    `json:"fullusername"`
    CompanyName  string    `json:"companyname"`
    CompanyRole  string    `json:"companyrole"`
    Phone        string    `json:"phone"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

