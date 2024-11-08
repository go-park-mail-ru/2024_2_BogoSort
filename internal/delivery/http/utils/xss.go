package utils

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeResponseAdvert(advert *dto.AdvertResponse, policy *bluemonday.Policy) {
	advert.Title = policy.Sanitize(advert.Title)
	advert.Description = policy.Sanitize(advert.Description)
}

func SanitizeRequestAdvert(advert *dto.AdvertRequest, policy *bluemonday.Policy) {
	advert.Title = policy.Sanitize(advert.Title)
	advert.Description = policy.Sanitize(advert.Description)
}

func SanitizeRequestSignup(credentials *dto.Signup, policy *bluemonday.Policy) {
	credentials.Email = policy.Sanitize(credentials.Email)
	credentials.Password = policy.Sanitize(credentials.Password)
}

func SanitizeRequestLogin(credentials *dto.Login, policy *bluemonday.Policy) {
	credentials.Email = policy.Sanitize(credentials.Email)
	credentials.Password = policy.Sanitize(credentials.Password)
}

func SanitizeRequestChangePassword(updatePassword *dto.UpdatePassword, policy *bluemonday.Policy) {
	updatePassword.OldPassword = policy.Sanitize(updatePassword.OldPassword)
	updatePassword.NewPassword = policy.Sanitize(updatePassword.NewPassword)
}

func SanitizeRequestUser(user *dto.User, policy *bluemonday.Policy) {
	user.Username = policy.Sanitize(user.Username)
	user.Phone = policy.Sanitize(user.Phone)
}

func SanitizeResponseUser(user *dto.User, policy *bluemonday.Policy) {
	user.Username = policy.Sanitize(user.Username)
	user.Phone = policy.Sanitize(user.Phone)
}

func SanitizeRequestUserUpdate(user *dto.UserUpdate, policy *bluemonday.Policy) {
	user.Email = policy.Sanitize(user.Email)
	user.Username = policy.Sanitize(user.Username)
	user.Phone = policy.Sanitize(user.Phone)
}
