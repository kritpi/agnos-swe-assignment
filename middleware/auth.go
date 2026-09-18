package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

const staffClaimsKey = "staff_claims"

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RequireStaffLogin(JWTsecret string) gin.HandlerFunc {
	keyFunc := func(*jwt.Token) (any, error) { return []byte(JWTsecret), nil }

	return func(c *gin.Context) {
		scheme, token, ok := strings.Cut(c.GetHeader("Authorization"), " ")
		token = strings.TrimSpace(token)
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			unauthorized(c, "missing or malformed bearer token")
			return
		}

		var claims domain.StaffClaims
		_, err := jwt.ParseWithClaims(token, &claims, keyFunc,
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithExpirationRequired(),
		)
		if errors.Is(err, jwt.ErrTokenExpired) {
			unauthorized(c, "token expired")
			return
		}
		if err != nil {
			unauthorized(c, "invalid token")
			return
		}

		c.Set(staffClaimsKey, &claims)
		c.Next()
	}
}

func StaffClaims(c *gin.Context) (*domain.StaffClaims, bool) {
	v, ok := c.Get(staffClaimsKey)
	if !ok {
		return nil, false
	}
	claims, ok := v.(*domain.StaffClaims)
	return claims, ok
}

func unauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Code: http.StatusUnauthorized, Message: msg})
}
