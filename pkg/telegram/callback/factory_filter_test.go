package callback

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_Filter(t *testing.T) {
	d := New("test", "name", "code")
	tests := []struct {
		name    string
		args    []M
		want    Filter
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty filter",
			want: Filter{
				data:   d,
				config: nil,
			},
			wantErr: assertNoError,
		},
		{
			name: "simple filter",
			args: []M{
				{
					"name": "John",
				},
			},
			want: Filter{
				data:   d,
				config: M{"name": "John"},
			},
			wantErr: assertNoError,
		},
		{
			name: "invalid key",
			args: []M{
				{
					"surname": "Doe",
				},
			},
			wantErr: assertErrorIs(ErrInvalidKey("surname")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Filter(tt.args...)
			if !tt.wantErr(t, err, fmt.Sprintf("Filter(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Filter(%v)", tt.args)
		})
	}
}

func TestFilter_check(t *testing.T) {
	d := New("test", "name", "code")

	type fields struct {
		data   *Data
		config M
	}
	tests := []struct {
		name    string
		fields  fields
		arg     M
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty filter",
			fields: fields{
				data:   d,
				config: nil,
			},
			wantErr: assertNoError,
		},
		{
			name: "correct filter",
			fields: fields{
				data: d,
				config: M{
					"name": "Foo",
				},
			},
			arg: M{
				"name": "Foo",
			},
			wantErr: assertNoError,
		},
		{
			name: "empty data", // TODO: new name
			fields: fields{
				data: d,
				config: M{
					"name": "A",
				},
			},
			arg:     M{},
			wantErr: assertErrorIs(ErrFilter), // because we specified name
		},
		{
			name: "unknown field",
			fields: fields{
				data: d,
				config: M{
					"name": "A",
				},
			},
			arg: M{
				"name": "A",
				"code": "B",
			},
			wantErr: assertNoError,
		},
		{
			name: "unknown empty field",
			fields: fields{
				data: d,
				config: M{
					"name":          "A",
					"unknown_field": "",
				},
			},
			arg: M{
				"name": "A",
			},
			wantErr: assertErrorIs(ErrFilter),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				data:   tt.fields.data,
				config: tt.fields.config,
			}
			tt.wantErr(t, f.check(tt.arg), fmt.Sprintf("check(%v)", tt.arg))
		})
	}
}
