package scheduleapi

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

var reGroupPattern = regexp.MustCompile(`^(\d\d)([ЁёА-я]+)-(\d+)$`)

type Group struct {
	Year   int    `json:"year"`
	Number int    `json:"number"`
	Spec   string `json:"spec"`
}

func (g *Group) UnmarshalText(text []byte) error {
	group, err := ParseGroup(unsafeBytesToString(text))
	if err != nil {
		return err
	}
	*g = group
	return nil
}

func (g Group) Value() (driver.Value, error) {
	return g.String(), nil
}

var ErrInvalidGroup = errors.New("invalid group format")

func (g Group) String() string {
	year := strconv.Itoa(g.Year)
	num := strconv.Itoa(g.Number)

	return year + strings.ToUpper(g.Spec) + "-" + num
}

func ParseGroup(strGroup string) (g Group, err error) {
	match := reGroupPattern.FindStringSubmatch(strGroup)
	if match == nil {
		return g, ErrInvalidGroup
	}

	yearInt, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	numInt, err := strconv.Atoi(match[3])
	if err != nil {
		return
	}

	g.Year = yearInt
	g.Number = numInt
	g.Spec = match[2]
	return

}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
