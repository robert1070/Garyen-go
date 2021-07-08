package util

import (
	"reflect"
	"unsafe"
)

// no copy to change slice to string
// use your own risk
func String(b []byte) (s string) {
	pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ps := (*reflect.StringHeader)(unsafe.Pointer(&s))

	ps.Data = pb.Data
	ps.Len = pb.Len

	return
}

// no copy to change string to slice
// use your own risk
func Slice(s string) (b []byte) {
	pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ps := (*reflect.StringHeader)(unsafe.Pointer(&s))

	pb.Data = ps.Data
	pb.Len = ps.Len
	pb.Cap = ps.Len

	return
}
