package dancok

type SelectResult[T any] struct {
	Items           []T
	TotalItemsCount int64
}
