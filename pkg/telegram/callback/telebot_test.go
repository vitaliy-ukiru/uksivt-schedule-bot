package callback

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestData_TeleBtn(t *testing.T) {
	d := New("test", "name", "code")
	type args struct {
		text    string
		kwParts M
		parts   []string
	}
	tests := []struct {
		name    string
		args    args
		wantBtn tele.Btn
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct",
			args: args{
				text: "Text",
				kwParts: M{
					"name": "Bar",
					"code": "187",
				},
			},
			wantBtn: tele.Btn{
				Unique: d.Prefix,
				Text:   "Text",
				Data:   "Bar|187",
			},
			wantErr: assertNoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBtn, err := d.TeleBtn(tt.args.text, tt.args.kwParts, tt.args.parts...)
			if !tt.wantErr(t, err, fmt.Sprintf("TeleBtn(%v, %v, %v)", tt.args.text, tt.args.kwParts, tt.args.parts)) {
				return
			}
			assert.Equalf(t, tt.wantBtn, gotBtn, "TeleBtn(%v, %v, %v)", tt.args.text, tt.args.kwParts, tt.args.parts)
		})
	}
}

func TestData_ParseTele(t *testing.T) {
	d := New("test", "name", "code")

	type args struct {
		callback *tele.Callback
	}
	tests := []struct {
		name    string
		args    args
		want    M
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct",
			args: args{
				callback: &tele.Callback{
					Data:   "Conor|iii",
					Unique: "test",
				},
			},
			want: M{
				"@":    "test",
				"name": "Conor",
				"code": "iii",
			},
			wantErr: assertNoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.ParseTele(tt.args.callback)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseTele(%v)", tt.args.callback)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseTele(%v)", tt.args.callback)
		})
	}
}
