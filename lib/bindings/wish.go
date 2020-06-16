package bindings

// CWish is used to create a new wish for a user
type CWish struct {
	Name        string `json:"name" binding:"required,max=256"`
	Description string `json:"description" binding:"max=1024"`
	Link        string `json:"link" binding:"url"`
	Image       string `json:"image" binding:"url"`
}

// RdWish is used to read wish's general information and delete a wish
type RdWish struct {
	ID uint `json:"id" binding:"required"`
}

// UWish is used to update wish's general information
type UWish struct {
	Name        string `json:"name" binding:"max=256"`
	Description string `json:"description" binding:"max=1024"`
	Link        string `json:"link" binding:"url"`
	Image       string `json:"image" binding:"url"`
}
