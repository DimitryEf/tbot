package golang

import (
	"reflect"
	"testing"
)

func TestConvertQueryToTopic(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want Topic
	}{
		{
			name: "simple",
			args: args{query: "Extract beginning of string (prefix)\n(tags: extract beginning string prefix)\n---\n\nt := string([]rune(s)[:5])"},
			want: Topic{
				Title:   "Extract beginning of string (prefix)",
				Code:    "t := string([]rune(s)[:5])",
				Checked: false,
				Tags: []Tag{
					{Name: "extract"},
					{Name: "beginning"},
					{Name: "string"},
					{Name: "prefix"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertQueryToTopic(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertQueryToTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
