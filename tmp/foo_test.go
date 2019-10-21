package main

import "testing"

func TestHelloWorld(t *testing.T) {
	t.Parallel()
	t.Run("foo", func(t *testing.T) {
		t.Run("bar", func(t *testing.T) {
			t.Run("quux", func(t *testing.T) {
				t.Errorf("ding")
			})
		})
	})
}

type Runner interface {
	Run(string, func(*testing.T)) bool
}

type Funner interface {
	Runner
}

func val() string {
	return "zing"
}

func TestX(t *testing.T) {
	var r Funner = t

	f := func(t *testing.T) {
		t.Run(val(), func(t *testing.T) {
			t.Parallel()
			t.Errorf("grr")
		})
	}

	returnf := func(x int) func(t *testing.T) {
		return f
	}

	ptrReturnf := &returnf
	r.Run(val(), (*ptrReturnf)(31337))
}

func TestSimple(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		t.Errorf("Simple")
	})
}

func TestComplex(t *testing.T) {
	d := t.Run

	d("complex", func(tt *testing.T) {
		t.Errorf("Complex :(")
	})
}

func TestComplex2(t *testing.T) {
	(*testing.T).Run(t, "complex2", func(tt *testing.T) {
		tt.Errorf("Complex2 :(")
	})
}

func TestComplex3(t *testing.T) {
	type TP = *testing.T

	TP.Run(t, "complex3", func(tt *testing.T) {
		tt.Errorf("Complex3 :(")
	})
}

func TestComplex4(t *testing.T) {

	testingval := func() *testing.T {
		return t
	}

	t.Run("complex4", func(tt *testing.T) {
		testingval().Errorf("lollipops")
		tt.Errorf("Complex4 :(")
	})
}

func TestComplex5(t *testing.T) {

	ef := func(*testing.T) func(format string, args ...interface{}) {
		return t.Errorf
	}

	t.Run("complex5", func(tt *testing.T) {
		ef(t)("lollipops")
		tt.Errorf("Complex5 :(")
	})
}

func TestComplex6(t *testing.T) {
	t.Run("complex6", func(tt *testing.T) {
		(*testing.T).Errorf(t, "lollipops")
		tt.Errorf("Complex6 :(")
	})
}
