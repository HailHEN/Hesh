CREATE SCHEMA IF NOT EXISTS core;
SET search_path TO core, public;
-- sets search path so we do not have to repeatedly set the schema

-- merchant table. Merchants can have many stores
CREATE TABLE merchant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_name VARCHAR(255) NOT NULL,
    business_phone VARCHAR(20) NOT NULL UNIQUE 
    CONSTRAINT valid_merchant_phone 
        CHECK (business_phone ~ '^\+[1-9]\d{6,14}$'), 
    business_description VARCHAR(255) NOT NULL,
    pos_provider VARCHAR(50) NOT NULL,
    pos_merchant_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE store (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchant(id) ON DELETE CASCADE,
    store_name VARCHAR(255) NOT NULL,
    
    store_phone VARCHAR(20) CONSTRAINT valid_store_phone CHECK (store_phone ~ '^\+[1-9]\d{6,14}$'),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE store_addresses (
    -- Making store_id the PRIMARY KEY enforces a strict 1:1 relationship. 
    -- A store_id can only appear once in this table.
    store_id UUID PRIMARY KEY REFERENCES store(id) ON DELETE CASCADE,
    
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    store_suburb VARCHAR(100) NOT NULL,    
    store_state VARCHAR(3) NOT NULL,       
    store_postcode VARCHAR(4) NOT NULL,    
    store_country VARCHAR(50) NOT NULL DEFAULT 'Australia',
    
    -- Geolocation coordinates
    latitude NUMERIC(9, 6),
    longitude NUMERIC(9, 6),
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_au_state CHECK (store_state IN ('NSW', 'VIC', 'QLD', 'WA', 'SA', 'TAS', 'ACT', 'NT'))
);
-- faster look up of merchants stores
CREATE INDEX idx_store_merchant_id ON store(merchant_id);
-- might change when web app is introduced
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL UNIQUE 
        CONSTRAINT valid_user_phone CHECK (phone_number ~ '^\+[1-9]\d{6,14}$'),
    wallet_token VARCHAR(255) UNIQUE, -- Used for Apple/Google Wallet push notifications
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- the amount of points accumulated by each user to each store
CREATE TABLE network_balances (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    merchant_id UUID REFERENCES merchant(id) ON DELETE CASCADE,
    points_balance INT NOT NULL DEFAULT 0 CHECK (points_balance >= 0),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, merchant_id)
);


-- tracks promotions/campaigns
-- gets consumed by worker
CREATE TABLE marketing_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchant(id) ON DELETE CASCADE,
    -- Optional: link to a specific store if it's a localized promo
    store_id UUID REFERENCES store(id) ON DELETE CASCADE, 
    
    campaign_name VARCHAR(255) NOT NULL,        -- e.g., "Surry Hills Winter Warm-up"
    message_body TEXT NOT NULL,                 -- The copy sent to the user

    channel VARCHAR(20)[] NOT NULL 
    CONSTRAINT valid_campaign_channels 
    CHECK (channel <@ ARRAY['sms', 'wallet_push', 'email']::VARCHAR[]),           -- 'wallet_push', 'sms', 'in_app'
    
    -- AI Targeting: WHO gets it?
    -- e.g., {"suburb": "Surry Hills", "days_since_last_visit": 30}
    target_segments JSONB NOT NULL DEFAULT '{}'::jsonb,
    
    -- Dynamic Outcome Payload: WHAT happens to the price/points?
    -- e.g., {"type": "discount", "target": "category", "match": "Food", "value": 0.15}
    -- e.g., {"type": "bonus_points", "multiplier": 2.0}
    reward_action JSONB NOT NULL,
    
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    campaign_status VARCHAR(20) NOT NULL DEFAULT 'pending'
    CONSTRAINT valid_campaign_status CHECK (campaign_status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- all transactions made via reward service
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pos_transaction_id VARCHAR(255) NOT NULL UNIQUE, 
    merchant_id UUID NOT NULL REFERENCES merchant(id),
    store_id UUID NOT NULL REFERENCES store(id), 
    user_id UUID NOT NULL REFERENCES users(id),
    applied_campaign_id UUID REFERENCES marketing_campaigns(id) ON DELETE SET NULL,
    total_amount NUMERIC(10, 2) NOT NULL, 
    points_earned INT NOT NULL DEFAULT 0 CHECK (points_earned >= 0),
    points_spent INT NOT NULL DEFAULT 0 CHECK (points_spent >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'AUD',
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- qucikly finds all transactions of a store 
CREATE INDEX idx_transactions_store ON transactions(store_id);
-- qucikly finds all transactions of a user 
CREATE INDEX idx_transactions_user ON transactions(user_id);

-- level 3 data
-- detailed transaction which includes information on what was purchased
-- stores should have all details of their products. Key for marketing
CREATE TABLE line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    sku VARCHAR(100),
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL, -- E.g., 'Beverage', 'Food', 'Apparel'
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10, 2) NOT NULL
);
-- menu of stores products
CREATE TABLE store_products(

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Links this product strictly to one specific store profile
    store_id UUID NOT NULL REFERENCES store(id) ON DELETE CASCADE,
    
    -- Product Attributes
    sku VARCHAR(100),                           
    item_name VARCHAR(255) NOT NULL,         
    category VARCHAR(100) NOT NULL,           
    item_description TEXT NOT NULL, -- useful for AI                     
    
    -- Pricing and Inventory Status
    unit_price NUMERIC(10, 2) NOT NULL,          -- Current retail price at this specific store
    is_available BOOLEAN NOT NULL DEFAULT TRUE,  -- Allows toggling out-of-stock items without deleting them
    -- points earned by buying this item
    points_earned INT DEFAULT 0 CHECK (points_earned >= 0),
    -- points cost to get this item
    points_cost INT CHECK (points_cost >= 0), 
    -- boundary for AI set by merchant
    max_points_discount_allowed INT CHECK (max_points_discount_allowed >= 0),
    -- flags if product can collect points
    is_redeemable BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);
