package repository

import (
    "context"
    "fmt"
    "hesh-core/internal/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMerchantRepository struct {
    db *pgxpool.Pool
}

func NewPostgresMerchantRepository(db *pgxpool.Pool) *PostgresMerchantRepository {
    return &PostgresMerchantRepository{db: db}
}

func (r *PostgresMerchantRepository) CreateMerchant(ctx context.Context, m *domain.Merchant) error {
    query := `
        INSERT INTO core.merchant (business_name, business_phone, business_description, pos_provider, pos_merchant_id, business_abn)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at`

    err := r.db.QueryRow(ctx, query, 
        m.BusinessName, m.BusinessPhone, m.BusinessDescription, m.PosProvider, m.PosMerchantID, m.BusinessABN,
    ).Scan(&m.ID, &m.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to create merchant: %w", err)
    }
    return nil
}

func (r *PostgresMerchantRepository) CreateStoreWithAddress(ctx context.Context, s *domain.Store) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to start store transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Updated: explicitly tracking and setting the is_active status flag
    storeQuery := `
        INSERT INTO core.store (merchant_id, store_name, is_active, store_phone)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`

    err = tx.QueryRow(ctx, storeQuery, s.MerchantID, s.StoreName, s.IsActive, s.StorePhone).Scan(&s.ID, &s.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert store: %w", err)
    }

    addressQuery := `
        INSERT INTO core.store_addresses (
            store_id, address_line1, address_line2, store_suburb, 
            store_state, store_postcode, store_country, latitude, longitude
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING updated_at`

    // Text Safety: If address line 2 is an empty string, handle it as a clean database NULL pointer
    var addressLine2 *string
    if s.Address.AddressLine2 != nil && *s.Address.AddressLine2 != "" {
        addressLine2 = s.Address.AddressLine2
    }

    err = tx.QueryRow(ctx, addressQuery,
        s.ID, s.Address.AddressLine1, addressLine2, s.Address.Suburb,
        s.Address.State, s.Address.Postcode, s.Address.Country, s.Address.Latitude, s.Address.Longitude,
    ).Scan(&s.Address.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to insert store address: %w", err)
    }

    if err = tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit store profile transaction: %w", err)
    }

    s.Address.StoreID = s.ID
    return nil
}