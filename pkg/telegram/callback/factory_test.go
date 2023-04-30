package callback

import (
	"fmt"
	"strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestData_New(t *testing.T) {
	d := New("test", "name", "code")
	type args struct {
		kwParts M
		parts   []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal keyword",
			args: args{
				kwParts: M{
					"name": "Foo",
					"code": "secret",
				},
			},
			want:    "test|Foo|secret",
			wantErr: assertNoError,
		},
		{
			name: "normal sequence",
			args: args{
				parts: []string{"John", "Do"},
			},
			want:    "test|John|Do",
			wantErr: assertNoError,
		},
		{
			name: "skip key",
			args: args{
				parts: []string{"A"},
			},
			wantErr: assertErrorIs(ErrValueNotPassed("code")),
		},
		{
			name: "too long",
			args: args{parts: []string{
				strings.Repeat("n", 32),
				strings.Repeat("b", 32),
			}},
			wantErr: assertErrorIs(ErrCallbackTooLong),
		},
		{
			name: "separator in value",
			args: args{
				parts: []string{"A|", ""},
			},
			wantErr: assertErrorIs(ErrSeparatorInValue),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.New(tt.args.kwParts, tt.args.parts...)
			if !tt.wantErr(t, err, fmt.Sprintf("New(%v, %v)", tt.args.kwParts, tt.args.parts)) {
				return
			}
			assert.Equalf(t, tt.want, got, "New(%v, %v)", tt.args.kwParts, tt.args.parts)
		})
	}

}

func assertNoError(t assert.TestingT, err error, i ...interface{}) bool {
	return assert.NoError(t, err, i...)
}

func assertErrorIs(target error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, target, i...)
	}
}

func TestData_Parse(t *testing.T) {
	d := New("test", "name", "code")
	tests := []struct {
		name    string
		arg     string
		want    M
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			arg:  "test|John|Do",
			want: M{
				"@":    "test",
				"name": "John",
				"code": "Do",
			},
			wantErr: assertNoError,
		},
		{

			name:    "invalid payload",
			arg:     "test",
			wantErr: assertErrorIs(ErrInvalidPayload),
		},
		{
			name:    "strange prefix",
			arg:     "not_test|A|B",
			wantErr: assertErrorIs(ErrInvalidPrefix),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Parse(tt.arg)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.arg)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.arg)
		})
	}
}
