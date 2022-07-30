package utils

import "github.com/letscrum/letscrum/pkg/settings"

// Setup Initialize the utils
func Setup() {
	jwtSecret = []byte(settings.AppSetting.JwtSecret)
}
