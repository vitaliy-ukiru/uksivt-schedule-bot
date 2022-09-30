package callback

import (
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func makeSet(arr []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range arr {
		m[s] = struct{}{}
	}

	return m
}

type Filter struct {
	data   *Data
	config M
}

type ErrInvalidKey string

func (d *Data) Filter(kw ...M) (Filter, error) {
	var m M

	if len(kw) > 0 {
		m = kw[0]
		set := makeSet(d.keys)
		for key := range m {
			if _, ok := set[key]; !ok {
				return Filter{}, ErrInvalidKey(key)
			}
		}
	}

	return Filter{
		data:   d,
		config: m,
	}, nil
}

func (d Data) MustFilter(kw ...M) Filter {
	f, err := d.Filter(kw...)
	if err != nil {
		panic(err)
	}

	return f
}

type HandlerFunc func(ctx tele.Context, data M) error

type Handleable interface {
	Handle(endpoint interface{}, h tele.HandlerFunc, m ...tele.MiddlewareFunc)
}

func (f Filter) Handle(group Handleable, h HandlerFunc, m ...tele.MiddlewareFunc) {
	group.Handle(f.data, func(c tele.Context) error {
		data, err := f.Process(c)
		if err != nil {
			if errors.Is(err, ErrFilter) {
				return nil
			}
			return err
		}

		return h(c, data)
	}, m...)
}

var ErrFilter = errors.New("filter is not equals")

func (f Filter) Process(c tele.Context) (M, error) {
	data, err := f.data.ParseTele(c.Callback())
	if err != nil {
		return nil, err
	}

	if f.config == nil { // not specific filter
		return data, nil
	}

	for key, value := range data {
		if f.config[key] != value {
			return nil, ErrFilter
		}
	}

	return data, nil

}

func (v ErrInvalidKey) Error() string { return fmt.Sprintf("invalid key: %q", v) }
func (v ErrInvalidKey) Key() string   { return string(v) }
