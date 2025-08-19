package stringUtil

func StrPtr(s string) *string {
	return &s
}

func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
