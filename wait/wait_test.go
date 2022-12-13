package wait

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	oops        = errors.New("oops")
	boolFnTrue  = func() bool { return true }
	boolFnFalse = func() bool { return false }
	errFnNil    = func() error { return nil }
	errFnNotNil = func() error { return oops }
	tFnNil      = func() (bool, error) { return true, nil }
	tFnNotNil   = func() (bool, error) { return false, oops }
)

func eqErr(t *testing.T, exp, err error) {
	t.Helper()

	if exp == nil || err == nil {
		if !errors.Is(exp, err) {
			t.Fatalf("exp: %v, err: %v", exp, err)
		}
		return
	}
	expect := exp.Error()
	actual := err.Error()
	if expect != actual {
		t.Fatalf("exp: %s, err: %s", expect, actual)
	}
}

func TestNoFunction(t *testing.T) {
	t.Parallel()

	ctx := On()
	err := ctx.Run()
	if !errors.Is(err, ErrNoFunction) {
		t.Fatalf("exp: %v, err: %v", ErrNoFunction, err)
	}
}

func TestBoolFunc(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opts []Option
		exp  error
	}{
		{
			name: "defaults ok",
			opts: []Option{BoolFunc(boolFnTrue)},
		},
		{
			name: "defaults fail",
			opts: []Option{BoolFunc(boolFnFalse)},
			exp:  ErrTimeoutExceeded,
		},
		{
			name: "iterations exceeded",
			opts: []Option{
				BoolFunc(boolFnFalse),
				Attempts(3),
			},
			exp: ErrAttemptsExceeded,
		},
		{
			name: "short timeout",
			opts: []Option{
				BoolFunc(boolFnFalse),
				Timeout(100 * time.Millisecond),
			},
			exp: ErrTimeoutExceeded,
		},
		{
			name: "short gap",
			opts: []Option{
				BoolFunc(boolFnFalse),
				Attempts(10),
				Gap(1 * time.Millisecond),
			},
			exp: ErrAttemptsExceeded,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := On(tc.opts...)
			err := c.Run()
			eqErr(t, tc.exp, err)
		})
	}
}

func TestErrorFunc(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opts []Option
		exp  error
	}{
		{
			name: "defaults ok",
			opts: []Option{ErrorFunc(errFnNil)},
		},
		{
			name: "defaults fail",
			opts: []Option{ErrorFunc(errFnNotNil)},
			exp:  fmt.Errorf("%v: %w", ErrTimeoutExceeded, oops),
		},
		{
			name: "attempts exceeded",
			opts: []Option{
				ErrorFunc(errFnNotNil),
				Attempts(3),
			},
			exp: fmt.Errorf("%v: %w", ErrAttemptsExceeded, oops),
		},
		{
			name: "short timeout",
			opts: []Option{
				ErrorFunc(errFnNotNil),
				Attempts(1000),
				Timeout(100 * time.Millisecond),
			},
			exp: fmt.Errorf("%v: %w", ErrTimeoutExceeded, oops),
		},
		{
			name: "short gap",
			opts: []Option{
				ErrorFunc(errFnNotNil),
				Attempts(10),
				Gap(1 * time.Millisecond),
			},
			exp: fmt.Errorf("%v: %w", ErrAttemptsExceeded, oops),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := On(tc.opts...)
			err := c.Run()
			eqErr(t, tc.exp, err)
		})
	}
}

func TestTestFunc(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opts []Option
		exp  error
	}{
		{
			name: "defaults ok",
			opts: []Option{TestFunc(tFnNil)},
		},
		{
			name: "defaults fail",
			opts: []Option{TestFunc(tFnNotNil)},
			exp:  fmt.Errorf("%v: %w", ErrTimeoutExceeded, oops),
		},
		{
			name: "default fail without error",
			opts: []Option{
				TestFunc(func() (bool, error) {
					return false, nil
				}),
			},
			exp: fmt.Errorf("%v: %w", ErrTimeoutExceeded, ErrConditionUnsatisfied),
		},
		{
			name: "attempts exceeded",
			opts: []Option{
				TestFunc(tFnNotNil),
				Attempts(3),
			},
			exp: fmt.Errorf("%v: %w", ErrAttemptsExceeded, oops),
		},
		{
			name: "short timeout",
			opts: []Option{
				TestFunc(tFnNotNil),
				Attempts(1000),
				Timeout(100 * time.Millisecond),
			},
			exp: fmt.Errorf("%v: %w", ErrTimeoutExceeded, oops),
		},
		{
			name: "short gap",
			opts: []Option{
				TestFunc(tFnNotNil),
				Attempts(10),
				Gap(1 * time.Millisecond),
			},
			exp: fmt.Errorf("%v: %w", ErrAttemptsExceeded, oops),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := On(tc.opts...)
			err := c.Run()
			eqErr(t, tc.exp, err)
		})
	}
}
