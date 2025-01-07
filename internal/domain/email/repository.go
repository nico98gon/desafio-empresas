package email

import (
	"database/sql"
	"fmt"
)

type EmailRepository struct {
	db *sql.DB
}

func NewEmailRepository(db *sql.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) Save(email EmailData) error {
	query := `
		INSERT INTO emails (message_id, date, sender, recipient, subject)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, email.MessageID, email.Date, email.From, email.To, email.Subject)
	if err != nil {
		return fmt.Errorf("error saving email: %w", err)
	}
	return nil
}