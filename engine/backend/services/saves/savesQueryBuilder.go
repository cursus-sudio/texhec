package saves

type listSavesQueryBuilder struct {
	query ListSavesQuery
}

func newQueryBuilder() ListSavesQueryBuilder {
	return listSavesQueryBuilder{
		query: ListSavesQuery{
			OrderedBy:    OrderedByCreated,
			SortOrder:    AscOrder,
			SavesPerPage: 10,
			CurrentPage:  0,
		},
	}
}

func (builder listSavesQueryBuilder) OrderByCreated() ListSavesQueryBuilder {
	builder.query.OrderedBy = OrderedByCreated
	return builder
}

func (builder listSavesQueryBuilder) OrderByLastModified() ListSavesQueryBuilder {
	builder.query.OrderedBy = OrderedByLastModified
	return builder
}

func (builder listSavesQueryBuilder) AscOrder() ListSavesQueryBuilder {
	builder.query.SortOrder = AscOrder
	return builder
}
func (builder listSavesQueryBuilder) DescOrder() ListSavesQueryBuilder {
	builder.query.SortOrder = DescOrder
	return builder
}

func (builder listSavesQueryBuilder) SavesPerPage(saves uint) ListSavesQueryBuilder {
	if saves < 1 {
		saves = 1
	}
	builder.query.SavesPerPage = saves
	return builder
}
func (builder listSavesQueryBuilder) CurrentPage(current uint) ListSavesQueryBuilder {
	builder.query.CurrentPage = current
	return builder
}

func (builder listSavesQueryBuilder) Build() ListSavesQuery {
	return builder.query
}
