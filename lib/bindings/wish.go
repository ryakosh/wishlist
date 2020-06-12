package bindings

// CWish is used to create Wish models
type CWish struct {
	Name        string `json:"name" binding:"required,max=256"`
	Description string `json:"description" binding:"max=1024"`
	Link        string `json:"link" binding:"url"`
	Image       string `json:"image" binding:"url"`
}

// RdWish is used to read and delete Wish models
type RdWish struct {
	ID uint `json:"id" binding:"required"`
}

// UWish is used to update Wish models
type UWish struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"max=256"`
	Description string `json:"description" binding:"max=1024"`
	Link        string `json:"link" binding:"url"`
	Image       string `json:"image" binding:"url"`
	FulfilledBy string `json:"fulfilled_by" binding:"username,max=64"`
}
