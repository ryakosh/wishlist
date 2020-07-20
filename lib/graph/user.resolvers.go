package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	dbmodel "github.com/ryakosh/wishlist/lib/db/model"
	"github.com/ryakosh/wishlist/lib/graph/generated"
	"github.com/ryakosh/wishlist/lib/graph/model"
)

func (r *userResolver) Friends(ctx context.Context, obj *model.User, page int, limit int) ([]*model.User, error) {
	var friends []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	if authedUser != obj.ID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	err = lib.Validator.Struct(struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=10"`
	}{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(page * limit) - limit).Limit(limit).Association("Friends").Find(&friends)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user's friends", d.Error)
	}

	for _, f := range friends {
		res = append(res, &model.User{
			ID:             f.ID,
			FirstName:      f.FirstName,
			LastName:       f.LastName,
			Friends:        f.ID,
			FriendRequests: f.ID,
		})
	}

	return res, nil
}

func (r *userResolver) FriendRequests(ctx context.Context, obj *model.User, page int, limit int) ([]*model.User, error) {
	var reqs []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	if authedUser != obj.ID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	err = lib.Validator.Struct(struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=10"`
	}{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(page * limit) - limit).Limit(limit).Association("FriendRequests").Find(&reqs)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user's friend requests", d.Error)
	}

	for _, r := range reqs {
		res = append(res, &model.User{
			ID:             r.ID,
			FirstName:      r.FirstName,
			LastName:       r.LastName,
			Friends:        r.ID,
			FriendRequests: r.ID,
		})
	}

	return res, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
