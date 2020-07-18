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

func (r *mutationResolver) GenToken(ctx context.Context, input model.Login) (*model.Token, error) {
	var user dbmodel.User

	err := lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, email, password").Where("id = ?", input.ID).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() || !dbmodel.VerifyPassword(input.Password, user.Password) {
		return nil, dbmodel.ErrUnmOrPwdIncorrect
	}

	return &model.Token{
		Token: lib.Encode(user.ID, user.Email),
	}, nil
}

func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerificationCode) (bool, error) {
	var user dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return false, err
	}

	err = lib.Validator.Struct(&input)
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

	isMatch, err := dbmodel.VerifyCode(authedUser, input.Code)
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

func (r *mutationResolver) RequestFriendship(ctx context.Context, input model.UserID) (*model.User, error) {
	var requestee dbmodel.User
	var friendsCount uint8
	var friendRequestsCount uint8

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	if authedUser == input.ID {
		return nil, dbmodel.ErrUserNotFound
	}

	d := r.DB.Select("id, first_name, last_name").Where("id = ?", input.ID).First(&requestee)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrUserNotFound
	}

	d = r.DB.Table("friendrequests").Where("user_id = ? AND requester_id = ?", input.ID, authedUser).Count(&friendRequestsCount)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	}

	d = r.DB.Table("friendships").Where("user_id = ? AND friend_id = ?", authedUser, input.ID).Count(&friendsCount)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	}

	if friendRequestsCount != 0 || friendsCount != 0 {
		return nil, dbmodel.ErrUserExists
	}

	err = r.DB.Model(&dbmodel.User{ID: input.ID}).Association("FriendRequests").Append(&dbmodel.User{ID: authedUser}).Error
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

func (r *mutationResolver) UnRequestFriendship(ctx context.Context, input model.UserID) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	err = r.DB.Model(&dbmodel.User{ID: input.ID}).Where("requester_id = ?", authedUser).Association("FriendRequests").Find(&requestees).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, dbmodel.ErrUserNotFound
	} else if err != nil {
		lib.LogError(lib.LPanic, "Could not delete friendship request", err)
	}

	err = r.DB.Model(&dbmodel.User{ID: input.ID}).Association("FriendRequests").Delete(&dbmodel.User{ID: authedUser}).Error
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

func (r *mutationResolver) AcceptFriendRequest(ctx context.Context, input model.UserID) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select("id").Where(
		"requester_id = ?", input.ID).Related(&requestees, "FriendRequests")
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

func (r *mutationResolver) RejectFriendshipRequest(ctx context.Context, input model.UserID) (*model.User, error) {
	var requestees []dbmodel.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select("id").Where(
		"requester_id = ?", input.ID).Related(&requestees, "FriendRequests")
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
		UserID:      authedUser,
		Name:        input.Name,
		Description: lib.DummyDefault(input.Description),
		Link:        lib.DummyDefault(input.Link),
		Image:       lib.DummyDefault(input.Image),
	}

	d := r.DB.Create(&wish)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not create wish", d.Error)
	}

	return &model.Wish{
		ID:          wish.ID,
		User:        authedUser,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
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

	d := r.DB.Select("id, name, user_id, description, link, image").First(&wish, input.ID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if wish.UserID != authedUser {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	d = r.DB.Model(&wish).Updates(&dbmodel.Wish{
		Name:        lib.DummyDefault(input.Name),
		Description: lib.DummyDefault(input.Description),
		Link:        lib.DummyDefault(input.Link),
		Image:       lib.DummyDefault(input.Image),
	})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not update wish", d.Error)
	}

	return &model.Wish{
		ID:          wish.ID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
	}, nil
}

func (r *mutationResolver) DeleteWish(ctx context.Context, input model.WishID) (int, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return 0, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return 0, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, user_id").First(&wish, input.ID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return 0, dbmodel.ErrWishNotFound
	}
	if wish.UserID != authedUser {
		return 0, dbmodel.ErrUserNotAuthorized
	}

	d = r.DB.Delete(wish)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete wish", d.Error)
	}

	return wish.ID, nil
}

func (r *mutationResolver) AddWantToFulfill(ctx context.Context, input model.WishID) (*model.Wish, error) {
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

	d := r.DB.Select("id, name, user_id, description, link, image").First(&wish, input.ID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser == wish.UserID || !dbmodel.AreFriends(wish.UserID, authedUser) {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	asso := r.DB.Model(&dbmodel.Wish{ID: input.ID}).Where("user_id = ?", authedUser).Association("WantToFulfill")

	if asso.Count() != 0 {
		return nil, dbmodel.ErrUserExists
	}

	err = r.DB.Model(&dbmodel.Wish{ID: input.ID}).Association("WantToFulfill").Append(&dbmodel.User{ID: authedUser}).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not add to WantToFulfill", err)
	}

	return &model.Wish{
		ID:          wish.ID,
		User:        wish.UserID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
	}, nil
}

func (r *mutationResolver) AddClaimer(ctx context.Context, input model.WishID) (int, error) {
	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return 0, err
	}

	err = lib.Validator.Struct(&input)
	if err != nil {
		return 0, lib.ErrValidationFailed
	}

	asso := r.DB.Model(&dbmodel.Wish{ID: input.ID}).Where("user_id = ?", authedUser).Association("WantToFulfill")
	if asso.Error != nil && !gorm.IsRecordNotFoundError(asso.Error) {
		lib.LogError(lib.LPanic, "Could not read wish's WantToFulfill", asso.Error)
	}

	if asso.Count() != 1 {
		return 0, dbmodel.ErrUserNotFound
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		asso := tx.Model(&dbmodel.Wish{ID: input.ID}).Association("Claimers").Append(&dbmodel.User{ID: authedUser})
		if asso.Error != nil {
			return asso.Error
		}

		asso = tx.Model(&dbmodel.Wish{ID: input.ID}).Association("WantToFulfill").Delete(&dbmodel.User{ID: authedUser})
		if asso.Error != nil {
			return asso.Error
		}

		return nil
	})
	if err != nil {
		lib.LogError(lib.LPanic, "Could not add to Claimers", err)
	}

	return input.ID, nil
}

