package main

import (
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
	us := make([]U, len(s))
	for i := range s {
		us[i] = f(s[i])
	}

	return us
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

func MapToAnyStatic[T, V any](args []T, v V) []any {
	return Map(args, func(_ T) any { return v })
}

func toString(x any) string {
	return fmt.Sprintf("%s", x)
}

func asString[T ~string](x T) string {
	return string(x)
}

type Lesser[K any] interface {
	comparable
	Less(K) int
}

func SortedMapKV[K Lesser[K], V any](m map[K]V) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.SortFunc(keys, func(a, b K) int {
		return a.Less(b)
	})
	return func(yield func(K, V) bool) {
		for _, key := range keys {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}
