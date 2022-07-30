package util

import "github.com/letscrum/letscrum/pkg/setting"

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
