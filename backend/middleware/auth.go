package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"hy-board-backend/config"
	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	cfg := config.GlobalConfig.Auth
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.JWTExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenStr string) (*JWTClaims, error) {
	cfg := config.GlobalConfig.Auth
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// UniProxyAuth restricts access to node endpoints using the configured Token header
func UniProxyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		expectedToken := config.GlobalConfig.Auth.UniProxyToken

		if token == "" || token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized node token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserAuth verifies JWT tokens and injects the User model into context
func UserAuth(requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User account disabled"})
			c.Abort()
			return
		}

		if requireAdmin && !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privilege required"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}

// GraphQLAuthContext parses the Authorization header if present, injecting user into context without aborting
func GraphQLAuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err == nil && user.Status == 1 {
			c.Set("user", &user)
		}
		c.Next()
	}
}
