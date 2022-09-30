package groups

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strconv"
)

type FileSystem struct {
	//        year       spec     num
	groups map[int]map[string][]int
}

func (f FileSystem) Years() []int {
	years := make([]int, 0, len(f.groups))
	for year := range f.groups {
		years = append(years, year)
	}
	sort.Ints(years)
	return years
}

func (f FileSystem) Specs(year int) []string {
	v, ok := f.groups[year]
	if !ok {
		return nil
	}
	specs := make([]string, 0, len(v))
	for spec := range v {
		specs = append(specs, spec)
	}
	return specs
}

func (f FileSystem) Numbers(year int, spec string) []int {
	specs, ok := f.groups[year]
	if !ok {
		return nil
	}
	numbers := append([]int(nil), specs[spec]...)
	sort.Ints(numbers)
	return numbers
}

func NewFS(groups map[int]map[string][]int) *FileSystem {
	return &FileSystem{groups: groups}
}

func NewFSFromReader(reader io.Reader) (*FileSystem, error) {
	m := make(map[string]map[string][]int)
	if err := json.NewDecoder(reader).Decode(&m); err != nil {
		return nil, err
	}
	g, err := toIntYears(m)
	if err != nil {
		return nil, err
	}
	return &FileSystem{groups: g}, nil
}

func NewFSFromFile(path string) (*FileSystem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewFSFromReader(file)
}

func toIntYears(groups map[string]map[string][]int) (map[int]map[string][]int, error) {
	result := make(map[int]map[string][]int)
	for y, spec := range groups {
		year, err := strconv.Atoi(y)
		if err != nil {
			return nil, err
		}
		result[year] = spec
	}
	return result, nil
}
