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

func (r *wishResolver) Owner(ctx context.Context, obj *model.Wish) (*model.User, error) {
	return r.user(ctx, obj.Owner)
}

func (r *wishResolver) FulfillmentClaimers(ctx context.Context, obj *model.Wish) (*model.Users, error) {
	return &model.Users{
		InObj:         obj,
		InAssociation: dbmodel.WishClaimersAsso,
	}, nil
}

func (r *wishResolver) Fulfillers(ctx context.Context, obj *model.Wish) (*model.Users, error) {
	return &model.Users{
		InObj:         obj,
		InAssociation: dbmodel.WishFulFillersAsso,
	}, nil
}

func (r *wishesResolver) Query(ctx context.Context, obj *model.Wishes, page int, limit int) ([]*model.Wish, error) {
	var wishes []dbmodel.Wish
	var res []*model.Wish

	err := lib.Validator.Struct(struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=10"`
	}{Page: page, Limit: limit})
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: obj.InObj.ID}).Select(
		"id, name, owner, description, link, image").Offset(
		(page * limit) - limit).Limit(limit).Association(string(dbmodel.UserWishesAsso)).Find(&wishes)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user's friends", d.Error)
	}

	for _, w := range wishes {
		res = append(res, &model.Wish{
			ID:                  w.ID,
			Owner:               w.Owner,
			Name:                w.Name,
			Description:         w.Description,
			Link:                w.Link,
			Image:               w.Image,
			FulfillmentClaimers: w.ID,
			Fulfillers:          w.ID,
		})
	}

	return res, nil
}

func (r *wishesResolver) Count(ctx context.Context, obj *model.Wishes) (int, error) {
	return r.DB.Model(&dbmodel.User{ID: obj.InObj.ID}).Association(string(dbmodel.UserWishesAsso)).Count(), nil
}

// Wish returns generated.WishResolver implementation.
func (r *Resolver) Wish() generated.WishResolver { return &wishResolver{r} }

// Wishes returns generated.WishesResolver implementation.
func (r *Resolver) Wishes() generated.WishesResolver { return &wishesResolver{r} }

type wishResolver struct{ *Resolver }
type wishesResolver struct{ *Resolver }
