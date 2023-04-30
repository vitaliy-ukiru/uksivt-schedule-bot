package callback

import (
	"errors"
	"fmt"
	"strings"
)

type Data struct {
	Prefix string
	Sep    string

	keys []string
}

const TelebotSeparator = "|"

var Separator = TelebotSeparator

func New(prefix string, keys ...string) *Data {
	return &Data{
		Prefix: prefix,
		keys:   keys,
		Sep:    Separator,
	}
}

func NewWithSeparator(prefix string, sep string, keys ...string) *Data {
	return &Data{
		Prefix: prefix,
		keys:   keys,
		Sep:    sep,
	}
}

const MaxCallbackDataSize = 64

type M map[string]string

type ErrValueNotPassed string

var ErrSeparatorInValue = errors.New("separator symbol in callback value")
var ErrCallbackTooLong = errors.New("callback data too long")

func (d Data) New(kwParts M, parts ...string) (string, error) {
	data, err := d.buildData(kwParts, parts)
	if err != nil {
		return "", err
	}

	result := strings.Join(data, d.Sep)
	if len(result) > MaxCallbackDataSize {
		return "", ErrCallbackTooLong
	}
	return result, nil
}

func (d *Data) buildData(kwArgs M, args []string) ([]string, error) {
	var data []string

	args = append([]string(nil), args...) // copy

	for _, key := range d.keys {
		v, ok := getKey(kwArgs, key)
		if !ok {
			if len(args) == 0 {
				return nil, ErrValueNotPassed(key)
			}
			// pop first item
			v, args = args[0], args[1:]
		}

		if strings.Contains(v, d.Sep) {
			return nil, ErrSeparatorInValue
		}

		data = append(data, v)
	}
	return data, nil

}

var (
	ErrInvalidPrefix     = errors.New("invalid callback data prefix")
	ErrInvalidPartsCount = errors.New("invalid parts count")
)

func (d *Data) Parse(data string) (M, error) {
	parts := strings.Split(data, d.Sep)
	prefix, parts := parts[0], parts[1:]
	return d.union(prefix, parts)

}

const PrefixMapKey = "@"

func (d *Data) union(prefix string, parts []string) (M, error) {
	if prefix != d.Prefix {
		return nil, ErrInvalidPrefix
	}

	if len(parts) != len(d.keys) {
		return nil, ErrInvalidPartsCount
	}

	result := M{
		PrefixMapKey: prefix,
	}

	for i, key := range d.keys {
		result[key] = parts[i]
	}

	return result, nil
}

func getKey[T comparable, K any](m map[T]K, key T) (v K, ok bool) {
	if m != nil {
		v, ok = m[key]
	}

	return
}

func (v ErrValueNotPassed) Error() string {
	return fmt.Sprintf("not passed value for key: %q", string(v))
}
func (v ErrValueNotPassed) Key() string { return string(v) }
