package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
	dbmodel "github.com/ryakosh/wishlist/lib/db/model"
	"github.com/ryakosh/wishlist/lib/email"
	"github.com/ryakosh/wishlist/lib/graph/generated"
	"github.com/ryakosh/wishlist/lib/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	var user dbmodel.User

	err := lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Where("id = ?", input.ID).Or("email = ?", input.Email).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if !d.RecordNotFound() {
		return nil, dbmodel.ErrUserExists
	}

	user = dbmodel.User{
		ID:        input.ID,
		Email:     input.Email,
		Password:  dbmodel.GenPasswordHash(input.Password),
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	d = r.DB.Create(&user)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not create user", d.Error)
	}

	code, err := dbmodel.CreateCode(input.ID)
	if err != nil {
		se, ok := err.(*dbmodel.ServerError)
		if ok {
			lib.LogError(lib.LError, "Could not generate email confirmation mail", se.Reason)
			return nil, email.ErrSendMail
		}

		return nil, err
	}

	mail, err := email.GenEmailConfirmMail(input.ID, code.View.(string))
	if err != nil {
		lib.LogError(lib.LError, "Could not generate email confirmation mail", err)
		return nil, email.ErrSendMail
	}

	err = email.Send(email.BotEmailEnv, input.Email, "لطفا ایمیل خود را تایید کنید [ویش لیست]", mail)
	if err != nil {
		lib.LogError(lib.LError, "Could not generate email confirmation mail", err)
		return nil, email.ErrSendMail
	}

	return &model.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Friends:        user.ID,
		FriendRequests: user.ID,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Updates(&dbmodel.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not update user", d.Error)
	}

	return &model.User{
		ID:             authedUser,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Friends:        authedUser,
		FriendRequests: authedUser,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (string, error) {
	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return "", err
	}

	d := r.DB.Delete(&dbmodel.User{ID: authedUser})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete user", d.Error)
	}

	return authedUser, nil
}

func (r *mutationResolver) GenToken(ctx context.Context, input model.Login) (string, error) {
	var user dbmodel.User

	err := lib.Validator.Struct(&input)
	if err != nil {
		return "", lib.ErrValidationFailed
	}

	d := r.DB.Select("id, email, password").Where("id = ?", input.ID).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() || !dbmodel.VerifyPassword(input.Password, user.Password) {
		return "", dbmodel.ErrUnmOrPwdIncorrect
	}

	return lib.Encode(user.ID, user.Email), nil
}

func (r *mutationResolver) VerifyEmail(ctx context.Context, code string) (bool, error) {
	var user dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return false, err
	}

	err = lib.Validator.Var(code, "max=14")
	if err != nil {
		return false, lib.ErrValidationFailed
	}

	d := r.DB.Select("is_email_verified").Where("id = ?", authedUser).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() {
		return false, dbmodel.ErrUserNotFound
	}

	if user.IsEmailVerified {
		return false, dbmodel.ErrEmailVerified
	}

	isMatch, err := dbmodel.VerifyCode(authedUser, code)
	if err != nil {
		return false, err
	}

	if isMatch.View.(bool) {
		d := db.DB.Model(&dbmodel.User{ID: authedUser}).Update("is_email_verified", true)
		if d.Error != nil {
			lib.LogError(lib.LPanic, "Could not update user", d.Error)
		}
	}

	return true, nil
}

func (r *mutationResolver) SendFriendRequest(ctx context.Context, id string) (*model.User, error) {
	var requestee dbmodel.User
	var friendsCount uint8
	var friendRequestsCount uint8

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "username,max=64")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	if authedUser == id {
		return nil, dbmodel.ErrUserNotFound
	}

	d := r.DB.Select("id, first_name, last_name").Where("id = ?", id).First(&requestee)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrUserNotFound
	}

	d = r.DB.Table("friendrequests").Where("user_id = ? AND requester_id = ?", id, authedUser).Count(&friendRequestsCount)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	}

	d = r.DB.Table("friendships").Where("user_id = ? AND friend_id = ?", authedUser, id).Count(&friendsCount)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	}

	if friendRequestsCount != 0 || friendsCount != 0 {
		return nil, dbmodel.ErrUserExists
	}

	err = r.DB.Model(&dbmodel.User{ID: id}).Association("FriendRequests").Append(&dbmodel.User{ID: authedUser}).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not request friendship", err)
	}

	return &model.User{
		ID:             requestee.ID,
		FirstName:      requestee.FirstName,
		LastName:       requestee.LastName,
		Friends:        requestee.ID,
		FriendRequests: requestee.ID,
	}, nil
}

func (r *mutationResolver) UnSendFriendRequest(ctx context.Context, id string) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "username,max=64")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	err = r.DB.Model(&dbmodel.User{ID: id}).Select("id, first_name, last_name").Where("requester_id = ?", authedUser).Association("FriendRequests").Find(&requestees).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, dbmodel.ErrUserNotFound
	} else if err != nil {
		lib.LogError(lib.LPanic, "Could not delete friendship request", err)
	}

	err = r.DB.Model(&dbmodel.User{ID: id}).Association("FriendRequests").Delete(&dbmodel.User{ID: authedUser}).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not delete friendship request", err)
	}

	return &model.User{
		ID:             requestees[0].ID,
		FirstName:      requestees[0].FirstName,
		LastName:       requestees[0].LastName,
		Friends:        requestees[0].ID,
		FriendRequests: requestees[0].ID,
	}, nil
}

