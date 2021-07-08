package mysql

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"Garyen-go/tidb/util"
)

// This file contains the data encoding and decoding functions.

// for go lint check
var (
	_          = LenEncIntSize(0)
	_          = WriteLenEncInt([]byte{1}, 0, 1)
	_          = LenNullString("")
	_          = WriteNullString([]byte{1}, 0, "")
	_          = writeEOFString([]byte{1}, 0, "")
	_          = WriteByte([]byte{1}, 0, byte(0))
	_          = WriteUint16([]byte{1, 2}, 0, 0)
	_          = WriteZeroes([]byte{1}, 0, 0)
	_, _, _, _ = ReadLenEncInt([]byte{1}, 0)
	_, _       = AppendBinaryValue([]byte{1}, 0, 1)
)

//
// Encoding methods.
//
// The same assumptions are made for all the encoding functions:
// - there is enough space to write the data in the buffer. If not, we
// will panic with out of bounds.
// - all functions start writing at 'pos' in the buffer, and return the next position.

// LenEncIntSize returns the number of bytes required to encode a
// variable-length integer.
func LenEncIntSize(i uint64) int {
	switch {
	case i < 251:
		return 1
	case i < 1<<16:
		return 3
	case i < 1<<24:
		return 4
	default:
		return 9
	}
}

// WriteLenEncInt write uint64 to []byte
func WriteLenEncInt(data []byte, pos int, i uint64) int {
	switch {
	case i < 251:
		data[pos] = byte(i)
		return pos + 1
	case i < 1<<16:
		data[pos] = 0xfc
		data[pos+1] = byte(i)
		data[pos+2] = byte(i >> 8)
		return pos + 3
	case i < 1<<24:
		data[pos] = 0xfd
		data[pos+1] = byte(i)
		data[pos+2] = byte(i >> 8)
		data[pos+3] = byte(i >> 16)
		return pos + 4
	default:
		data[pos] = 0xfe
		data[pos+1] = byte(i)
		data[pos+2] = byte(i >> 8)
		data[pos+3] = byte(i >> 16)
		data[pos+4] = byte(i >> 24)
		data[pos+5] = byte(i >> 32)
		data[pos+6] = byte(i >> 40)
		data[pos+7] = byte(i >> 48)
		data[pos+8] = byte(i >> 56)
		return pos + 9
	}
}

// AppendLenEncInt append LenEncInt []byte to data
func AppendLenEncInt(data []byte, i uint64) []byte {
	switch {
	case i <= 250:
		return append(data, byte(i))

	case i <= 0xffff:
		return append(data, 0xfc, byte(i), byte(i>>8))

	case i <= 0xffffff:
		return append(data, 0xfd, byte(i), byte(i>>8), byte(i>>16))

	case i <= 0xffffffffffffffff:
		return append(data, 0xfe, byte(i), byte(i>>8), byte(i>>16), byte(i>>24),
			byte(i>>32), byte(i>>40), byte(i>>48), byte(i>>56))
	}

	return data
}

// LenNullString return length Null terminated string
func LenNullString(value string) int {
	return len(value) + 1
}

// WriteNullString write NULL terminated strign to []byte
func WriteNullString(data []byte, pos int, value string) int {
	pos += copy(data[pos:], value)
	data[pos] = 0
	return pos + 1
}

func writeEOFString(data []byte, pos int, value string) int {
	pos += copy(data[pos:], value)
	return pos
}

// WriteByte write byte to []byte
func WriteByte(data []byte, pos int, value byte) int {
	data[pos] = value
	return pos + 1
}

// WriteUint16 write uint16 to []byte
func WriteUint16(data []byte, pos int, value uint16) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	return pos + 2
}

// AppendUint16 append uint16 to []byte
func AppendUint16(data []byte, n uint16) []byte {
	data = append(data, byte(n))
	data = append(data, byte(n>>8))
	return data
}

// WriteUint32 write uint32 to []byte
func WriteUint32(data []byte, pos int, value uint32) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	data[pos+2] = byte(value >> 16)
	data[pos+3] = byte(value >> 24)
	return pos + 4
}

// AppendUint32 append uint32 to []byte
func AppendUint32(data []byte, n uint32) []byte {
	data = append(data, byte(n))
	data = append(data, byte(n>>8))
	data = append(data, byte(n>>16))
	data = append(data, byte(n>>24))
	return data
}

// WriteUint64 write uint64 to []byte
func WriteUint64(data []byte, pos int, value uint64) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	data[pos+2] = byte(value >> 16)
	data[pos+3] = byte(value >> 24)
	data[pos+4] = byte(value >> 32)
	data[pos+5] = byte(value >> 40)
	data[pos+6] = byte(value >> 48)
	data[pos+7] = byte(value >> 56)
	return pos + 8
}

// AppendLenEncStringBytes append bytes of len enc string  to data
func AppendLenEncStringBytes(data, b []byte) []byte {
	data = AppendLenEncInt(data, uint64(len(b)))
	data = append(data, b...)
	return data
}

// WriteZeroes write 0 to []byte
func WriteZeroes(data []byte, pos int, len int) int {
	for i := 0; i < len; i++ {
		data[pos+i] = 0
	}
	return pos + len
}

