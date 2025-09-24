package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"errors"
)

func (r *Repository) GetUserByID(id int) (ds.User, error) {
	user := ds.User{}
	if id < 0 {
		return ds.User{}, errors.New("invalid id, it must be >= 0")
	}
	sub := r.db.Where("id = ?", id).Find(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.User{}, ErrorNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByLogin(login string) (ds.User, error) {
	user := ds.User{}
	sub := r.db.Where("login = ?", login).Find(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.User{}, ErrorNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) CreateUser(userJSON apitypes.UserJSON) (ds.User, error) {
	user := apitypes.UserFromJSON(userJSON)
	if user.Login == "" {
		return ds.User{}, errors.New("login is empty")
	}
	if user.Password == "" {
		return ds.User{}, errors.New("password is empty")
	}
	if _, err := r.GetUserByLogin(user.Login); err == nil {
		return ds.User{}, errors.New("user already exists")
	}
	currUser, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.User{}, err
	}
	if user.IsModerator && !currUser.IsModerator {
		return ds.User{}, errors.New("you are not a moderator")
	}
	sub := r.db.Create(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	return user, nil
}

func (r *Repository) ChangeProfile(id int, userJSON apitypes.UserJSON) (ds.User, error) {
	user := apitypes.UserFromJSON(userJSON)
	currUser, err := r.GetUserByID(id)
	if err != nil {
		return ds.User{}, err
	}
	if user.IsModerator && !currUser.IsModerator {
		user.IsModerator = false
	}
	err = r.db.Model(&currUser).Updates(user).Error
	if err != nil {
		return ds.User{}, err
	}
	return currUser, nil
}
