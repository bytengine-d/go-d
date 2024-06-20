package coder

type NumberCoder interface {
	Transform(number int64) string
	From(code string) int64
}
