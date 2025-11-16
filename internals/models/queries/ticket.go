package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) AddNewTicket(req structures.Ticket) error {
	query := `INSERT INTO tickets(
		admin_id,
		name,
		subject,
		phone,
		email,
		message
	) VALUES(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6 
	)`
	_, err := q.Pool.Exec(context.Background(), query, req.AdminID, req.Name, req.Subject, req.Mobile, req.Email, req.Message)
	return err
}

func (q *Query) GetAllTickets(adminId string) (*[]structures.Ticket, error) {
	query := `SELECT admin_id, name, subject, phone, email, message FROM tickets WHERE admin_id=$1`

	res, err := q.Pool.Query(context.Background(), query, adminId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var tickets []structures.Ticket
	for res.Next() {
		var ticket structures.Ticket
		if err := res.Scan(
			&ticket.AdminID,
			&ticket.Name,
			&ticket.Subject,
			&ticket.Mobile,
			&ticket.Email,
			&ticket.Message,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &tickets, nil
}
