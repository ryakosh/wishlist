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

func (r *userResolver) Wishes(ctx context.Context, obj *model.User) (*model.Wishes, error) {
	return &model.Wishes{
		InObj: obj,
	}, nil
}

func (r *userResolver) Friends(ctx context.Context, obj *model.User) (*model.Users, error) {
	return &model.Users{
		InObj:         obj,
		InAssociation: dbmodel.UserFriendsAsso,
	}, nil
}

func (r *userResolver) FriendRequests(ctx context.Context, obj *model.User) (*model.Users, error) {
	return &model.Users{
		InObj:         obj,
		InAssociation: dbmodel.UserFriendRequestsAsso,
	}, nil
}

func (r *usersResolver) Query(ctx context.Context, obj *model.Users, page int, limit int) ([]*model.User, error) {
	var users []dbmodel.User
	var res []*model.User
	var d *gorm.DB

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
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

	switch o := obj.InObj.(type) {
	case *model.User:
		if authedUser != o.ID {
			return nil, dbmodel.ErrUserNotAuthorized
		}

		d = r.DB.Model(&dbmodel.User{ID: authedUser})
	case *model.Wish:
		if authedUser != o.Owner {
			return nil, dbmodel.ErrUserNotAuthorized
		}

		d = r.DB.Model(&dbmodel.Wish{ID: o.ID})
	default:
		lib.LogError(lib.LPanic, "Model object type assertion failed", nil)
	}

	d.Select("id, first_name, last_name").Offset(
		(page * limit) - limit).Limit(limit).Association(string(obj.InAssociation)).Find(&users)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read users", d.Error)
	}

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

func (r *usersResolver) Count(ctx context.Context, obj *model.Users) (int, error) {
	var d *gorm.DB

	switch o := obj.InObj.(type) {
	case *model.User:
		d = r.DB.Model(&dbmodel.User{ID: o.ID})
	case *model.Wish:
		d = r.DB.Model(&dbmodel.Wish{ID: o.ID})

	default:
		lib.LogError(lib.LPanic, "Model object type assertion failed", nil)
	}

	return d.Association(string(obj.InAssociation)).Count(), nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

// Users returns generated.UsersResolver implementation.
func (r *Resolver) Users() generated.UsersResolver { return &usersResolver{r} }

type userResolver struct{ *Resolver }
type usersResolver struct{ *Resolver }
