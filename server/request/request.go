package request

import (
	"net/url"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/helpers"
)

const (
	QUERY_SORT                = "_sort"
	QUERY_PAGINATION_PAGE     = "_page"
	QUERY_PAGINATION_PER_PAGE = "_per_page"
	QUERY_SEARCH_COLUMN       = "_column"
	QUERY_SEARCH_WHERE        = "_where"
	QUERY_SEARCH_OR           = "_or"
	QUERY_SEARCH_AND          = "_and"
)

type (
	Filters struct {
		Value        string
		IsQuery      bool
		IsPagination bool
		IsRegx       bool
	}

	RequestQueries struct {
		Cfg   *helpers.ConfigServerApi // only use middleware to verify db config is true and response
		Query url.Values
		Model *models.Models[map[string]any]
	}
)

func (rq *RequestQueries) GetResponse() any {

	var rules = []string{
		QUERY_PAGINATION_PAGE,
		QUERY_PAGINATION_PER_PAGE,
		QUERY_SEARCH_AND,
		QUERY_SEARCH_COLUMN,
		QUERY_SEARCH_OR,
		QUERY_SEARCH_WHERE,
		QUERY_SORT,
	}

	return nil
}

func (filt *Filters) New() Filters {
	return Filters{}
}
