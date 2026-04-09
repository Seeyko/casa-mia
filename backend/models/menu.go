package models

type MenuCategory struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Section   string     `json:"section"` // "carte" or "traiteur"
	SortOrder int        `json:"sort_order"`
	Items     []MenuItem `json:"items,omitempty"`
}

type MenuItem struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       *string `json:"price"` // nullable, string to handle "2€/pièce" etc.
	ImagePath   string  `json:"image_path"`
	SortOrder   int     `json:"sort_order"`
	Available   bool    `json:"available"`
	Badge       string  `json:"badge"` // "NEW", "★", etc.
	Note        string  `json:"note"`  // extra note like seasonal info
}

type MenuResponse struct {
	Categories []MenuCategory `json:"categories"`
}
