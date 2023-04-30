package scheduleapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGroup(t *testing.T) {
	type args struct {
		strGroup string
	}
	tests := []struct {
		name    string
		args    args
		wantG   Group
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:  "Correct",
			args:  args{strGroup: "20П-1"},
			wantG: Group{Year: 20, Spec: "П", Number: 1},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
		{
			name: "Invalid",
			args: args{strGroup: "1foo-bar"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrInvalidGroup, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotG, err := ParseGroup(tt.args.strGroup)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseGroup(%v)", tt.args.strGroup)) {
				return
			}
			assert.Equalf(t, tt.wantG, gotG, "ParseGroup(%v)", tt.args.strGroup)
		})
	}
}
