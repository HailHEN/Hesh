package repository

import (
    "context"
    "fmt"
    "log"
    "errors"
    "hesh-core/internal/domain"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    
)

type PostgresTransactionRepository struct {
    db *pgxpool.Pool
}

func NewPostgresTransactionRepository(db *pgxpool.Pool) *PostgresTransactionRepository {
    return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) RecordPurchase(ctx context.Context, t *domain.Transaction) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("unable to start transaction block: %w", err)
    }
    defer func() {
        if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
            // Send it to your application logger instead of forcing it into your execution returns
            log.Printf("WARN: database transaction rollback failed: %v", err)
        }
    }()

    txQuery := `
        INSERT INTO core.transactions (
            pos_transaction_id, merchant_id, store_id, user_id, 
            applied_campaign_id, total_amount, points_earned, points_spent, currency, timestamp
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at`

    // UUID Safeties: Convert empty strings to nil pointers for Postgres
    var storeID, userID, campaignID *string
    if t.StoreID != "" {
        storeID = &t.StoreID
    }
    if t.UserID != "" {
        userID = &t.UserID
    }
    if t.AppliedCampaignID != "" {
        campaignID = &t.AppliedCampaignID
    }

    err = tx.QueryRow(ctx, txQuery,
        t.PosTransactionID, t.MerchantID, storeID, userID,
        campaignID, t.TotalAmount, t.PointsEarned, t.PointsSpent, t.Currency, t.Timestamp,
    ).Scan(&t.ID, &t.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to write primary transaction record: %w", err)
    }

    if len(t.LineItems) > 0 {
        batch := &pgx.Batch{}
        itemQuery := `
            INSERT INTO core.line_items (transaction_id, sku, item_name, category, quantity, unit_price)
            VALUES ($1, $2, $3, $4, $5, $6)`

        for i := range t.LineItems {
            batch.Queue(itemQuery, t.ID, t.LineItems[i].SKU, t.LineItems[i].ItemName, t.LineItems[i].Category, t.LineItems[i].Quantity, t.LineItems[i].UnitPrice)
        }

        br := tx.SendBatch(ctx, batch)
        for i := 0; i < len(t.LineItems); i++ {
            if _, err := br.Exec(); err != nil {
                _ = br.Close()
                return fmt.Errorf("failed to insert line item %d: %w", i, err)
            }
        }
        if err := br.Close(); err != nil {
            return fmt.Errorf("failed closing line items batch: %w", err)
        }
    }

    if len(t.AppliedRewards) > 0 {
        batch := &pgx.Batch{}
        rewardQuery := `
            INSERT INTO core.applied_rewards (transaction_id, reward_perk_id, store_product_id, quantity, points_deducted_per_unit)
            VALUES ($1, $2, $3, $4, $5)`

        for i := range t.AppliedRewards {
            var perkID, productID *string
            if t.AppliedRewards[i].RewardPerkID != "" {
                perkID = &t.AppliedRewards[i].RewardPerkID
            }
            if t.AppliedRewards[i].StoreProductID != "" {
                productID = &t.AppliedRewards[i].StoreProductID
            }

            batch.Queue(rewardQuery, t.ID, perkID, productID, t.AppliedRewards[i].Quantity, t.AppliedRewards[i].PointsDeductedPerUnit)
        }

        br := tx.SendBatch(ctx, batch)
        for i := 0; i < len(t.AppliedRewards); i++ {
            if _, err := br.Exec(); err != nil {
                _ = br.Close()
                return fmt.Errorf("failed to insert applied reward %d: %w", i, err)
            }
        }
        if err := br.Close(); err != nil {
            return fmt.Errorf("failed closing applied rewards batch: %w", err)
        }
    }

    if err = tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction ledger block: %w", err)
    }

    return nil
}