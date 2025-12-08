package queries

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	Pool *pgxpool.Pool
}

func NewQuery(pool *pgxpool.Pool) *Query { return &Query{Pool: pool} }

func (qr *Query) InitializeDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	queries := []string{
		// ============================================================
		// Extensions
		// ============================================================
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,

		// ============================================================
		// Sequences
		// ============================================================
		`CREATE SEQUENCE IF NOT EXISTS admin_seq START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS master_distributor_seq START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS distributor_seq START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS user_seq START 1;`,

		// ============================================================
		// Common trigger function
		// ============================================================
		`CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END; $$ LANGUAGE plpgsql;`,

		// ============================================================
		// Admins
		// ============================================================
		`CREATE TABLE IF NOT EXISTS admins (
			admin_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_unique_id TEXT UNIQUE NOT NULL DEFAULT ('A' || LPAD(nextval('admin_seq')::TEXT, 7, '0')),
			admin_name TEXT NOT NULL,
			admin_phone TEXT UNIQUE NOT NULL,
			admin_email TEXT UNIQUE NOT NULL,
			admin_password TEXT NOT NULL,
			admin_wallet_balance NUMERIC(20,2) NOT NULL DEFAULT 0,
			admin_blocked BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`DROP TRIGGER IF EXISTS trg_admins_updated_at ON admins;`,
		`CREATE TRIGGER trg_admins_updated_at BEFORE UPDATE ON admins
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// ============================================================
		// Master Distributors
		// ============================================================
		`CREATE TABLE IF NOT EXISTS master_distributors (
			master_distributor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_id UUID NOT NULL,
			master_distributor_unique_id TEXT UNIQUE NOT NULL DEFAULT ('MD' || LPAD(nextval('master_distributor_seq')::TEXT, 7, '0')),
			master_distributor_name TEXT NOT NULL,
			master_distributor_phone TEXT UNIQUE NOT NULL,
			master_distributor_email TEXT UNIQUE NOT NULL,
			master_distributor_password TEXT NOT NULL,
			master_distributor_aadhar_number TEXT NOT NULL DEFAULT '',
			master_distributor_pan_number TEXT NOT NULL DEFAULT '',
			master_distributor_date_of_birth TEXT NOT NULL DEFAULT '',
			master_distributor_gender TEXT NOT NULL DEFAULT '',
			master_distributor_city TEXT NOT NULL DEFAULT '',
			master_distributor_state TEXT NOT NULL DEFAULT '',
			master_distributor_address TEXT NOT NULL DEFAULT '',
			master_distributor_pincode TEXT NOT NULL DEFAULT '',
			business_name TEXT NOT NULL DEFAULT '',
			business_type TEXT NOT NULL DEFAULT '',
			gst_number TEXT NOT NULL DEFAULT '',
			master_distributor_wallet_balance NUMERIC(20,2) NOT NULL DEFAULT 0,
			master_distributor_blocked BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE
		);`,
		`DROP TRIGGER IF EXISTS trg_master_distributors_updated_at ON master_distributors;`,
		`CREATE TRIGGER trg_master_distributors_updated_at BEFORE UPDATE ON master_distributors
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// ============================================================
		// Distributors
		// ============================================================
		`CREATE TABLE IF NOT EXISTS distributors (
			distributor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			distributor_unique_id TEXT UNIQUE NOT NULL DEFAULT ('D' || LPAD(nextval('distributor_seq')::TEXT, 7, '0')),
			master_distributor_id UUID NOT NULL,
			admin_id UUID NOT NULL,
			distributor_name TEXT NOT NULL,
			distributor_phone TEXT UNIQUE NOT NULL,
			distributor_email TEXT UNIQUE NOT NULL,
			distributor_password TEXT NOT NULL,
			distributor_aadhar_number TEXT NOT NULL DEFAULT '',
			distributor_pan_number TEXT NOT NULL DEFAULT '',
			distributor_date_of_birth TEXT NOT NULL DEFAULT '',
			distributor_gender TEXT NOT NULL DEFAULT '',
			distributor_city TEXT NOT NULL DEFAULT '',
			distributor_state TEXT NOT NULL DEFAULT '',
			distributor_address TEXT NOT NULL DEFAULT '',
			distributor_pincode TEXT NOT NULL DEFAULT '',
			business_name TEXT NOT NULL DEFAULT '',
			business_type TEXT NOT NULL DEFAULT '',
			gst_number TEXT NOT NULL DEFAULT '',
			distributor_wallet_balance NUMERIC(20,2) NOT NULL DEFAULT 0,
			distributor_blocked BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (master_distributor_id) REFERENCES master_distributors(master_distributor_id) ON DELETE CASCADE,
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE
		);`,
		`DROP TRIGGER IF EXISTS trg_distributors_updated_at ON distributors;`,
		`CREATE TRIGGER trg_distributors_updated_at BEFORE UPDATE ON distributors
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// ============================================================
		// Users
		// ============================================================
		`CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_id UUID NOT NULL,
			master_distributor_id UUID NOT NULL,
			distributor_id UUID NOT NULL,
			user_unique_id TEXT UNIQUE NOT NULL DEFAULT ('R' || LPAD(nextval('user_seq')::TEXT, 7, '0')),
			user_name TEXT NOT NULL,
			user_phone TEXT UNIQUE NOT NULL,
			user_email TEXT UNIQUE NOT NULL,
			user_password TEXT NOT NULL,
			user_aadhar_number TEXT NOT NULL DEFAULT '',
			user_pan_number TEXT NOT NULL DEFAULT '',
			user_date_of_birth TEXT NOT NULL DEFAULT '',
			user_gender TEXT NOT NULL DEFAULT '',
			user_city TEXT NOT NULL DEFAULT '',
			user_state TEXT NOT NULL DEFAULT '',
			user_address TEXT NOT NULL DEFAULT '',
			user_pincode TEXT NOT NULL DEFAULT '',
			business_name TEXT NOT NULL DEFAULT '',
			business_type TEXT NOT NULL DEFAULT '',
			gst_number TEXT NOT NULL DEFAULT '',
			user_mpin TEXT NOT NULL DEFAULT '',
			user_kyc_status BOOLEAN NOT NULL DEFAULT FALSE,
			user_wallet_balance NUMERIC(20,2) NOT NULL DEFAULT 0,
			user_blocked BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE,
			FOREIGN KEY (master_distributor_id) REFERENCES master_distributors(master_distributor_id) ON DELETE CASCADE,
			FOREIGN KEY (distributor_id) REFERENCES distributors(distributor_id) ON DELETE CASCADE
		);`,
		`DROP TRIGGER IF EXISTS trg_users_updated_at ON users;`,
		`CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// ============================================================
		// Payout Service
		// ============================================================
		`CREATE TABLE IF NOT EXISTS payout_service (
			payout_transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			operator_transaction_id TEXT DEFAULT '',
			order_id TEXT DEFAULT '',
			user_id UUID NOT NULL,
			mobile_number TEXT NOT NULL,
			account_number TEXT NOT NULL,
			ifsc_code TEXT NOT NULL,
			bank_name TEXT NOT NULL,
			beneficiary_name TEXT NOT NULL,
			amount NUMERIC(20,2) NOT NULL DEFAULT 0,
			commision NUMERIC(20,2) NOT NULL DEFAULT 0,
			transfer_type TEXT NOT NULL CHECK (transfer_type IN ('IMPS','NEFT')),
			transaction_status TEXT NOT NULL CHECK (transaction_status IN ('PENDING','SUCCESS','FAILED','REFUND')),
			remarks TEXT DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);`,
		`DROP TRIGGER IF EXISTS trg_payout_service_updated_at ON payout_service;`,
		`CREATE TRIGGER trg_payout_service_updated_at BEFORE UPDATE ON payout_service
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// ============================================================
		// Fund Requests
		// ============================================================
		`CREATE TABLE IF NOT EXISTS fund_requests (
			admin_id UUID NOT NULL,
			request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			requester_id UUID NOT NULL,
			requester_name TEXT NOT NULL,
			requester_unique_id TEXT NOT NULL,
			requester_type TEXT NOT NULL CHECK (requester_type IN ('USER','DISTRIBUTOR','MASTER_DISTRIBUTOR')),
			amount NUMERIC(20,2) NOT NULL DEFAULT 0,
			bank_name TEXT NOT NULL,
			request_date TEXT NOT NULL,
			utr_number TEXT UNIQUE NOT NULL,
			request_status TEXT NOT NULL CHECK (request_status IN ('PENDING','APPROVED','REJECTED')),
			remarks TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE
		);`,
		`DROP TRIGGER IF EXISTS trg_fund_requests_updated_at ON fund_requests;`,
		`CREATE TRIGGER trg_fund_requests_updated_at BEFORE UPDATE ON fund_requests
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`CREATE TABLE IF NOT EXISTS revert_history(
			revert_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			unique_id TEXT NOT NULL,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			amount TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// ============================================================
		// OTPs
		// ============================================================
		`CREATE TABLE IF NOT EXISTS otps (
			otp CHAR(4) NOT NULL DEFAULT LPAD((FLOOR(random()*10000))::INT::TEXT, 4, '0'),
			phone TEXT,
			email TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE OR REPLACE FUNCTION delete_expired_otps() RETURNS trigger AS $$
		BEGIN
			DELETE FROM otps WHERE created_at < NOW() - INTERVAL '5 minutes';
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`,
		`DROP TRIGGER IF EXISTS trg_delete_expired_otps ON otps;`,
		`CREATE TRIGGER trg_delete_expired_otps
			AFTER INSERT ON otps
			FOR EACH ROW EXECUTE FUNCTION delete_expired_otps();`,

		// ============================================================
		// Unified Transactions
		// ============================================================
		`CREATE TABLE IF NOT EXISTS transactions (
			transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			transactor_id UUID NOT NULL,
			receiver_id UUID,
			transactor_name TEXT,
			receiver_name TEXT,
			transactor_type TEXT NOT NULL CHECK (transactor_type IN ('ADMIN','MASTER_DISTRIBUTOR','DISTRIBUTOR','USER')),
			receiver_type TEXT CHECK (receiver_type IN ('ADMIN','MASTER_DISTRIBUTOR','DISTRIBUTOR','USER')),
			transaction_type TEXT NOT NULL CHECK (transaction_type IN ('CREDIT','DEBIT')),
			transaction_service TEXT NOT NULL CHECK (transaction_service IN ('FUND_REQUEST','TOPUP','PAYOUT','COMMISSION','FUND_TRANSFER')),
			amount NUMERIC(20,2) NOT NULL CHECK (amount >= 0),
			transaction_status TEXT CHECK (transaction_status IN ('PENDING','SUCCESS','FAILED')),
			remarks TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`DROP TRIGGER IF EXISTS trg_transactions_updated_at ON transactions;`,
		`CREATE TRIGGER trg_transactions_updated_at
			BEFORE UPDATE ON transactions
			FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_transactor
			ON transactions (transactor_type, transactor_id, created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_service
			ON transactions (transaction_service);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_status
			ON transactions (transaction_status);`,
		`CREATE TABLE IF NOT EXISTS banks(
			bank_name TEXT NOT NULL,
			ifsc_code TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS beneficiaries(
			beneficiary_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			mobile_number TEXT NOT NULL,
			bank_name TEXT NOT NULL,
			ifsc_code TEXT NOT NULL,
			account_number TEXT NOT NULL,
			beneficiary_name TEXT NOT NULL,
			beneficiary_phone TEXT NOT NULL,
			beneficiary_verified BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tickets(
			ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_id UUID NOT NULL,
			name TEXT NOT NULL,
			subject TEXT NOT NULL,
			email TEXT NOT NULL,
			phone TEXT NOT NULL,
			message TEXT NOT NULL,
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE
		)`,
	}

	tx, err := qr.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatalf("failed to start SQL transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for i, q := range queries {
		if _, execErr := tx.Exec(ctx, q); execErr != nil {
			log.Fatalf("failed to execute query %d: %v\nSQL:\n%s", i+1, execErr, q)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	log.Println("database initialized successfully")
}
