package callback

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Callback interface {
	Prefix() string
	Separator() string

	Keys() []string
}

var ErrNilCallback = errors.New("callback is nil")
var ErrInvalidCallbackType = errors.New("callback is not struct")

func Compile(cb Callback) (string, error) {
	rv := reflect.ValueOf(cb)
	rt := rv.Type()
	if rv.IsNil() || !rv.IsValid() {
		return "", ErrNilCallback
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", ErrInvalidCallbackType
	}

	data := make(map[string]string)
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		if !ft.IsExported() {
			continue
		}

		if !fv.IsValid() || !fv.CanSet() {
			return "", ErrInvalidField(ft.Name)
		}

		//if fv.Kind() != reflect.String {
		//	return "", ErrInvalidFieldType(ft.Name)
		//}

		tag, ok := ft.Tag.Lookup(StructTag)
		if tag == "-" {
			continue
		}

		if !ok {
			tag = ft.Name
		}

		switch obj := fv.Interface().(type) {
		case string:
			data[tag] = obj
		case int8, int16, int32, int64, int:
			i := fv.Int()
			data[tag] = strconv.FormatInt(i, 10)
		case uint8, uint16, uint32, uint64, uint:
			i := fv.Uint()
			data[tag] = strconv.FormatUint(i, 10)
		case float32, float64:
			i := fv.Float()
			data[tag] = strconv.FormatFloat(i, 'f', -1, 64)
		case bool:
			b := fv.Bool()
			data[tag] = strconv.FormatBool(b)
		}
	}

	if len(cb.Keys()) != len(data) {
		return "", ErrInvalidPartsCount
	}

	result := []string{cb.Prefix()}

	for _, key := range cb.Keys() {
		v, ok := getKey(data, key)
		if !ok {
			return "", ErrValueNotPassed(key)
		}

		if strings.Contains(v, cb.Separator()) {
			return "", ErrSeparatorInValue
		}

		result = append(result, v)
	}

	return strings.Join(result, cb.Separator()), nil
}

func Decompile(callback string, dest Callback) error {
	parts := strings.Split(callback, dest.Separator())
	prefix, parts := parts[0], parts[1:]
	if prefix != dest.Prefix() {
		return ErrInvalidPrefix
	}

	if len(parts) != len(dest.Keys()) {
		return ErrInvalidPartsCount
	}

	src := M{
		PrefixMapKey: prefix,
	}

	for i, key := range dest.Keys() {
		src[key] = parts[i]
	}

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

		switch ft.Type.Kind() {
		case reflect.String:
			fv.SetString(v)
		case reflect.Int:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			fv.SetInt(i)
		case reflect.Uint:
			i, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			fv.SetUint(i)
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			fv.SetBool(b)
		default:
			return ErrInvalidFieldType(ft.Name)
		}
	}

	return nil

}
