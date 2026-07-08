package domain

import (
    "context"
    "encoding/json"
    "time"
)

type ProgramRepository interface {
    CreateProduct(ctx context.Context, p *StoreProduct) error
    CreatePointProgram(ctx context.Context, prog *PointProgram) error
    CreateRewardPerk(ctx context.Context, perk *RewardPerk) error
    CreateCampaign(ctx context.Context, c *MarketingCampaign) error
}

type StoreProduct struct {
    ID                       string    `json:"id"`
    StoreID                  string    `json:"store_id"`
    SKU                      string    `json:"sku,omitempty"`
    ItemName                 string    `json:"item_name"`
    Category                 string    `json:"category"`
    ItemDescription          string    `json:"item_description"`
    UnitPrice                float64   `json:"unit_price"`
    IsAvailable              bool      `json:"is_available"`
    PointsEarned             int       `json:"points_earned"`
    PointsCost               int       `json:"points_cost,omitempty"`
    MaxPointsDiscountAllowed int       `json:"max_points_discount_allowed,omitempty"`
    IsRedeemable             bool      `json:"is_redeemable"`
    CreatedAt                time.Time `json:"created_at"`
    UpdatedAt                time.Time `json:"updated_at"`
}

type PointProgram struct {
    ID              string    `json:"id"`
    MerchantID      string    `json:"merchant_id"`
    PointsPerDollar float64   `json:"points_per_dollar"`
    IsActive        bool      `json:"is_active"`
    CreatedAt       time.Time `json:"created_at"`
}

type RewardPerk struct {
    ID             string    `json:"id"`
    PointProgramID string    `json:"point_program_id"`
    RewardName     string    `json:"reward_name"`
    PointsRequired int       `json:"points_required"`
    Description    string    `json:"description,omitempty"`
    RewardAction   string    `json:"reward_action,omitempty"`
    IsActive       bool      `json:"is_active"`
    CreatedAt       time.Time `json:"created_at"`
}

type MarketingCampaign struct {
    ID             string          `json:"id"`
    MerchantID     string          `json:"merchant_id"`
    StoreID        string          `json:"store_id,omitempty"` // Can be blank for global merchant campaigns
    IsArchived     bool            `json:"is_archived"`        // Added to match final schema soft-delete
    CampaignName   string          `json:"campaign_name"`
    MessageBody    string          `json:"message_body"`
    Channels       []string        `json:"channels"`
    TargetSegments json.RawMessage `json:"target_segments"`    // Upgraded for native JSONB handling
    RewardAction   json.RawMessage `json:"reward_action"`       // Upgraded for native JSONB handling
    ScheduledAt    time.Time       `json:"scheduled_at"`
    ExpiresAt      *time.Time      `json:"expires_at,omitempty"` // Pointer handles optional expiration NULLs safely
    CampaignStatus string          `json:"campaign_status"`      // Casing aligned with DB column name
    CreatedAt      time.Time       `json:"created_at"`
}