package models

import "time"

type Warehouse struct {
    ID            int64     `json:"id"`
    Name          string    `json:"name"`
    Address       string    `json:"address"`
    ContactPerson string    `json:"contact_person"`
    ContactPhone  string    `json:"contact_phone"`
    TotalCapacity int       `json:"total_capacity"`
    CurrentUsage  int       `json:"current_usage"`
    IsActive      bool      `json:"is_active"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type WarehouseCategory struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}

type WarehouseInventory struct {
    ID            int64     `json:"id"`
    WarehouseID   int64     `json:"warehouse_id"`
    CategoryID    int64     `json:"category_id"`
    ItemType      string    `json:"item_type"` // vending_machine, toy, capsule
    ItemName      string    `json:"item_name"`
    Description   string    `json:"description"`
    Quantity      int       `json:"quantity"`
    MinStockLevel int       `json:"min_stock_level"`
    MaxStockLevel int       `json:"max_stock_level"`
    UnitPrice     float64   `json:"unit_price"`
    SKU           string    `json:"sku"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    
    // Joined fields
    WarehouseName    string `json:"warehouse_name"`
    WarehouseAddress string `json:"warehouse_address"`
    CategoryName     string `json:"category_name"`
}

type WarehouseSupply struct {
    ID           int64     `json:"id"`
    WarehouseID  int64     `json:"warehouse_id"`
    SupplierName string    `json:"supplier_name"`
    SupplyDate   time.Time `json:"supply_date"`
    ExpectedDate time.Time `json:"expected_date"`
    Status       string    `json:"status"` // ordered, in_transit, delivered, cancelled
    TotalAmount  float64   `json:"total_amount"`
    Notes        string    `json:"notes"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type SupplyItem struct {
    ID               int64   `json:"id"`
    SupplyID         int64   `json:"supply_id"`
    InventoryItemID  int64   `json:"inventory_item_id"`
    QuantityOrdered  int     `json:"quantity_ordered"`
    QuantityReceived int     `json:"quantity_received"`
    UnitPrice        float64 `json:"unit_price"`
    TotalPrice       float64 `json:"total_price"`
    CreatedAt        time.Time `json:"created_at"`
}

type WarehouseShipment struct {
    ID               int64     `json:"id"`
    WarehouseID      int64     `json:"warehouse_id"`
    ShipmentType     string    `json:"shipment_type"` // to_location, to_courier, return, other
    TargetLocationID *int64    `json:"target_location_id"`
    CourierInfo      string    `json:"courier_info"`
    ShipmentDate     time.Time `json:"shipment_date"`
    Status           string    `json:"status"` // preparing, shipped, delivered, cancelled
    Notes            string    `json:"notes"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type ShipmentItem struct {
    ID                int64 `json:"id"`
    ShipmentID        int64 `json:"shipment_id"`
    InventoryItemID   int64 `json:"inventory_item_id"`
    VendingMachineID  *int64 `json:"vending_machine_id"`
    Quantity          int   `json:"quantity"`
    CreatedAt         time.Time `json:"created_at"`
}