// ReadLenEncInt read info of len encoded int, return length, next pos(skip len self to data), is null, handle result
// https://dev.mysql.com/doc/internals/en/integer.html#packet-Protocol::FixedLengthInteger
func ReadLenEncInt(data []byte, pos int) (uint64, int, bool, bool) {
	isNull := false
	if pos >= len(data) {
		return 0, 0, isNull, false
	}
	switch data[pos] {
	// 251: NULL
	case 0xfb:
		isNull = true
		return 0, pos + 1, isNull, true
	case 0xfc:
		// Encoded in the next 2 bytes.
		if pos+2 >= len(data) {
			return 0, 0, isNull, false
		}
		return uint64(data[pos+1]) |
			uint64(data[pos+2])<<8, pos + 3, isNull, true
	case 0xfd:
		// Encoded in the next 3 bytes.
		if pos+3 >= len(data) {
			return 0, 0, isNull, false
		}
		return uint64(data[pos+1]) |
			uint64(data[pos+2])<<8 |
			uint64(data[pos+3])<<16, pos + 4, isNull, true
	case 0xfe:
		// Encoded in the next 8 bytes.
		if pos+8 >= len(data) {
			return 0, 0, isNull, false
		}
		return uint64(data[pos+1]) |
			uint64(data[pos+2])<<8 |
			uint64(data[pos+3])<<16 |
			uint64(data[pos+4])<<24 |
			uint64(data[pos+5])<<32 |
			uint64(data[pos+6])<<40 |
			uint64(data[pos+7])<<48 |
			uint64(data[pos+8])<<56, pos + 9, isNull, true
	}
	// 0-250
	return uint64(data[pos]), pos + 1, isNull, true
}

func ReadByte(data []byte, pos int) (byte, int, bool) {
	if pos >= len(data) {
		return 0, 0, false
	}

	return data[pos], pos + 1, true
}

func ReadBytes(data []byte, pos int, size int) ([]byte, int, bool) {
	if pos+size-1 >= len(data) {
		return nil, 0, false
	}

	return data[pos : pos+size], pos + size, true
}

// ReadBytesCopy returns a copy of the bytes in the packet.
// Useful to remember contents of ephemeral packets.
func ReadBytesCopy(data []byte, pos int, size int) ([]byte, int, bool) {
	if pos+size-1 >= len(data) {
		return nil, 0, false
	}

	result := make([]byte, size)
	copy(result, data[pos:pos+size])

	return result, pos + size, true
}

// AppendBinaryValue encode binary-type value of prepare binary protocol according to type of value
func AppendBinaryValue(data []byte, fieldType uint8, value interface{}) ([]byte, error) {
	// constructor phase
	var t []byte
	switch v := value.(type) {
	case int8:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case int16:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case int32:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case int64:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case int:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case uint8:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case uint16:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case uint32:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case uint64:
		t = make([]byte, 8)
		WriteUint64(t, 0, v)
	case uint:
		t = make([]byte, 8)
		WriteUint64(t, 0, uint64(v))
	case float64:
		switch fieldType {
		case TypeFloat:
			t = make([]byte, 4)
			WriteUint32(t, 0, math.Float32bits(float32(v)))
		case TypeNewDecimal:
			tmp := strconv.FormatFloat(v, 'f', -1, 64)
			t = util.Slice(tmp)
		default:
			t = make([]byte, 8)
			WriteUint64(t, 0, math.Float64bits(v))
		}
	case []byte:
		t = v
	case string:
		switch fieldType {
		case TypeDatetime, TypeTimestamp:
			if v == "0000-00-00 00:00:00" {
				t = append(t, 0)
			} else {
				ts, err := time.Parse("2006-01-02 15:04:05", v)
				if err != nil {
					var mysqlTypeStr string
					if fieldType == TypeDatetime {
						mysqlTypeStr = "TypeDatetime"
					} else {
						mysqlTypeStr = "TypeTimestamp"
					}
					return nil, fmt.Errorf("invalid %s %s", mysqlTypeStr, v)
				}

				t = append(t, 11)
				t = AppendUint16(t, uint16(ts.Year()))
				t = append(t, byte(int(ts.Month())), byte(ts.Day()), byte(ts.Hour()), byte(ts.Minute()), byte(ts.Second()))
				microseconds := uint32(ts.Nanosecond() / 1000)
				t = AppendUint32(t, microseconds)
			}
		case TypeDate:
			// format: 2006-01-02
			ts, err := time.Parse("2006-01-02", v)
			if err != nil {
				t = append(t, 0)
			} else {
				t = append(t, 4)
				t = AppendUint16(t, uint16(ts.Year()))
				t = append(t, byte(int(ts.Month())), byte(ts.Day()))
			}
		case TypeDuration:
			timeValue, err := stringToMysqlTime(v)
			if err != nil {
				return nil, err
			}

			t = mysqlTimeToBinaryResult(timeValue)
		default:
			t = util.Slice(v)
		}
	default:
		return data, fmt.Errorf("invalid type %T", value)
	}

	// append phase
	switch fieldType {
	case TypeTiny:
		data = append(data, t[0])
		return data, nil
	case TypeShort, TypeYear:
		data = append(data, t[:2]...)
		return data, nil
	case TypeFloat, TypeInt24, TypeLong:
		data = append(data, t[:4]...)
		return data, nil
	case TypeLonglong, TypeDouble:
		data = append(data, t[:8]...)
		return data, nil
	case TypeNewDecimal, TypeJSON, TypeString, TypeVarString, TypeVarchar, TypeBit, TypeTinyBlob, TypeMediumBlob, TypeLongBlob, TypeBlob:
		tmp := make([]byte, 0, len(t)+9)
		data = append(data, AppendLenEncStringBytes(tmp, t)...)
		return data, nil
	case TypeEnum, TypeSet, TypeDate, TypeDatetime, TypeDuration, TypeTimestamp, TypeNewDate:
		data = append(data, t...)
		return data, nil
	default:
		return data, fmt.Errorf("not supported field type")
	}
}