func (r *mutationResolver) AcceptFriendRequest(ctx context.Context, id string) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "username,max=64")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select("id, first_name, last_name").Where(
		"requester_id = ?", id).Related(&requestees, "FriendRequests")
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user's friend requests", d.Error)
	}

	if len(requestees) != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&dbmodel.User{ID: authedUser}).Association("Friends").Append(requestees[0]).Error
		if err != nil {
			return err
		}

		err = tx.Model(&dbmodel.User{ID: requestees[0].ID}).Association("Friends").Append(&dbmodel.User{ID: authedUser}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&dbmodel.User{ID: authedUser}).Association("FriendRequests").Delete(requestees[0]).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		lib.LogError(lib.LPanic, " Could not accept friendship", err)
	}

	return &model.User{
		ID:             requestees[0].ID,
		FirstName:      requestees[0].FirstName,
		LastName:       requestees[0].LastName,
		Friends:        requestees[0].ID,
		FriendRequests: requestees[0].ID,
	}, nil
}

func (r *mutationResolver) RejectFriendRequest(ctx context.Context, id string) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "username,max=64")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select("id, first_name, last_name").Where(
		"requester_id = ?", id).Related(&requestees, "FriendRequests")
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user's friend requests", d.Error)
	}

	if len(requestees) != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err = r.DB.Model(&dbmodel.User{ID: authedUser}).Association("FriendRequests").Delete(requestees[0]).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not reject friendship", err)
	}

	return &model.User{
		ID:             requestees[0].ID,
		FirstName:      requestees[0].FirstName,
		LastName:       requestees[0].LastName,
		Friends:        requestees[0].ID,
		FriendRequests: requestees[0].ID,
	}, nil
}

func (r *mutationResolver) CreateWish(ctx context.Context, input model.NewWish) (*model.Wish, error) {
	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	wish := dbmodel.Wish{
		Owner:       authedUser,
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		Image:       input.Image,
	}

	d := r.DB.Create(&wish)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not create wish", d.Error)
	}

	return &model.Wish{
		ID:                  wish.ID,
		Owner:               authedUser,
		Name:                wish.Name,
		Description:         wish.Description,
		Link:                wish.Link,
		Image:               wish.Image,
		FulfillmentClaimers: wish.ID,
		Fulfillers:          wish.ID,
	}, nil
}

func (r *mutationResolver) UpdateWish(ctx context.Context, input model.UpdateWish) (*model.Wish, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, name, owner, description, link, image").First(&wish, input.ID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if wish.Owner != authedUser {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	d = r.DB.Model(&wish).Updates(&dbmodel.Wish{
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		Image:       input.Image,
	})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not update wish", d.Error)
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

func (r *mutationResolver) DeleteWish(ctx context.Context, id int) (int, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return 0, err
	}

	err = lib.Validator.Var(id, "min=0")
	if err != nil {
		return 0, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, owner").First(&wish, id)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return 0, dbmodel.ErrWishNotFound
	}
	if wish.Owner != authedUser {
		return 0, dbmodel.ErrUserNotAuthorized
	}

	d = r.DB.Delete(wish)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete wish", d.Error)
	}

	return wish.ID, nil
}

func (r *mutationResolver) AddWantToFulfill(ctx context.Context, id int) (*model.Wish, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "min=0")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, name, owner, description, link, image").First(&wish, id)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser == wish.Owner || !dbmodel.AreFriends(wish.Owner, authedUser) {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	asso := r.DB.Model(&dbmodel.Wish{ID: id}).Where("user_id = ?", authedUser).Association("WantToFulfill")

	if asso.Count() != 0 {
		return nil, dbmodel.ErrUserExists
	}

	err = r.DB.Model(&dbmodel.Wish{ID: id}).Association("WantToFulfill").Append(&dbmodel.User{ID: authedUser}).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not add to WantToFulfill", err)
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

func (r *mutationResolver) ClaimFulfillment(ctx context.Context, id int) (*model.Wish, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Var(id, "min=0")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, name, owner, description, link, image").First(&wish, id)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	asso := r.DB.Model(&dbmodel.Wish{ID: id}).Where("user_id = ?", authedUser).Association("WantToFulfill")
	if asso.Error != nil && !gorm.IsRecordNotFoundError(asso.Error) {
		lib.LogError(lib.LPanic, "Could not read wish's WantToFulfill", asso.Error)
	}

	if asso.Count() != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		asso := tx.Model(&dbmodel.Wish{ID: id}).Association("Claimers").Append(&dbmodel.User{ID: authedUser})
		if asso.Error != nil {
			return asso.Error
		}

		asso = tx.Model(&dbmodel.Wish{ID: id}).Association("WantToFulfill").Delete(&dbmodel.User{ID: authedUser})
		if asso.Error != nil {
			return asso.Error
		}

		return nil
	})
	if err != nil {
		lib.LogError(lib.LPanic, "Could not add to Claimers", err)
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

func (r *mutationResolver) AcceptFulfillmentClaim(ctx context.Context, input model.FulfillmentClaimer) (*model.Wish, error) {
	err := lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.handleClaimer(ctx, input.WishID, input.ClaimerID, dbmodel.WishFulFillersAsso)
}

func (r *mutationResolver) RejectFulfillmentClaim(ctx context.Context, input model.FulfillmentClaimer) (*model.Wish, error) {
	err := lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.handleClaimer(ctx, input.WishID, input.ClaimerID, dbmodel.WishWantToFulfillAsso)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	err := lib.Validator.Var(id, "username,max=64")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.user(ctx, id)
}

func (r *queryResolver) Wish(ctx context.Context, id int) (*model.Wish, error) {
	err := lib.Validator.Var(id, "min=0")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	return r.wish(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
