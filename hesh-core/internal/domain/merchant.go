package domain

import (
    "context"
    "time"
)

type MerchantRepository interface {
    CreateMerchant(ctx context.Context, m *Merchant) error
    CreateStoreWithAddress(ctx context.Context, s *Store) error
}

type Merchant struct {
    ID                  string    `json:"id"`
    BusinessName        string    `json:"business_name"`
    BusinessPhone       string    `json:"business_phone"`
    BusinessDescription string    `json:"business_description"`
    PosProvider         string    `json:"pos_provider"`
    PosMerchantID       string    `json:"pos_merchant_id"`
    BusinessABN         string    `json:"business_abn"`
    CreatedAt           time.Time `json:"created_at"`
}

type Store struct {
    ID         string       `json:"id"`
    MerchantID string       `json:"merchant_id"`
    StoreName  string       `json:"store_name"`
    IsActive   bool         `json:"is_active"`
    StorePhone string       `json:"store_phone,omitempty"`
    Address    StoreAddress `json:"address"`
    CreatedAt  time.Time    `json:"created_at"`
}

type StoreAddress struct {
    StoreID      string   `json:"store_id"`
    AddressLine1 string   `json:"address_line1"`
    AddressLine2 *string  `json:"address_line2,omitempty"` // Pointer handles optional NULL rows
    Suburb       string   `json:"suburb"`
    State        string   `json:"state"`
    Postcode     string   `json:"postcode"`
    Country      string   `json:"country"`
    Latitude     *float64 `json:"latitude,omitempty"`  // Pointer protects against 0.0 default values
    Longitude    *float64 `json:"longitude,omitempty"` // Pointer protects against 0.0 default values
    UpdatedAt    time.Time `json:"updated_at"`
}