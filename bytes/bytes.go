package bytes

import (
	"bytes"
	"sort"
)

func Sort(a [][]byte) {
	sort.Sort(byteSlices(a))
}

func SearchBytes(a [][]byte, x []byte) int {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if bytes.Compare(a[h], x) < 0 {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func Contains(a [][]byte, x []byte) bool {
	n := SearchBytes(a, x)
	return n < len(a) && bytes.Equal(a[n], x)
}

func Union(a, b [][]byte) [][]byte {
	n := len(b)
	if len(a) > len(b) {
		n = len(a)
	}
	other := make([][]byte, 0, n)

	for {
		if len(a) > 0 && len(b) > 0 {
			if cmp := bytes.Compare(a[0], b[0]); cmp == 0 {
				other, a, b = append(other, a[0]), a[1:], b[1:]
			} else if cmp == -1 {
				other, a = append(other, a[0]), a[1:]
			} else {
				other, b = append(other, b[0]), b[1:]
			}
		} else if len(a) > 0 {
			other, a = append(other, a[0]), a[1:]
		} else if len(b) > 0 {
			other, b = append(other, b[0]), b[1:]
		} else {
			return other
		}
	}
}

func Intersect(a, b [][]byte) [][]byte {
	n := len(b)
	if len(a) > len(b) {
		n = len(a)
	}
	other := make([][]byte, 0, n)

	for len(a) > 0 && len(b) > 0 {
		if cmp := bytes.Compare(a[0], b[0]); cmp == 0 {
			other, a, b = append(other, a[0]), a[1:], b[1:]
		} else if cmp == -1 {
			a = a[1:]
		} else {
			b = b[1:]
		}
	}
	return other
}
func Clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	buf := make([]byte, len(b))
	copy(buf, b)
	return buf
}
func CloneSlice(a [][]byte) [][]byte {
	other := make([][]byte, len(a))
	for i := range a {
		other[i] = Clone(a[i])
	}
	return other
}

func IsSorted(a [][]byte) bool {
	return sort.IsSorted(byteSlices(a))
}

type byteSlices [][]byte

func (a byteSlices) Len() int           { return len(a) }
func (a byteSlices) Less(i, j int) bool { return bytes.Compare(a[i], a[j]) == -1 }
func (a byteSlices) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
