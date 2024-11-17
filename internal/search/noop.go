package search

import "context"

func NewNoopSearcher() Searcher {
	return noopSearcher{}
}

type noopSearcher struct{}

func (s noopSearcher) Close() {

}

func (s noopSearcher) Index(context.Context, IndexArgs) error {
	return nil
}

func (s noopSearcher) BatchIndex(context.Context, BatchIndexArgs) error {
	return nil
}

func (s noopSearcher) Search(context.Context, SearchArgs) (SearchResults, error) {
	return SearchResults{}, nil
}
