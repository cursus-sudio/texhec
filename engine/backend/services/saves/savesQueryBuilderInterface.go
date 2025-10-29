package saves

import "errors"

const (
	OrderedByCreated int = iota
	OrderedByLastModified

	AscOrder int = iota
	DescOrder
)

var (
	ErrPageHasToHaveAtLeast5Saves  error = errors.New("save query: page has to have at least 5 saves")
	ErrPageHasToHaveAtMost100Saves error = errors.New("save query: page has to have at most 100 saves")
	ErrInvalidEnumValues           error = errors.New("save query: enums are not recognized")
)

type ListSavesQuery struct {
	OrderedBy    int
	SortOrder    int
	SavesPerPage uint
	CurrentPage  uint
}

// can return:
// - ErrPageHasToHaveAtLeast5Saves
func (listSavesQuery *ListSavesQuery) Valid() []error {
	errs := []error{}
	switch listSavesQuery.OrderedBy {
	case OrderedByCreated:
	case OrderedByLastModified:
	default:
		errs = append(errs, ErrInvalidEnumValues)
	}
	switch listSavesQuery.SortOrder {
	case AscOrder:
	case DescOrder:
	default:
		errs = append(errs, ErrInvalidEnumValues)
	}
	if listSavesQuery.SavesPerPage < 5 {
		errs = append(errs, ErrPageHasToHaveAtLeast5Saves)
	}
	if listSavesQuery.SavesPerPage > 100 {
		errs = append(errs, ErrPageHasToHaveAtMost100Saves)
	}
	return errs
}

type ListSavesQueryBuilder interface {
	// default
	OrderByCreated() ListSavesQueryBuilder
	OrderByLastModified() ListSavesQueryBuilder

	// default
	AscOrder() ListSavesQueryBuilder
	DescOrder() ListSavesQueryBuilder

	// if saves per page are bellow 1 then its threated as 1
	SavesPerPage(uint) ListSavesQueryBuilder // default 10
	CurrentPage(uint) ListSavesQueryBuilder  // default 0

	Build() ListSavesQuery
}
