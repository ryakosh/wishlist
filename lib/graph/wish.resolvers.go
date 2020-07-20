package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
	dbmodel "github.com/ryakosh/wishlist/lib/db/model"
	"github.com/ryakosh/wishlist/lib/graph/generated"
	"github.com/ryakosh/wishlist/lib/graph/model"
)

func (r *wishResolver) Owner(ctx context.Context, obj *model.Wish) (*model.User, error) {
	return r.user(ctx, obj.Owner)
}

func (r *wishResolver) FulfillmentClaimers(ctx context.Context, obj *model.Wish, page int, limit int) ([]*model.User, error) {
	err := lib.Validator.Struct(struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=10"`
	}{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.wishUsers(ctx, obj.FulfillmentClaimers, dbmodel.WishClaimersAsso, page, limit)
}

func (r *wishResolver) Fulfillers(ctx context.Context, obj *model.Wish, page int, limit int) ([]*model.User, error) {
	err := lib.Validator.Struct(struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=10"`
	}{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.wishUsers(ctx, obj.Fulfillers, dbmodel.WishFulFillersAsso, page, limit)
}

func (r *Resolver) wishUsers(ctx context.Context, wishID int, asso db.Association, page int, limit int) ([]*model.User, error) {
	var wish dbmodel.Wish
	var users []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	d := r.DB.Select("id, owner").First(&wish, wishID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.Owner {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	r.DB.Model(&wish).Select("id, first_name, last_name").Offset(
		(page * limit) - limit).Limit(limit).Association(string(asso)).Find(&users)

	for _, u := range users {
		res = append(res, &model.User{
			ID:             u.ID,
			FirstName:      u.FirstName,
			LastName:       u.LastName,
			Friends:        u.ID,
			FriendRequests: u.ID,
		})
	}

	return res, nil
}

// Wish returns generated.WishResolver implementation.
func (r *Resolver) Wish() generated.WishResolver { return &wishResolver{r} }

type wishResolver struct{ *Resolver }
