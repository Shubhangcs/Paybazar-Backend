package queries

import (
	"context"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetAllMasterDistributorsByID(adminId string) (*[]structures.MasterDistributorGetResponse, error) {
	const query = `
	SELECT
		md.master_distributor_unique_id,
		md.master_distributor_name,
		md.master_distributor_email,
		md.master_distributor_phone,
		COALESCE(mdw.balance::TEXT, '0') AS master_distributor_wallet_balance
	FROM master_distributors md
	LEFT JOIN master_distributor_wallets mdw
		ON mdw.master_distributor_id = md.master_distributor_id
	WHERE md.admin_id = $1
	ORDER BY md.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []structures.MasterDistributorGetResponse
	for rows.Next() {
		var r structures.MasterDistributorGetResponse
		if err := rows.Scan(
			&r.MasterDistributorUniqueID,
			&r.MasterDistributorName,
			&r.MasterDistributorEmail,
			&r.MasterDistributorPhone,
			&r.MasterDistributorWalletBalance,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &out, nil
}

func (q *Query) GetAllDistributorsByMasterDistributorID(masterDistributorId string) (*[]structures.DistributorGetResponse, error) {
	const query = `
	SELECT
		d.distributor_unique_id,
		d.distributor_name,
		d.distributor_email,
		d.distributor_phone,
		COALESCE(dw.balance::TEXT, '0') AS distributor_wallet_balance
	FROM distributors d
	LEFT JOIN distributor_wallets dw
		ON dw.distributor_id = d.distributor_id
	WHERE d.master_distributor_id = $1
	ORDER BY d.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, masterDistributorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []structures.DistributorGetResponse
	for rows.Next() {
		var r structures.DistributorGetResponse
		if err := rows.Scan(
			&r.DistributorUniqueID,
			&r.DistributorName,
			&r.DistributorEmail,
			&r.DistributorPhone,
			&r.DistributorWalletBalance,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &out, nil
}

func (q *Query) GetAllUsersByDistributorID(distributorId string) (*[]structures.UserGetResponse, error) {
	const query = `
	SELECT
		u.user_unique_id,
		u.user_name,
		u.user_email,
		u.user_phone,
		COALESCE(uw.balance::TEXT, '0') AS user_wallet_balance
	FROM users u
	LEFT JOIN user_wallets uw
		ON uw.user_id = u.user_id
	WHERE u.distributor_id = $1
	ORDER BY u.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := q.Pool.Query(ctx, query, distributorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []structures.UserGetResponse
	for rows.Next() {
		var r structures.UserGetResponse
		if err := rows.Scan(
			&r.UserUniqueID,
			&r.UserName,
			&r.UserEmail,
			&r.UserPhone,
			&r.UserWalletBalance,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &out, nil
}
