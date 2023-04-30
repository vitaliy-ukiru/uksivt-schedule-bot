package callback

import (
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Filter struct {
	data   *Data
	config M
}

type ErrInvalidKey string

func (v ErrInvalidKey) Error() string { return fmt.Sprintf("invalid key: %q", string(v)) }
func (v ErrInvalidKey) Key() string   { return string(v) }

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

func (d *Data) MustFilter(kw ...M) Filter {
	f, err := d.Filter(kw...)
	if err != nil {
		panic(err)
	}

	return f
}

var ErrFilter = errors.New("filter is not equals")

func (f Filter) Process(c tele.Context) (M, error) {
	data, err := f.data.ParseTele(c.Callback())
	if err != nil {
		return nil, err
	}

	return data, f.check(data)
}

func (f Filter) check(data M) error {
	if f.config == nil { // not specific filter
		return nil
	}

	for key, value := range f.config {
		if v, ok := data[key]; !ok || v != value {
			return ErrFilter
		}
	}

	return nil
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

func makeSet(arr []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range arr {
		m[s] = struct{}{}
	}

	return m
}
