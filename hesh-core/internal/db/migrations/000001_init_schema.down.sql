SET search_path TO core, public;

-- ============================================================================
-- 1. DROP INDEXES
-- ============================================================================
DROP INDEX IF EXISTS core.idx_reward_perks_program;
DROP INDEX IF EXISTS core.idx_campaigns_merchant;
DROP INDEX IF EXISTS core.idx_rewards_merchant;
DROP INDEX IF EXISTS core.idx_store_products_category;
DROP INDEX IF EXISTS core.idx_store_products_store;
DROP INDEX IF EXISTS core.idx_transactions_user;
DROP INDEX IF EXISTS core.idx_transactions_store;
DROP INDEX IF EXISTS core.idx_applied_rewards_tx;
DROP INDEX IF EXISTS core.idx_store_merchant_id;

-- ============================================================================
-- 2. DROP TABLES (Strict Reverse Dependency Order)
-- ============================================================================
-- Child records containing foreign keys to transactions, perks, and products
DROP TABLE IF EXISTS core.applied_rewards CASCADE;
DROP TABLE IF EXISTS core.line_items CASCADE;

-- Transaction records
DROP TABLE IF EXISTS core.transactions CASCADE;

-- Core configurations, catalogs, and transactional trackers
DROP TABLE IF EXISTS core.reward_perks CASCADE;
DROP TABLE IF EXISTS core.point_programs CASCADE;
DROP TABLE IF EXISTS core.store_products CASCADE;
DROP TABLE IF EXISTS core.marketing_campaigns CASCADE;
DROP TABLE IF EXISTS core.network_balances CASCADE;

-- Foundation profiles and consumer profiles
DROP TABLE IF EXISTS core.users CASCADE;
DROP TABLE IF EXISTS core.store_addresses CASCADE;
DROP TABLE IF EXISTS core.store CASCADE;
DROP TABLE IF EXISTS core.merchant CASCADE;

-- ============================================================================
-- 3. DROP SCHEMA
-- ============================================================================
DROP SCHEMA IF EXISTS core CASCADE;