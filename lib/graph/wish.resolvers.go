package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ryakosh/wishlist/lib"
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

// Wish returns generated.WishResolver implementation.
func (r *Resolver) Wish() generated.WishResolver { return &wishResolver{r} }

type wishResolver struct{ *Resolver }
