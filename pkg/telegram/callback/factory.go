package callback

import (
	"errors"
	"fmt"
	"strings"
)

type M map[string]string

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

func (d *Data) New(kwParts M, parts ...string) (string, error) {
	data, err := d.buildData(kwParts, parts)
	if err != nil {
		return "", err
	}

	result := strings.Join(data, d.Sep)
	if len(result) > MaxCallbackDataSize {
		return "", ErrCallbackTooLong
	}
	return d.Prefix + d.Sep + result, nil
}

var (
	ErrSeparatorInValue = errors.New("separator symbol in callback value")
	ErrCallbackTooLong  = errors.New("callback data too long")
)

type ErrValueNotPassed string

func (v ErrValueNotPassed) Error() string {
	return fmt.Sprintf("not passed value for key: %q", string(v))
}
func (v ErrValueNotPassed) Key() string { return string(v) }

func (d *Data) buildData(kwArgs M, args []string) ([]string, error) {
	data := make([]string, 0, len(kwArgs)+len(args))

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
	ErrInvalidPayload    = errors.New("invalid payload format")
)

func (d *Data) Parse(data string) (M, error) {
	parts := strings.Split(data, d.Sep)
	if len(parts) < 2 {
		return nil, ErrInvalidPayload
	}
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
