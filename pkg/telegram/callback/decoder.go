package callback

import (
	"errors"
	"fmt"
	"reflect"
)

const StructTag = "cb"

var ErrInvalidType = errors.New("dest must be non-nil pointer")

func Decode(src M, dest any) error {
	rv := reflect.ValueOf(dest)
	rt := rv.Type().Elem()
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrInvalidType
	}
	rv = rv.Elem()

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		if !ft.IsExported() {
			continue
		}

		if !fv.IsValid() || !fv.CanSet() {
			return ErrInvalidField(ft.Name)
		}

		if fv.Kind() != reflect.String {
			return ErrInvalidFieldType(ft.Name)
		}

		tag, ok := ft.Tag.Lookup(StructTag)
		if tag == "-" {
			continue
		}

		if !ok {
			tag = ft.Name
		}

		v, ok := src[tag]
		if !ok {
			return ErrNotHasValue(ft.Name)
		}

		fv.SetString(v)
	}

	return nil
}

type (
	ErrInvalidField     string
	ErrInvalidFieldType string
	ErrNotHasValue      string
)

func (e ErrInvalidField) Error() string {
	return fmt.Sprintf("invalid field: %q", string(e))
}
func (e ErrInvalidFieldType) Error() string {
	return fmt.Sprintf("field: %q not string", string(e))
}

func (e ErrNotHasValue) Error() string {
	return fmt.Sprintf("no value for field: %q", string(e))
}
