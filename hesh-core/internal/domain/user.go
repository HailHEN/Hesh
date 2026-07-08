package domain

import (
    "context"
    "time"
)

type UserRepository interface {
    GetOrCreateUser(ctx context.Context, phoneNumber string) (*User, error)
    UpsertNetworkBalance(ctx context.Context, balance *NetworkBalance) error
    UpdateUserEmail(ctx context.Context, userID string, email string) error
}

type User struct {
    ID          string    `json:"id"`
    PhoneNumber string    `json:"phone_number"`
    Email       string    `json:"email"`
    WalletToken string    `json:"wallet_token"`
    CreatedAt   time.Time `json:"created_at"`
}

type NetworkBalance struct {
    UserID        string    `json:"user_id"`
    MerchantID    string    `json:"merchant_id"`
    PointsBalance int       `json:"points_balance"`
    UpdatedAt     time.Time `json:"updated_at"`
}