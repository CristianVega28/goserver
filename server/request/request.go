package request

import (
	"net/http"
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
	QUERY_SEARCH_LTN          = "_ltn"
	QUERY_SEARCH_GTN          = "_gtn"
)

type (
	Filters struct {
		Value        []string
		Key          string
		Builder      *models.Builder
		IsQuery      bool
		IsPagination bool
		IsRegx       bool
	}

	RequestQueries struct {
		Cfg     *helpers.ConfigServerApi // only use middleware to verify db config is true and response
		Query   url.Values
		Request *http.Request
		Model   *models.Models[map[string]any]
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
		QUERY_SEARCH_GTN,
		QUERY_SEARCH_LTN,
	}

	var present_in_rules bool = false
	for _, v := range rules {
		if rq.Query.Has(v) {
			present_in_rules = true
		}
	}

	if !present_in_rules {
		return rq.Model.SelectAll()
	}

	for key, values := range rq.Request.URL.Query() {
		builder := rq.Model.Builder([]string{})

		filters := Filters{
			Value:   values,
			Key:     key,
			Builder: builder,
		}

		switch key {
		case QUERY_SEARCH_COLUMN:
			filters.FilterSearchColumn()
		}
	}

	return nil
}

// func (filt *Filters) New() Filters {
// 	return Filters{}
// }

func (fl *Filters) FilterSearchColumn() {

}
