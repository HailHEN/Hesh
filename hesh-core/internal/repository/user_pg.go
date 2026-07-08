package repository

import (
	"context"
	"fmt"
	"hesh-core/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}
/* Implementing the interface functions defined in the user domain */

/* 
DO UPDATE ... is here in order to be able to return something. Safe for login etc
*/
func (r *PostgresUserRepository) GetOrCreateUser(ctx context.Context, phoneNumber string) (*domain.User, error) {
	query := `
        INSERT INTO core.users (phone_number)
        VALUES ($1)
        ON CONFLICT (phone_number) 
        DO UPDATE SET phone_number = core.users.phone_number
        RETURNING id, phone_number, COALESCE(email, ''), , COALESCE(wallet_token, ''), created_at`

	var u domain.User
	err := r.db.QueryRow(ctx, query, phoneNumber).Scan(&u.ID, &u.PhoneNumber, &u.Email, &u.WalletToken, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to set user profile: %w", err)
	}
	return &u, nil
}

func (r *PostgresUserRepository) UpsertNetworkBalance(ctx context.Context, balance *domain.NetworkBalance) error {
	query := `
		INSERT INTO core.network_balances (user_id, merchant_id, points_balance)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, merchant_id)
		DO UPDATE SET points_balance = core.network_balances.points_balance + EXCLUDED.points_balance, updated_at = NOW()
		RETURNING points_balance, updated_at`

	err := r.db.QueryRow(ctx, query, balance.UserID, balance.MerchantID, balance.PointsBalance).
		Scan(&balance.PointsBalance, &balance.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user network points balance: %w", err)
	}
	return nil
}

// user registers or updates email address
// when sigining after basic sign up (done in their own time)
func (r *PostgresUserRepository) UpdateUserEmail(ctx context.Context, userID string, email string) error {
    query := `
        UPDATE core.users 
        SET email = $1 
        WHERE id = $2`

    // nothing needs to be returned aside from error
    _, err := r.db.Exec(ctx, query, email, userID)
    if err != nil {
        return fmt.Errorf("failed to update user email in database: %w", err)
    }

    return nil
}