package domain

import (
    "context"
    "time"
)

type TransactionRepository interface {
    RecordPurchase(ctx context.Context, t *Transaction) error
}

type Transaction struct {
    ID                string          `json:"id"`
    PosTransactionID  string          `json:"pos_transaction_id"`
    MerchantID        string          `json:"merchant_id"`
    StoreID           string          `json:"store_id,omitempty"`
    UserID            string          `json:"user_id,omitempty"`
    AppliedCampaignID string          `json:"applied_campaign_id,omitempty"`
    TotalAmount       float64         `json:"total_amount"`
    PointsEarned      int             `json:"points_earned"`
    PointsSpent       int             `json:"points_spent"`
    Currency          string          `json:"currency"`
    Timestamp         time.Time       `json:"timestamp"`
    LineItems         []LineItem      `json:"line_items"`
    AppliedRewards    []AppliedReward `json:"applied_rewards,omitempty"`
    CreatedAt         time.Time       `json:"created_at"`
}

type LineItem struct {
    ID            string  `json:"id"`
    TransactionID string  `json:"transaction_id"`
    SKU           string  `json:"sku"`
    ItemName      string  `json:"item_name"`
    Category      string  `json:"category"`
    Quantity      int     `json:"quantity"`
    UnitPrice     float64 `json:"unit_price"`
}

type AppliedReward struct {
    ID                    string    `json:"id"`
    TransactionID         string    `json:"transaction_id"`
    RewardPerkID          string    `json:"reward_perk_id,omitempty"`
    StoreProductID        string    `json:"store_product_id,omitempty"`
    Quantity              int       `json:"quantity"`
    PointsDeductedPerUnit int       `json:"points_deducted_per_unit"`
    CreatedAt             time.Time `json:"created_at"`
}