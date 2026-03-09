package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, constant.MsgAuthHeaderRequired, nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, constant.MsgInvalidAuthFormat, nil)
			c.Abort()
			return
		}

		token, err := security.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, constant.MsgInvalidToken, nil)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
			c.Set("jwt_token", token)
		}

		c.Next()
	}
}

func JWTAuthWithRefresh(tokenStore port.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, constant.MsgAuthHeaderRequired, nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, constant.MsgInvalidAuthFormat, nil)
			c.Abort()
			return
		}

		token, err := security.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, constant.MsgInvalidToken, nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, constant.MsgInvalidTokenClaim, nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		if security.ShouldRefresh(token) && tokenStore != nil {
			userID := security.GetUserID(c)
			role := security.GetRole(c)
			if userID > 0 && role != "" {
				newToken, refreshErr := tokenStore.RefreshToken(context.Background(), userID, role)
				if refreshErr == nil {
					c.Header("X-New-Token", newToken)
				}
			}
		}

		c.Next()
	}
}

func OptionalJWTAuth() gin.HandlerFunc {
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

		token, err := security.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}
