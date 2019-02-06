package contextionary

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/creativesoftwarefdn/weaviate/database/schema/kind"
	"github.com/fatih/camelcase"
)

// SearchResult is a single search result. See wrapping Search Results for the Type
type SearchResult struct {
	Name      string
	Kind      kind.Kind
	Certainty float32
}

// SearchResults is grouping of SearchResults for a SchemaSearch
type SearchResults struct {
	Type    SearchType
	Results []SearchResult
}

// Len of the result set
func (r SearchResults) Len() int {
	return len(r.Results)
}

// SchemaSearch can be used to search for related classes and properties, see
// documentation of SearchParams for more details on how to use it and
// documentation on SearchResults for more details on how to use the return
// value
func (mi *MemoryIndex) SchemaSearch(p SearchParams) (SearchResults, error) {
	result := SearchResults{}
	if err := p.Validate(); err != nil {
		return result, fmt.Errorf("invalid search params: %s", err)
	}

	centroid, err := mi.centroidFromNameAndKeywords(p)
	if err != nil {
		return result, fmt.Errorf("could not build centroid from name and keywords: %s", err)
	}

	rawResults, err := mi.knnSearch(*centroid)
	if err != nil {
		return result, fmt.Errorf("could not perform knn search: %s", err)
	}

	if p.SearchType == SearchTypeClass {
		return mi.handleClassSearch(p, rawResults)
	}

	// since we have passed validation we know that anything that's not a class
	// search must be a property search
	return mi.handlePropertySearch(p, rawResults)
}

func (mi *MemoryIndex) centroidFromNameAndKeywords(p SearchParams) (*Vector, error) {
	nameVector, err := mi.camelCaseWordToVector(p.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid name in search: %s", err)
	}

	if len(p.Keywords) == 0 {
		return nameVector, nil
	}

	vectors := make([]Vector, len(p.Keywords)+1, len(p.Keywords)+1)
	weights := make([]float32, len(p.Keywords)+1, len(p.Keywords)+1)
	// set last vector to className which always has weight=1
	vectors[len(vectors)-1] = *nameVector
	weights[len(vectors)-1] = 1

	for i, keyword := range p.Keywords {
		kwVector, err := mi.wordToVector(keyword.Keyword)
		if err != nil {
			return nil, fmt.Errorf("invalid keyword in search: %s", err)
		}
		vectors[i] = *kwVector
		weights[i] = keyword.Weight
	}

	return ComputeWeightedCentroid(vectors, weights)
}

func (mi *MemoryIndex) camelCaseWordToVector(w string) (*Vector, error) {
	parts := camelcase.Split(w)
	if len(parts) == 1 {
		// no camelcasing, no need to build a centroid
		return mi.wordToVector(w)
	}

	vectors := make([]Vector, len(parts), len(parts))
	weights := make([]float32, len(parts), len(parts))
	for i, part := range parts {
		v, err := mi.wordToVector(part)
		if err != nil {
			return nil, fmt.Errorf("invalid camelCased compound word: %s", err)
		}

		vectors[i] = *v
		weights[i] = 1 // on camel-casing all parts are weighted equally
	}

	return ComputeWeightedCentroid(vectors, weights)
}

func (mi *MemoryIndex) wordToVector(w string) (*Vector, error) {
	w = strings.ToLower(w)
	itemIndex := mi.WordToItemIndex(w)
	if ok := itemIndex.IsPresent(); !ok {
		return nil, fmt.Errorf(
			"the word '%s' is not present in the contextionary and therefore not a valid search term", w)
	}

	vector, err := mi.GetVectorForItemIndex(itemIndex)
	if err != nil {
		return nil, fmt.Errorf("could not get vector for word '%s' with itemIndex '%d': %s",
			w, itemIndex, err)
	}

	return vector, nil
}

func (mi *MemoryIndex) handleClassSearch(p SearchParams, search rawResults) (SearchResults, error) {
	return SearchResults{
		Type:    p.SearchType,
		Results: search.extractClassNames(p),
	}, nil
}

func (mi *MemoryIndex) handlePropertySearch(p SearchParams, search rawResults) (SearchResults, error) {
	return SearchResults{
		Type:    p.SearchType,
		Results: search.extractPropertyNames(p),
	}, nil
}

func (mi *MemoryIndex) knnSearch(vector Vector) (rawResults, error) {
	list, distances, err := mi.GetNnsByVector(vector, 10000, 3)
	if err != nil {
		return nil, fmt.Errorf("could not get nearest neighbors for vector '%v': %s", vector, err)
	}

	results := make(rawResults, len(list), len(list))
	for i := range list {
		word, err := mi.ItemIndexToWord(list[i])
		if err != nil {
			return results, fmt.Errorf("got a result from kNN search, but don't have a word for this index: %s", err)
		}

		results[i] = rawResult{
			name:     word,
			distance: distances[i],
		}
	}

	return results, nil
}

// rawResult is a helper struct to contain the results of the kNN-search. It
// does not yet contain the desired output. This means the names can be both
// classes/properties and arbitrary words. Furthermore the certainty has not
// yet been normalized , so it is merely the raw kNN distance
type rawResult struct {
	name     string
	distance float32
}

type rawResults []rawResult

func (r rawResults) extractClassNames(p SearchParams) []SearchResult {
	var results []SearchResult
	regex := regexp.MustCompile(fmt.Sprintf("^\\$%s\\[([A-Za-z]+)\\]$", p.Kind.AllCapsName()))

	for _, rawRes := range r {
		if regex.MatchString(rawRes.name) {
			certainty := distanceToCertainty(rawRes.distance)
			if certainty < p.Certainty {
				continue
			}

			results = append(results, SearchResult{
				Name:      regex.FindStringSubmatch(rawRes.name)[1], //safe because we ran .MatchString before
				Certainty: certainty,
				Kind:      p.Kind,
			})
		}
	}

	return results
}

func (r rawResults) extractPropertyNames(p SearchParams) []SearchResult {
	var results []SearchResult
	regex := regexp.MustCompile("^\\$[A-Za-z]+\\[([A-Za-z]+)\\]$")

	propsMap := map[string][]SearchResult{}

	for _, rawRes := range r {
		if regex.MatchString(rawRes.name) {
			name := regex.FindStringSubmatch(rawRes.name)[1] //safe because we ran .MatchString before
			certainty := distanceToCertainty(rawRes.distance)
			if certainty < p.Certainty {
				continue
			}

			res := SearchResult{
				Name:      name,
				Certainty: certainty,
			}
			if _, ok := propsMap[name]; !ok {
				propsMap[name] = []SearchResult{res}
			} else {
				propsMap[name] = append(propsMap[name], res)
			}
		}
	}

	// now calculate mean of duplicate results
	for _, resultsPerName := range propsMap {
		results = append(results, SearchResult{
			Name:      resultsPerName[0].Name,
			Certainty: meanCertainty(resultsPerName),
		})
	}

	return results
}

func meanCertainty(rs []SearchResult) float32 {
	var compound float32
	for _, r := range rs {
		compound += r.Certainty
	}

	return compound / float32(len(rs))
}

func distanceToCertainty(d float32) float32 {
	return 1 - d/12
}