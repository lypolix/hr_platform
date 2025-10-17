package port

type Paginated[T any] struct {
	Data       []T
	TotalPages int
}
