package models

import (
	"encoding/json"
	"time"
)

type OpeningSlot struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type DayHours struct {
	Slots []OpeningSlot `json:"slots"`
}

// OpeningHours maps day names (lundi, mardi, ...) to their hours
type OpeningHours map[string]*DayHours

type Location struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Address      string         `json:"address"`
	Phone        string         `json:"phone"`
	OpeningHours json.RawMessage `json:"opening_hours"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type LocationStatus struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	IsOpen     bool   `json:"is_open"`
	NextChange string `json:"next_change"` // e.g. "Ouvre mardi à 9h" or "Ferme à 21h"
}
