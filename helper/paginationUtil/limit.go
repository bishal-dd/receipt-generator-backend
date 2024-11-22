package paginationUtil


func Limit(first *int) int {
    if first == nil || *first <= 0 {
        return 10 // Default limit
    }
    return *first
}
