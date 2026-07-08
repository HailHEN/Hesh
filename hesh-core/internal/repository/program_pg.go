package repository

import (
    "context"
    "fmt"
    "hesh-core/internal/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProgramRepository struct {
    db *pgxpool.Pool
}

func NewPostgresProgramRepository(db *pgxpool.Pool) *PostgresProgramRepository {
    return &PostgresProgramRepository{db: db}
}

func (r *PostgresProgramRepository) CreateProduct(ctx context.Context, p *domain.StoreProduct) error {
    query := `
        INSERT INTO core.store_products (
            store_id, sku, item_name, category, item_description, 
            unit_price, is_available, points_earned, points_cost, 
            max_points_discount_allowed, is_redeemable
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at`

    err := r.db.QueryRow(ctx, query,
        p.StoreID, p.SKU, p.ItemName, p.Category, p.ItemDescription,
        p.UnitPrice, p.IsAvailable, p.PointsEarned, p.PointsCost,
        p.MaxPointsDiscountAllowed, p.IsRedeemable,
    ).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

    if err != nil {
        return fmt.Errorf("failed to insert store product: %w", err)
    }
    return nil
}
// Only allowed to create a program if merchant has paid or something etc
func (r *PostgresProgramRepository) CreatePointProgram(ctx context.Context, prog *domain.PointProgram) error {
    query := `
        INSERT INTO core.point_programs (merchant_id, points_per_dollar, is_active)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`

    err := r.db.QueryRow(ctx, query, prog.MerchantID, prog.PointsPerDollar, prog.IsActive).
        Scan(&prog.ID, &prog.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to set up point program: %w", err)
    }
    return nil
}

func (r *PostgresProgramRepository) CreateRewardPerk(ctx context.Context, perk *domain.RewardPerk) error {
    query := `
        INSERT INTO core.reward_perks (point_program_id, reward_name, points_required, reward_description, reward_action, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at`

    err := r.db.QueryRow(ctx, query,
        perk.PointProgramID, perk.RewardName, perk.PointsRequired, 
        perk.Description, perk.RewardAction, perk.IsActive,
    ).Scan(&perk.ID, &perk.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to save reward perk: %w", err)
    }
    return nil
}

func (r *PostgresProgramRepository) CreateCampaign(ctx context.Context, c *domain.MarketingCampaign) error {
    query := `
        INSERT INTO core.marketing_campaigns (
            merchant_id, store_id, is_archived, campaign_name, message_body, channels, 
            target_segments, reward_action, scheduled_at, expires_at, campaign_status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at`

    // UUID Safety Net: Swap an empty string StoreID for a safe database NULL pointer
    var storeID *string
    if c.StoreID != "" {
        storeID = &c.StoreID
    }

    // pgx handles *time.Time pointers natively, seamlessly converting nil to NULL
    err := r.db.QueryRow(ctx, query,
        c.MerchantID, storeID, c.IsArchived, c.CampaignName, c.MessageBody, c.Channels,
        c.TargetSegments, c.RewardAction, c.ScheduledAt, c.ExpiresAt, c.CampaignStatus,
    ).Scan(&c.ID, &c.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to create marketing campaign: %w", err)
    }
    return nil
}