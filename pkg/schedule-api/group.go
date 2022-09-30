package scheduleapi

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var reGroupPattern = regexp.MustCompile(`^(\d\d)([ЁёА-я]+)-(\d+)$`)

type Group struct {
	Year   int    `json:"year"`
	Number int    `json:"number"`
	Spec   string `json:"spec"`
}

func (g Group) Value() (driver.Value, error) {
	return g.String(), nil
}

var ErrInvalidGroup = errors.New("invalid group format")

func (g *Group) UnmarshalJSON(bytes []byte) error {
	if bytes[0] == '"' && bytes[len(bytes)-1] == '"' { // its string
		var strMessage string

		if err := json.Unmarshal(bytes, &strMessage); err != nil {
			return err
		}
		group, err := ParseGroup(strMessage)
		if err != nil {
			return err
		}
		g.Year = group.Year
		g.Spec = group.Spec
		g.Number = group.Number
		return nil
	}

	var group struct {
		Year   int    `json:"year"`
		Number int    `json:"number"`
		Spec   string `json:"spec"`
	}

	if err := json.Unmarshal(bytes, &group); err != nil {
		return err
	}

	*g = group
	return nil
}

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
