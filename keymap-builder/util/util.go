package util

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"slices"
)

func Reduce[T, A any](s []T, acc A, f func(A, T) A) A {
	for _, x := range s {
		acc = f(acc, x)
	}

	return acc
}

func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	for i := range s {
		out[i] = f(s[i])
	}

	return out
}

func MapString[T any](s string, f func(byte) T) []T {
	out := make([]T, len(s))
	for i := range s {
		out[i] = f(s[i])
	}

	return out
}

func MapEnumerated[T, U any](s []T, f func(int, T) U) []U {
	us := make([]U, len(s))
	for i := range s {
		us[i] = f(i, s[i])
	}

	return us
}

func MapToAny[T any](args []T) []any {
	return Map(args, func(a T) any { return a })
}

func ToString(x any) string {
	if s, ok := x.(string); ok {
		return s
	}

	return fmt.Sprintf("%s", x)
}

func AsString[T ~string](x T) string {
	return string(x)
}

type Lesser[K any] interface {
	comparable
	Less(K) int
}

func LesserFn[T Lesser[T]](a, b T) int {
	return a.Less(b)
}

func SortedMapFunc[K comparable, V any](m map[K]V, fn func(K, K) int) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.SortFunc(keys, fn)
	return func(yield func(K, V) bool) {
		for _, key := range keys {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}

func SortedMap[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return func(yield func(K, V) bool) {
		for _, key := range keys {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Panicif(cond bool, extras ...any) {
	if cond {
		if len(extras) > 0 {
			if v, ok := extras[0].(string); ok {
				panic(fmt.Sprintf(v, extras[1:]...))
			}
		}
		panic("condition failed")
	}
}
