package models

type DashboardStats struct {
    TotalMachines  int     `json:"total_machines"`
    ActiveMachines int     `json:"active_machines"`
    TotalRevenue   float64 `json:"total_revenue"`
    PendingTasks   int     `json:"pending_tasks"`
    TotalUsers     int     `json:"total_users"`
    TotalVideos    int     `json:"total_videos"`
}
