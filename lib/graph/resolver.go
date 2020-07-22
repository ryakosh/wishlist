package graph

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
	dbmodel "github.com/ryakosh/wishlist/lib/db/model"
	"github.com/ryakosh/wishlist/lib/graph/model"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	DB *gorm.DB
}

func (r *Resolver) handleClaimer(ctx context.Context, wishID int,
	claimer string, appendTo db.Association) (*model.Wish, error) {
	var wish dbmodel.Wish

	authedUser := dbmodel.AuthedUserFromCtx(ctx)

	d := r.DB.Select("id, name, owner, description, link, image").First(&wish, wishID)
	if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.Owner {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	count := r.DB.Model(&wish).Where("user_id = ?", claimer).Association("Claimers").Count()
	if count != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		asso := tx.Model(&wish).Association(string(appendTo)).Append(&dbmodel.User{ID: claimer})
		if asso.Error != nil {
			return asso.Error
		}

		asso = tx.Model(&wish).Association(string(dbmodel.WishClaimersAsso)).Delete(&dbmodel.User{ID: claimer})
		if asso.Error != nil {
			return asso.Error
		}

		return nil
	})
	if err != nil {
		lib.LogError(lib.LPanic, "Could not accept fulfillment claim", err)
	}

	return &model.Wish{
		ID:                  wish.ID,
		Owner:               wish.Owner,
		Name:                wish.Name,
		Description:         wish.Description,
		Link:                wish.Link,
		Image:               wish.Image,
		FulfillmentClaimers: wish.ID,
		Fulfillers:          wish.ID,
	}, nil
}

func (r *Resolver) user(ctx context.Context, id string) (*model.User, error) {
	var user dbmodel.User

	d := r.DB.Select("id, first_name, last_name").Where("id = ?", id).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrUserNotFound
	}

	return &model.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Friends:        user.ID,
		FriendRequests: user.ID,
	}, nil
}

func (r *Resolver) wish(ctx context.Context, wishID int) (*model.Wish, error) {
	var wish dbmodel.Wish

	d := r.DB.Omit("created_at, updated_at").First(&wish, wishID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}
	return &model.Wish{
		ID:                  wish.ID,
		Owner:               wish.Owner,
		Name:                wish.Name,
		Description:         wish.Description,
		Link:                wish.Link,
		Image:               wish.Image,
		FulfillmentClaimers: wish.ID,
		Fulfillers:          wish.ID,
	}, nil
}
