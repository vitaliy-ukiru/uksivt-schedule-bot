package group

import (
	"strconv"
	"strings"
)

type Group struct {
	Year   int
	Number int
	Spec   string
}

func (g Group) String() string {
	year := strconv.Itoa(g.Year)
	num := strconv.Itoa(g.Number)

	return year + strings.ToUpper(g.Spec) + "-" + num
}
