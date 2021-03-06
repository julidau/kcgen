package kcsl

import (
	"errors"
	"math"

	"go.starlark.net/starlark"
)

var mathBuiltins = starlark.StringDict{
	"pi": starlark.Float(math.Pi),
	"sin": starlark.NewBuiltin("sin", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("sin", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Sin(float64(theta))), nil
	}),
	"sinh": starlark.NewBuiltin("sinh", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("sinh", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Sinh(float64(theta))), nil
	}),
	"asin": starlark.NewBuiltin("asin", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("asin", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Asin(float64(theta))), nil
	}),
	"asinh": starlark.NewBuiltin("asinh", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("asinh", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Asinh(float64(theta))), nil
	}),
	"cos": starlark.NewBuiltin("cos", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("cos", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Cos(float64(theta))), nil
	}),
	"tan": starlark.NewBuiltin("tan", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("tan", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Tan(float64(theta))), nil
	}),
	"tanh": starlark.NewBuiltin("tanh", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("tanh", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Tanh(float64(theta))), nil
	}),
	"atan": starlark.NewBuiltin("atan", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var theta starlark.Float
		if err := starlark.UnpackArgs("atan", args, kwargs, "theta", &theta); err != nil {
			return starlark.None, err
		}
		return starlark.Float(math.Atan(float64(theta))), nil
	}),

	"shl": starlark.NewBuiltin("shl", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var base starlark.Int
		var shiftAmount int
		if err := starlark.UnpackArgs("shl", args, kwargs, "base", &base, "shift", &shiftAmount); err != nil {
			return starlark.None, err
		}
		b, ok := base.Uint64()
		if !ok {
			return starlark.None, errors.New("cannot represent base as unsigned integer")
		}
		return starlark.MakeUint64(b << uint(shiftAmount)), nil
	}),
	"shr": starlark.NewBuiltin("shr", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var base starlark.Int
		var shiftAmount int
		if err := starlark.UnpackArgs("shr", args, kwargs, "base", &base, "shift", &shiftAmount); err != nil {
			return starlark.None, err
		}
		b, ok := base.Uint64()
		if !ok {
			return starlark.None, errors.New("cannot represent base as unsigned integer")
		}
		return starlark.MakeUint64(b >> uint(shiftAmount)), nil
	}),
	"_not": starlark.NewBuiltin("_not", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var base starlark.Int
		if err := starlark.UnpackArgs("_not", args, kwargs, "base", &base); err != nil {
			return starlark.None, err
		}
		b, ok := base.Uint64()
		if !ok {
			return starlark.None, errors.New("cannot represent base as unsigned integer")
		}
		return starlark.MakeUint64(^b), nil
	}),
	"_and": starlark.NewBuiltin("_and", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var o1, o2 starlark.Int
		if err := starlark.UnpackArgs("_and", args, kwargs, "op1", &o1, "op2", &o2); err != nil {
			return starlark.None, err
		}
		i1, ok := o1.Uint64()
		if !ok {
			return starlark.None, errors.New("cannot represent op1 as unsigned integer")
		}
		i2, ok := o2.Uint64()
		if !ok {
			return starlark.None, errors.New("cannot represent op2 as unsigned integer")
		}
		return starlark.MakeUint64(i1 & i2), nil
	}),
}
