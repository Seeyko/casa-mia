package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.migrate(); err != nil {
		return nil, err
	}

	log.Println("PostgreSQL database initialized")
	return database, nil
}

func (d *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS locations (
		id             SERIAL PRIMARY KEY,
		name           TEXT NOT NULL,
		slug           TEXT UNIQUE NOT NULL,
		address        TEXT NOT NULL,
		phone          TEXT NOT NULL,
		opening_hours  JSONB NOT NULL DEFAULT '{}',
		order_method   TEXT NOT NULL DEFAULT 'phone',
		order_info     TEXT NOT NULL DEFAULT '',
		created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS news (
		id         SERIAL PRIMARY KEY,
		title      TEXT NOT NULL,
		content    TEXT NOT NULL DEFAULT '',
		image_path TEXT DEFAULT '',
		published  BOOLEAN NOT NULL DEFAULT false,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS menu_categories (
		id         SERIAL PRIMARY KEY,
		name       TEXT NOT NULL,
		section    TEXT NOT NULL DEFAULT 'carte',
		sort_order INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS menu_items (
		id          SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES menu_categories(id) ON DELETE CASCADE,
		name        TEXT NOT NULL,
		description TEXT DEFAULT '',
		price       TEXT DEFAULT '',
		image_path  TEXT DEFAULT '',
		sort_order  INTEGER NOT NULL DEFAULT 0,
		available   BOOLEAN NOT NULL DEFAULT true,
		badge       TEXT DEFAULT '',
		note        TEXT DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS admin_users (
		id             SERIAL PRIMARY KEY,
		username       TEXT UNIQUE NOT NULL,
		password_hash  TEXT NOT NULL,
		reset_token    TEXT DEFAULT '',
		reset_expires  TIMESTAMPTZ,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_menu_items_category ON menu_items(category_id);
	CREATE INDEX IF NOT EXISTS idx_menu_categories_section ON menu_categories(section);
	CREATE INDEX IF NOT EXISTS idx_news_published ON news(published);

	ALTER TABLE locations ADD COLUMN IF NOT EXISTS closure_start   DATE;
	ALTER TABLE locations ADD COLUMN IF NOT EXISTS closure_end     DATE;
	ALTER TABLE locations ADD COLUMN IF NOT EXISTS closure_message TEXT NOT NULL DEFAULT '';
	`

	_, err := d.DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database schema migration completed")
	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
