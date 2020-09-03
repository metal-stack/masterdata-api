package v1

import "testing"

func Test_onlyOneOfPtrsSet(t *testing.T) {
	type dummy struct {
		a *string
		b *int
		c *dummy
	}
	s := "abc"
	b := 42
	no := dummy{
		nil, nil, nil,
	}
	one := dummy{
		nil, nil, &dummy{},
	}
	two := dummy{
		&s, nil, &dummy{},
	}
	all := dummy{
		&s, &b, &dummy{},
	}

	type args struct {
		ptrs []interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all nil",
			args: args{
				ptrs: []interface{}{no.a, no.b, no.c},
			},
			want: false,
		},
		{
			name: "two != nil",
			args: args{
				ptrs: []interface{}{two.a, two.b, two.c},
			},
			want: false,
		},
		{
			name: "all != nil",
			args: args{
				ptrs: []interface{}{all.a, all.b, all.c},
			},
			want: false,
		},
		{
			name: "only one != nil",
			args: args{
				ptrs: []interface{}{one.a, one.b, one.c},
			},
			want: true,
		},
		{
			name: "one != nil, value",
			args: args{
				ptrs: []interface{}{one.a, one.b, one.c, "value"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := onlyOneOfPtrsSet(tt.args.ptrs...); got != tt.want {
				t.Errorf("onlyOneOfPtrsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
