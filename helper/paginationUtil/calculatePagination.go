package paginationUtil


func CalculatePagination(first *int, after *string, ) (int, int, error) {
	limit := Limit(first)
    offset, err := Offset(after)
    if err != nil {
        return 0, 0, err
    }
    return offset, limit, nil
}