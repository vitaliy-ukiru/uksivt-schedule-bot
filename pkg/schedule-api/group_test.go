package scheduleapi

import (
	"reflect"
	"testing"
)

func TestParseGroup(t *testing.T) {
	type args struct {
		strGroup string
	}
	tests := []struct {
		name    string
		args    args
		wantG   Group
		wantErr bool
	}{
		{
			name:  "Correct",
			args:  args{strGroup: "20П-1"},
			wantG: Group{Year: 20, Spec: "П", Number: 1},
		},
		{
			name:    "Invalid",
			args:    args{strGroup: "1foo-bar"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotG, err := ParseGroup(tt.args.strGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotG, tt.wantG) {
				t.Errorf("ParseGroup() gotG = %v, want %v", gotG, tt.wantG)
			}
		})
	}
}
