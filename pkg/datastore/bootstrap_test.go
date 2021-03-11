package datastore

import (
	"reflect"
	"testing"
)

func Test_splitYamlDocs(t *testing.T) {
	type args struct {
		doc string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				doc: "",
			},
			want: nil,
		},
		{
			name: "doc w/o ---",
			args: args{
				doc: `a:
  b: asdasd
`,
			},
			want: []string{
				`a:
  b: asdasd
`,
			},
		},
		{
			name: "doc w/ ---",
			args: args{
				doc: `---
a:
  b: asdasd
`,
			},
			want: []string{
				`a:
  b: asdasd
`,
			},
		},
		{
			name: "multi doc",
			args: args{
				doc: `---
a:
  b: asdasd
---
c:
  d: asdcxx
`,
			},
			want: []string{
				`a:
  b: asdasd
`,
				`c:
  d: asdcxx
`,
			},
		},
		{
			name: "multi doc chaos",
			args: args{
				doc: `
---   
 
a:
  b: asdasd


---        
c:
  d: asdcxx
`,
			},
			want: []string{
				`a:
  b: asdasd


`,
				`c:
  d: asdcxx
`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := splitYamlDocs(tt.args.doc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitYamlDocs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