func (r *mutationResolver) AcceptClaimer(ctx context.Context, wishID model.WishID, claimer model.UserID) (*model.Wish, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&wishID)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	err = lib.Validator.Struct(&claimer)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, name, user_id, description, link, image").First(&wish, wishID.ID)
	if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.UserID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	count := r.DB.Model(&wish).Where("user_id = ?", claimer.ID).Association("Claimers").Count()
	if count != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		asso := tx.Model(&wish).Association("Fulfillers").Append(&dbmodel.User{ID: claimer.ID})
		if asso.Error != nil {
			return asso.Error
		}

		asso = tx.Model(&wish).Association("Claimers").Delete(&dbmodel.User{ID: claimer.ID})
		if asso.Error != nil {
			return asso.Error
		}

		return nil
	})
	if err != nil {
		lib.LogError(lib.LPanic, "Could not accept fulfillment claim", err)
	}

	return &model.Wish{
		ID:          wish.ID,
		User:        wish.UserID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
	}, nil
}

func (r *mutationResolver) RejectClaimer(ctx context.Context, wishID model.WishID, claimer model.UserID) (*model.Wish, error) {
	var wish dbmodel.Wish

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	err = lib.Validator.Struct(&wishID)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	err = lib.Validator.Struct(&claimer)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, user_id").First(&wish, wishID.ID)
	if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.UserID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	count := r.DB.Model(&wish).Where("user_id = ?", claimer.ID).Association("Claimers").Count()
	if count != 1 {
		return nil, dbmodel.ErrUserNotFound
	}

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		asso := tx.Model(&wish).Association("WantToFulfill").Append(&dbmodel.User{ID: claimer.ID})
		if asso.Error != nil {
			return asso.Error
		}

		asso = tx.Model(&wish).Association("Claimers").Delete(&dbmodel.User{ID: claimer.ID})
		if asso.Error != nil {
			return asso.Error
		}

		return nil
	})
	if err != nil {
		lib.LogError(lib.LPanic, "Could not reject fulfillment claim", err)
	}

	return &model.Wish{
		ID:          wish.ID,
		User:        wish.UserID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user dbmodel.User

	err := lib.Validator.Var(id, "username")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

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

func (r *queryResolver) Wish(ctx context.Context, input model.WishID) (*model.Wish, error) {
	var wish dbmodel.Wish

	err := lib.Validator.Struct(&input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Omit("fulfilled_by, created_at, updated_at").First(&wish, input.ID)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}
	return &model.Wish{
		ID:          wish.ID,
		User:        wish.UserID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
		Claimers:    wish.ID,
		Fulfillers:  wish.ID,
	}, nil
}

func (r *userResolver) Friends(ctx context.Context, obj *model.User, input *model.Page) ([]*model.User, error) {
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

	err = lib.Validator.Struct(input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(input.Page * 10) - 10).Limit(10).Association("Friends").Find(&friends)
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

func (r *userResolver) FriendRequests(ctx context.Context, obj *model.User, input *model.Page) ([]*model.User, error) {
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

	err = lib.Validator.Struct(input)
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Model(&dbmodel.User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(input.Page * 10) - 10).Limit(10).Association("FriendRequests").Find(&reqs)
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

func (r *wishResolver) User(ctx context.Context, obj *model.Wish) (*model.User, error) {
	var user dbmodel.User

	err := lib.Validator.Var(obj.User, "username")
	if err != nil {
		return nil, lib.ErrValidationFailed
	}

	d := r.DB.Select("id, first_name, last_name").Where("id = ?", obj.User).First(&user)
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

func (r *wishResolver) Claimers(ctx context.Context, obj *model.Wish) ([]*model.User, error) {
	var wish dbmodel.Wish
	var claimers []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	d := r.DB.Select("id, user_id").First(&wish, obj.Claimers)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.UserID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	r.DB.Model(&dbmodel.Wish{ID: obj.Claimers}).Association("Claimers").Find(&claimers)

	for _, c := range claimers {
		res = append(res, &model.User{
			ID:             c.ID,
			FirstName:      c.FirstName,
			LastName:       c.LastName,
			Friends:        c.ID,
			FriendRequests: c.ID,
		})
	}

	return res, nil
}

func (r *wishResolver) Fulfillers(ctx context.Context, obj *model.Wish) ([]*model.User, error) {
	var wish dbmodel.Wish
	var fulfillers []dbmodel.User
	var res []*model.User

	c := lib.GinCtxFromCtx(ctx)
	authedUser, err := dbmodel.Authenticate(c)
	if err != nil {
		return nil, err
	}

	d := r.DB.Select("id, user_id").First(&wish, obj.Fulfillers)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", d.Error)
	} else if d.RecordNotFound() {
		return nil, dbmodel.ErrWishNotFound
	}

	if authedUser != wish.UserID {
		return nil, dbmodel.ErrUserNotAuthorized
	}

	r.DB.Model(&dbmodel.Wish{ID: obj.Fulfillers}).Association("Fulfillers").Find(&fulfillers)

	for _, f := range fulfillers {
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

// Wish returns generated.WishResolver implementation.
func (r *Resolver) Wish() generated.WishResolver { return &wishResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type wishResolver struct{ *Resolver }
