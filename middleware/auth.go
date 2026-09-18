package middleware

import (
	"context"
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

// StaffVerifier checks that a token's staff member still belongs to its
// hospital.
type StaffVerifier interface {
	VerifyStaff(ctx context.Context, hospitalID, staffID string) error
}

func RequireStaffLogin(JWTsecret string, staff StaffVerifier) gin.HandlerFunc {
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

		// the token was valid at login; make sure the membership still is
		err = staff.VerifyStaff(c.Request.Context(), claims.HospitalID, claims.StaffID)
		if errors.Is(err, domain.ErrStaffNotMember) {
			unauthorized(c, domain.ErrStaffNotMember.Error())
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"})
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
