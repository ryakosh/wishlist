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

func (r *wishResolver) User(ctx context.Context, obj *model.Wish) (*model.User, error) {
	return r.user(ctx, obj.User)
}

func (r *wishResolver) Claimers(ctx context.Context, obj *model.Wish, page int, limit int) ([]*model.User, error) {
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

	return r.wishUsers(ctx, obj.Claimers, dbmodel.WishClaimersAsso, page, limit)
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

// Wish returns generated.WishResolver implementation.
func (r *Resolver) Wish() generated.WishResolver { return &wishResolver{r} }

type wishResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *Resolver) wishUsers(ctx context.Context, wishID int, asso db.Association, page int, limit int) ([]*model.User, error) {
	var wish dbmodel.Wish
	var users []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	d := r.DB.Select("id, user_id").First(&wish, wishID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.UserID {
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
