package result

type Pager struct {
	Data   interface{}
	Total  int
	Limit  int
	Offset int
}
