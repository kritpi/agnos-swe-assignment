package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/kritpi/agnos-swe-assignment/internal/core/domain"
)

// VerifyStaff checks that staffID is still a member of hospitalID, so a token
// stops working once its staff member is removed or moved.
func (s *service) VerifyStaff(ctx context.Context, hospitalID, staffID string) error {
	_, err := s.repo.FindStaffByID(ctx, hospitalID, staffID)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrStaffNotMember
	}
	if err != nil {
		return fmt.Errorf("find staff: %w", err)
	}
	return nil
}
