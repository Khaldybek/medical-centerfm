package migrations

import (
	"gorm.io/gorm"
)

type CreateUsersTable struct{}

func (m *CreateUsersTable) ID() string {
	return "000005_create_users"
}

func (m *CreateUsersTable) Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'patient',
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`).Error
}

func (m *CreateUsersTable) Rollback(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS users").Error
} 