-- quick search on all the product for a given store
CREATE INDEX idx_store_products_store ON store_products(store_id);
-- quick search on all a stores product, given the category
CREATE INDEX idx_store_products_category ON store_products(store_id, category);

-- ####################################################################################################################
-- In the future, each store can have a different default points program 
-- ####################################################################################################################

-- defines the stores average day to day royalty program
CREATE TABLE point_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchant(id) ON DELETE CASCADE,
    
    -- Earning Rules
    points_per_dollar NUMERIC(4, 2) NOT NULL DEFAULT 1.00, 
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- default reward perks that runs continously
CREATE TABLE reward_perks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    point_program_id UUID NOT NULL REFERENCES point_programs(id) ON DELETE CASCADE,
    
    -- Reward Details
    reward_name VARCHAR(255) NOT NULL,                        -- e.g. "Free Signature Coffee", "5 dollars off"
    points_required INT NOT NULL CHECK (points_required > 0), -- e.g. 100 points
    reward_description TEXT,                                   -- Context for the AI chatbot
    reward_action TEXT,  
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE applied_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    
    -- It can be a custom item from reward_perks OR a directly redeemable store product
    reward_perk_id UUID REFERENCES reward_perks(id) ON DELETE SET NULL,
    store_product_id UUID REFERENCES store_products(id) ON DELETE SET NULL,
    
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    points_deducted_per_unit INT NOT NULL CHECK (points_deducted_per_unit >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Safety check: Ensure at least one reference is populated
    CONSTRAINT check_reward_source CHECK (reward_perk_id IS NOT NULL OR store_product_id IS NOT NULL)
);


CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_store_addresses_modtime BEFORE UPDATE ON store_addresses FOR EACH ROW EXECUTE FUNCTION update_modified_column();
CREATE TRIGGER update_store_products_modtime BEFORE UPDATE ON store_products FOR EACH ROW EXECUTE FUNCTION update_modified_column();
CREATE TRIGGER update_network_balances_modtime BEFORE UPDATE ON network_balances FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- quickly finds all applied_rewards of a transaction (multiple rewards can be redeemed or used in a single transaction)
CREATE INDEX idx_applied_rewards_tx ON applied_rewards(transaction_id);

-- quick serach on all the point programs given the merchant
CREATE INDEX idx_rewards_merchant ON point_programs(merchant_id);
-- quick serach of all merchants campaigns  
CREATE INDEX idx_campaigns_merchant ON marketing_campaigns(merchant_id);
CREATE INDEX idx_reward_perks_program ON reward_perks(point_program_id) WHERE is_active = TRUE;