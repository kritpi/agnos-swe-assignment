package entity

// toNullString maps an empty string to nil so it is stored as NULL. Storing an
// empty string instead would collide on the UNIQUE national_id / passport_id
// columns.
func toNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// fromNullString maps a NULL column back to an empty string.
func fromNullString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
