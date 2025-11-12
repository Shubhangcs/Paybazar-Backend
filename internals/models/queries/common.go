package queries

import (
	"context"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetAllMasterDistributorsByID(adminId string) (*[]structures.MasterDistributorGetResponse, error) {
	const query = `
		SELECT 
			master_distributor_id,
			master_distributor_unique_id,
			master_distributor_name,
			master_distributor_email,
			master_distributor_phone,
			master_distributor_wallet_balance
		FROM 
			master_distributors
		WHERE 
			admin_id = $1
		ORDER BY 
			created_at DESC;
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
			&r.MasterDistributorID,
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &out, nil
}

func (q *Query) GetAllDistributorsByMasterDistributorID(masterDistributorId string) (*[]structures.DistributorGetResponse, error) {
	const query = `
		SELECT
			distributor_id,
			distributor_unique_id,
			distributor_name,
			distributor_email,
			distributor_phone,
			distributor_wallet_balance
		FROM
			distributors
		WHERE
			master_distributor_id = $1
		ORDER BY
			created_at DESC;
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
			&r.DistributorID,
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &out, nil
}

func (q *Query) GetAllUsersByDistributorID(distributorId string) (*[]structures.UserGetResponse, error) {
	const query = `
		SELECT
			user_id,
			user_unique_id,
			user_name,
			user_email,
			user_phone,
			user_wallet_balance
		FROM
			users
		WHERE
			distributor_id = $1
		ORDER BY
			created_at DESC;
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
			&r.UserID,
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &out, nil
}
