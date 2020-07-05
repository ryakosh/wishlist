package views

// CWish is used to respond when a new wish is created for a user
type CWish struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
}

// RWish is used to respond about a wish's general information
type RWish struct {
	ID          uint64 `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
}

// UWish is used to respond when a wish's general information gets updated
type UWish struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
}

type WishID struct {
	ID uint64 `json:"id"`
}
