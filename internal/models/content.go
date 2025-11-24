package models

import "time"

type Video struct {
    ID          int64     `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    URL         string    `json:"url" db:"url"`
    UserID      int64     `json:"user_id" db:"user_id"`
    Status      string    `json:"status" db:"status"`
    Duration    int       `json:"duration" db:"duration"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type VideoView struct {
    ID           int64     `json:"id" db:"id"`
    VideoID      int64     `json:"video_id" db:"video_id"`
    EquipmentID  int64     `json:"equipment_id" db:"equipment_id"`
    ViewersCount int16     `json:"viewers_count" db:"vq"`
    DateTime     time.Time `json:"datetime" db:"adddatetime"`
    GPSLat       string    `json:"gps_lat" db:"gpsposlat"`
    GPSLon       string    `json:"gps_lon" db:"gpsposlon"`
}
