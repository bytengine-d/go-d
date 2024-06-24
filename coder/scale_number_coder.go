package coder

import (
	"encoding/binary"
	"github.com/bytengine-d/go-d/lang"
	"math"
	"math/rand"
	"strings"
)

var (
	NumberChars    = []rune("1234567890")
	LowercaseChars = []rune("abcdefghijklmnopqrstuvwxyz")
	UppercaseChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	NumberAndLowercaseChars             = append(NumberChars, LowercaseChars...)
	LowercaseAndUppercaseChars          = append(LowercaseChars, UppercaseChars...)
	NumberAndLowercaseAndUppercaseChars = append(NumberAndLowercaseChars, UppercaseChars...)
)

const (
	NumberCharsLen    = 10
	LowercaseCharsLen = 26
	UppercaseCharsLen = 26

	NumberAndLowercaseLen             = NumberCharsLen + LowercaseCharsLen
	LowercaseAndUppercaseLen          = LowercaseCharsLen + UppercaseCharsLen
	NumberAndLowercaseAndUppercaseLen = NumberCharsLen + LowercaseCharsLen + UppercaseCharsLen
)

// region Global functions
func BuildTable(tableSize int, chares ...[]rune) string {
	totalLen := 0
	limits := make([]int32, len(chares))
	for idx, chars := range chares {
		totalLen += len(chars)
		limits[idx] = int32(totalLen)
	}
	bitset := make([]bool, totalLen)
	total := int32(totalLen)
	var idx int32 = 0
	buf := strings.Builder{}
	for i := 0; i < tableSize; i++ {
		idx = rand.Int31n(total)
		if bitset[idx] {
			i--
			continue
		}
		bitset[idx] = true
		subIdx := idx
		for j, limit := range limits {
			if idx < limit {
				buf.WriteRune(chares[j][subIdx])
				break
			}
			subIdx = idx - limit
		}
	}
	return buf.String()
}

func NewScaleConvertor(table string) *ScaleNumberCoder {
	charTable := []rune(table)
	return &ScaleNumberCoder{
		scale:     int64(len(charTable)),
		charTable: charTable,
	}
}

// endregion

// region ScaleNumberCoder
type ScaleNumberCoder struct {
	scale     int64
	charTable []rune
}

func (coder *ScaleNumberCoder) Transform(number int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(number))
	scale := coder.scale
	condition := coder.scale - 1
	charTable := coder.charTable
	var remainder int64
	buf := strings.Builder{}
	for number > condition {
		remainder = number % scale
		buf.WriteRune(charTable[remainder])
		number = number / scale
	}
	buf.WriteRune(charTable[number])

	return lang.Reverse(buf.String())
}

func (coder *ScaleNumberCoder) From(code string) int64 {
	var value int64 = 0
	var idx int64 = 0
	codeChars := []rune(code)
	codeLen := len(codeChars)
	scale := float64(coder.scale)
	charTable := coder.charTable
	for i, char := range codeChars {
		idx = lang.Index(charTable, char, func(v1, v2 rune) bool {
			return v1 == v2
		})
		value += idx * int64(math.Pow(scale, float64(codeLen-i-1)))
	}
	return value
}

// endregion
