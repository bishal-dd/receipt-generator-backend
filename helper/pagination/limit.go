package pagination


func Limit(limit *int) int {
	if limit == nil {
		return 10 
	}
	return *limit
}