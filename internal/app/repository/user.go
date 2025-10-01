package repository

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Repository) GetUserByID(id uuid.UUID) (ds.User, error) {
	user := ds.User{}
	sub := r.db.Where("uuid = ?", id).Find(&user)
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
	user.UUID = uuid.New()

	logrus.Info(user.Login, user.Password, user.IsModerator)

	sub := r.db.Create(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	return user, nil
}

func (r *Repository) ChangeProfile(login string, userJSON apitypes.UserJSON) (ds.User, error) {
	user := apitypes.UserFromJSON(userJSON)
	currUser, err := r.GetUserByLogin(login)
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

func (r *Repository) SignIn(userJSON apitypes.UserJSON) (string, error) {
	user, err := r.GetUserByLogin(userJSON.Login)
	if err != nil {
		return "", err
	}

	if user.Password != userJSON.Password {
		return "", errors.New("invalid password")
	}

	token, err := GenerateToken(user.UUID, user.IsModerator)
	if err != nil {
		return "", err
	}
	
	err = r.SaveJWTToken(user.UUID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *Repository) SignOut(id uuid.UUID) error {
	err := r.DeleteJWTToken(id)
	if err != nil {
		return err
	}
	return nil
}

func GenerateToken(id uuid.UUID, isModerator bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = id.String()
	claims["is_moderator"] = isModerator
	claims["exp"] = time.Hour * 1

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *Repository) SaveJWTToken(uuid uuid.UUID, token string) error {
	expiration := time.Hour * 1

	err := r.rd.Set(uuid.String(), token, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteJWTToken(uuid uuid.UUID) error {
	err := r.rd.Del(uuid.String()).Err()
	if err != nil {
		return err
	}
	return nil
}







func (r *Repository) FillWithUsers() {
	users := []ds.User{
		{
			Login: "admin",
			Password: "admin",
			IsModerator: true,
		},
		{
			Login: "user",
			Password: "user",
			IsModerator: false,
		},
	}
	for _, user := range users {
		r.CreateUser(apitypes.UserToJSON(user))
	}
}
