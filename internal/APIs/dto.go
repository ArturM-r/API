package APIs

type Update struct {
	Completed *bool   `json:"completed"`
	Title     *string `json:"title"`
}
type QueryFilter struct {
	completed *bool
	title     *string
	limit     *int
	offset    *int
}
type QueryFilterFromHandler struct {
	Completed *bool   `json:"completed"`
	Title     *string `json:"title"`
	Limit     *int    `json:"limit"`
	Offset    *int    `json:"offset"`
}
