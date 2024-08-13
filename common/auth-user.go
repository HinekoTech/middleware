package common

import (
	"errors"
	"strings"
	"sync"

	"github.com/FourWD/middleware/orm"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	mu        sync.RWMutex
	blacklist []string
)

type UserAuthorization struct {
	IsSuccess bool
	Code      string
	Message   string
}

func CheckUserAuthorization(c *fiber.Ctx, db *gorm.DB, excludePath ...[]string) UserAuthorization {
	path := getLastPathComponent(c.Path())
	defaultExcludePath := []string{"login", "logout", "register", "wake-up", "warmup"}
	Print("CheckUserAuthorization", path)

	if StringExistsInList(path, defaultExcludePath) {
		return UserAuthorization{IsSuccess: true, Code: "200", Message: "ok"}
	}

	if len(excludePath) > 1 {
		if StringExistsInList(path, defaultExcludePath) || StringExistsInList(path, excludePath[0]) {
			return UserAuthorization{IsSuccess: true, Code: "200", Message: "ok"}
		}
	}

	bearerToken := c.Get("Authorization")
	token := strings.Replace(bearerToken, "Bearer ", "", 1)
	if token == "" {
		PrintError("LogUserLogin invalid request", token)
		return UserAuthorization{IsSuccess: false, Code: "401", Message: "invalid request"}
	}

	var logUserLogin orm.LogUserLogin
	if err := db.Where("token = ?", token).Order("created_at DESC").First(&logUserLogin).Error; err != nil {
		PrintError("LogUserLogin not found", token)
		return UserAuthorization{IsSuccess: false, Code: "401", Message: "log_login not found"}
	}

	if logUserLogin.ID == "" {
		PrintError("LogUserLogin unauthorized", token)
		return UserAuthorization{IsSuccess: false, Code: "401", Message: "unauthorized"}
	}

	return UserAuthorization{IsSuccess: true, Code: "200", Message: "ok"}
}

func getLastPathComponent(path string) string {
	components := strings.Split(path, "/")
	lastComponent := components[len(components)-1]
	if lastComponent == "favicon.ico" {
		return ""
	}
	return lastComponent
}

func Logout(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	if token == "" {
		return errors.New("no token")
	}

	// return Database.Model(&orm.JwtBlacklist{}).Create(&orm.JwtBlacklist{
	// 	ID:    uuid.NewString(),
	// 	Md5:   MD5(token),
	// 	Token: token,
	// }).Error
	return addJwtBlacklist(token)
}

func addJwtBlacklist(token string) error {
	mu.Lock()
	defer mu.Unlock()

	// Check if the blacklist has reached its max size
	maxBlacklistSize := 100
	if len(blacklist) >= maxBlacklistSize {
		// Remove the oldest token (first in the slice)
		blacklist = blacklist[1:]
	}

	// Add the new token to the end of the slice
	blacklist = append(blacklist, token)
	return nil
}
