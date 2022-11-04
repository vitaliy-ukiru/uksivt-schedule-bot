package groups

import (
	"encoding/json"
	"io"
	"os"
	"sort"

	"github.com/pkg/errors"
)

type InMemoryService struct {
	//        year       spec     num
	groups map[int]map[string][]int
}

func (i InMemoryService) Years() []int {
	years := make([]int, 0, len(i.groups))
	for year := range i.groups {
		years = append(years, year)
	}
	sort.Ints(years)
	return years
}

func (i InMemoryService) Specs(year int) []string {
	v, ok := i.groups[year]
	if !ok {
		return nil
	}
	specs := make([]string, 0, len(v))
	for spec := range v {
		specs = append(specs, spec)
	}
	sort.Strings(specs)
	return specs
}

func (i InMemoryService) Numbers(year int, spec string) []int {
	specs, ok := i.groups[year]
	if !ok {
		return nil
	}
	numbers := append([]int(nil), specs[spec]...)
	sort.Ints(numbers)
	return numbers
}

func NewInMemoryService(groups map[int]map[string][]int) *InMemoryService {
	return &InMemoryService{groups: groups}
}

func NewInMemoryFromReader(reader io.Reader) (*InMemoryService, error) {
	m := make(map[int]map[string][]int)
	if err := json.NewDecoder(reader).Decode(&m); err != nil {
		return nil, errors.Wrap(err, "in-mem group service: cannot decode")
	}

	return &InMemoryService{groups: m}, nil
}

func NewInMemoryFromFile(path string) (*InMemoryService, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "in-mem group service: cannot open file")
	}
	defer file.Close()
	return NewInMemoryFromReader(file)
}
