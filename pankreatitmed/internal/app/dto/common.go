package dto

type List[T any] struct {
	Items []T `json:"items"`
}
