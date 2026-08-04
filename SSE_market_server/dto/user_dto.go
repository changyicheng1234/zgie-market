// token结构体
package dto

import "loginTest/model"

type UserDto struct {
	UserID    int    `json:"userID"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarURL"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Identity  string `json:"identity"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		UserID:    user.UserID,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Email:     user.Email,
		Identity:  user.Identity,
	}
}
