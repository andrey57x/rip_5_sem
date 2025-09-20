package apitypes

import "Backend/internal/app/ds"

// UserJSON model
// @Description Model for user
// @Tags users
type UserJSON struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
}

func UserToJSON(user ds.User) UserJSON {
	return UserJSON{
		ID:          user.ID,
		Login:       user.Login,
		Password:    user.Password,
		IsModerator: user.IsModerator,
	}
}

func UserFromJSON(userJSON UserJSON) ds.User {
	return ds.User{
		Login:       userJSON.Login,
		Password:    userJSON.Password,
		IsModerator: userJSON.IsModerator,
	}
}
