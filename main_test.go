package main

import (
	"crypto/subtle"
	"testing"
)

func UnsafeCompare(x string, y string) (result int) {
	if len(x) != len(y) {
		return 0
	}

	for idx, i := range x {
		if string(i) != y[idx:idx+1] {
			return 0
		}
	}
	return 1
}

func SafeCompare(x string, y string) (result int) {
	xbytes := []byte(x)
	ybytes := []byte(y)
	return subtle.ConstantTimeCompare(xbytes, ybytes)
}

func benchmarkUnsafeCompare(x string, y string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		UnsafeCompare(x, y)
	}
}

func benchmarkSafeCompare(x string, y string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeCompare(x, y)
	}
}

func BenchmarkUnsafeCompareSame(b *testing.B) { benchmarkUnsafeCompare("1234567890", "1234567890", b) }

func BenchmarkUnsafeCompareLen1(b *testing.B) { benchmarkUnsafeCompare("1234567890X", "1234567890", b) }
func BenchmarkUnsafeCompareLen2(b *testing.B) {
	benchmarkUnsafeCompare("1234567890XX", "1234567890", b)
}
func BenchmarkUnsafeCompareLen3(b *testing.B) {
	benchmarkUnsafeCompare("1234567890XXX", "1234567890", b)
}

func BenchmarkUnsafeCompare0(b *testing.B) { benchmarkUnsafeCompare("XXXXXXXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare1(b *testing.B) { benchmarkUnsafeCompare("1XXXXXXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare2(b *testing.B) { benchmarkUnsafeCompare("12XXXXXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare3(b *testing.B) { benchmarkUnsafeCompare("123XXXXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare4(b *testing.B) { benchmarkUnsafeCompare("1234XXXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare5(b *testing.B) { benchmarkUnsafeCompare("12345XXXXX", "1234567890", b) }
func BenchmarkUnsafeCompare6(b *testing.B) { benchmarkUnsafeCompare("123456XXXX", "1234567890", b) }
func BenchmarkUnsafeCompare7(b *testing.B) { benchmarkUnsafeCompare("1234567XXX", "1234567890", b) }
func BenchmarkUnsafeCompare8(b *testing.B) { benchmarkUnsafeCompare("12345678XX", "1234567890", b) }
func BenchmarkUnsafeCompare9(b *testing.B) { benchmarkUnsafeCompare("123456789X", "1234567890", b) }

func BenchmarkSafeCompareSame(b *testing.B) { benchmarkSafeCompare("1234567890", "1234567890", b) }
func BenchmarkSafeCompare0(b *testing.B)    { benchmarkSafeCompare("XXXXXXXXXX", "1234567890", b) }
func BenchmarkSafeCompare1(b *testing.B)    { benchmarkSafeCompare("1XXXXXXXXX", "1234567890", b) }
func BenchmarkSafeCompare2(b *testing.B)    { benchmarkSafeCompare("12XXXXXXXX", "1234567890", b) }
func BenchmarkSafeCompare3(b *testing.B)    { benchmarkSafeCompare("123XXXXXXX", "1234567890", b) }
func BenchmarkSafeCompare4(b *testing.B)    { benchmarkSafeCompare("1234XXXXXX", "1234567890", b) }
func BenchmarkSafeCompare5(b *testing.B)    { benchmarkSafeCompare("12345XXXXX", "1234567890", b) }
func BenchmarkSafeCompare6(b *testing.B)    { benchmarkSafeCompare("123456XXXX", "1234567890", b) }
func BenchmarkSafeCompare7(b *testing.B)    { benchmarkSafeCompare("1234567XXX", "1234567890", b) }
func BenchmarkSafeCompare8(b *testing.B)    { benchmarkSafeCompare("12345678XX", "1234567890", b) }
func BenchmarkSafeCompare9(b *testing.B)    { benchmarkSafeCompare("123456789X", "1234567890", b) }
