package model

type Wish struct {
	ID          int    `json:"id"`
	User        string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Claimers    int    `json:"claimers"`
	Fulfillers  int    `json:"fulfillers"`
}

type NewWish struct {
	Name        string  `json:"name" validate:"min=1,max=256"`
	Description *string `json:"description" validate:"omitempty,max=1024"`
	Link        *string `json:"link" validate:"omitempty,url"`
	Image       *string `json:"image" validate:"omitempty,url"`
}

type UpdateWish struct {
	ID          int     `json:"id" validate:"min=0"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=256"`
	Description *string `json:"description" validate:"omitempty,max=1024"`
	Link        *string `json:"link" validate:"omitempty,url"`
	Image       *string `json:"image" validate:"omitempty,url"`
}
