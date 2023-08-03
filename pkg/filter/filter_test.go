package filter

import "testing"

func TestInverse(t *testing.T) {
	t.Run("Inverse of function that filters out odd numbers", func(t *testing.T) {
		f := Inverse[int](func(i int) bool { return i%2 != 0 })
		for i, want := range []bool{true, false, true, false, true} {
			if got := f(i); got != want {
				t.Errorf("f(%d) = %v, want %v", i, got, want)
			}
		}
	})

	t.Run("Inverse of function that filters out empty strings", func(t *testing.T) {
		g := Inverse[string](func(s string) bool { return s == "" })
		for _, tc := range []struct {
			in   string
			want bool
		}{
			{in: "nonempty", want: true},
			{in: "", want: false},
		} {
			if got := g(tc.in); got != tc.want {
				t.Errorf("g(%q) = %v, want %v", tc.in, got, tc.want)
			}
		}
	})
}
