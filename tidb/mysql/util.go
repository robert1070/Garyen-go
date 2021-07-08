package mysql

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"time"
	"unicode/utf8"
)

// Abs int64 abs
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func Pstack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)

	return string(buf[0:n])
}

// CalcPassword check scrambled password received from client.
func CalcPassword(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	_, _ = crypt.Write(password)
	stage1 := crypt.Sum(nil)

	// scrambleHash = SHA1(scramble + SHA1(stage1Hash))
	// inner Hash
	crypt.Reset()
	_, _ = crypt.Write(stage1)
	hash := crypt.Sum(nil)

	// outer Hash
	crypt.Reset()
	_, _ = crypt.Write(scramble)
	_, _ = crypt.Write(hash)
	scramble = crypt.Sum(nil)

	// token = scrambleHash XOR stage1Hash
	for i := range scramble {
		scramble[i] ^= stage1[i]
	}

	return scramble
}

// CalcCachingSha2Password: Hash password using MySQL 8+ method (SHA256)
func CalcCachingSha2Password(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// XOR(SHA256(password), SHA256(SHA256(SHA256(password)), scramble))

	crypt := sha256.New()
	_, _ = crypt.Write(password)
	message1 := crypt.Sum(nil)

	crypt.Reset()
	_, _ = crypt.Write(message1)
	message1Hash := crypt.Sum(nil)

	crypt.Reset()
	_, _ = crypt.Write(message1Hash)
	_, _ = crypt.Write(scramble)
	message2 := crypt.Sum(nil)

	for i := range message1 {
		message1[i] ^= message2[i]
	}

	return message1
}

// seed must be in the range of ascii
func RandomBuf(size int) []byte {
	buf := make([]byte, size)

	rand.Seed(time.Now().UTC().UnixNano())

	min, max := 30, 127
	for i := 0; i < size; i++ {
		buf[i] = byte(min + rand.Intn(max-min))
	}

	return buf
}

func LengthEncodedInt(b []byte) (num uint64, isNull bool, n int) {
	if len(b) == 0 {
		return 0, true, 1
	}

	switch b[0] {
	// 251: NULL
	case 0xfb:
		return 0, true, 1

		// 252: value of following 2
	case 0xfc:
		return uint64(b[1]) | uint64(b[2])<<8, false, 3

		// 253: value of following 3
	case 0xfd:
		return uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16, false, 4

		// 254: value of following 8
	case 0xfe:
		return uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16 |
				uint64(b[4])<<24 | uint64(b[5])<<32 | uint64(b[6])<<40 |
				uint64(b[7])<<48 | uint64(b[8])<<56,
			false, 9
	}

	// 0-250: value of first byte
	return uint64(b[0]), false, 1
}

func PutLengthEncodedInt(n uint64) []byte {
	switch {
	case n <= 250:
		return []byte{byte(n)}

	case n <= 0xffff:
		return []byte{0xfc, byte(n), byte(n >> 8)}

	case n <= 0xffffff:
		return []byte{0xfd, byte(n), byte(n >> 8), byte(n >> 16)}

	case n <= 0xffffffffffffffff:
		return []byte{0xfe, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24),
			byte(n >> 32), byte(n >> 40), byte(n >> 48), byte(n >> 56)}
	}
	return nil
}

// LengthEncodedString returns the string read as a bytes slice, whether the value is NULL,
// the number of bytes read and an error, in case the string is longer than
// the input slice
func LengthEncodedString(b []byte) ([]byte, bool, int, error) {
	// Get length
	num, isNull, n := LengthEncodedInt(b)
	if num < 1 {
		return b[n:n], isNull, n, nil
	}

	n += int(num)

	// Check data length
	if len(b) >= n {
		return b[n-int(num) : n : n], false, n, nil
	}
	return nil, false, n, io.EOF
}

func SkipLengthEncodedString(b []byte) (int, error) {
	// Get length
	num, _, n := LengthEncodedInt(b)
	if num < 1 {
		return n, nil
	}

	n += int(num)

	// Check data length
	if len(b) >= n {
		return n, nil
	}

	return n, io.EOF
}

func PutLengthEncodedString(b []byte) []byte {
	data := make([]byte, 0, len(b)+9)
	data = append(data, PutLengthEncodedInt(uint64(len(b)))...)
	data = append(data, b...)

	return data
}

func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

func Uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}

func Uint64ToBytes(n uint64) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32),
		byte(n >> 40),
		byte(n >> 48),
		byte(n >> 56),
	}
}

func FormatBinaryDate(n int, data []byte) ([]byte, error) {
	switch n {
	case 0:
		return []byte("0000-00-00"), nil
	case 4:
		return []byte(fmt.Sprintf("%04d-%02d-%02d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3])), nil
	default:
		return nil, fmt.Errorf("invalid date packet length %d", n)
	}
}

func FormatBinaryDateTime(n int, data []byte) ([]byte, error) {
	switch n {
	case 0:
		return []byte("0000-00-00 00:00:00"), nil
	case 4:
		return []byte(fmt.Sprintf("%04d-%02d-%02d 00:00:00",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3])), nil
	case 7:
		return []byte(fmt.Sprintf(
			"%04d-%02d-%02d %02d:%02d:%02d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3],
			data[4],
			data[5],
			data[6])), nil
	case 11:
		return []byte(fmt.Sprintf(
			"%04d-%02d-%02d %02d:%02d:%02d.%06d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3],
			data[4],
			data[5],
			data[6],
			binary.LittleEndian.Uint32(data[7:11]))), nil
	default:
		return nil, fmt.Errorf("invalid datetime packet length %d", n)
	}
}

func FormatBinaryDateTime2(n int, data []byte) ([]byte, error) {
	if n == 0 {
		return []byte("0000-00-00"), nil
	}

	var sign byte
	if data[0] == 1 {
		sign = byte('-')
	}

	switch n {
	case 8:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
		)), nil
	case 12:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d.%06d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
			binary.LittleEndian.Uint32(data[8:12]),
		)), nil
	default:
		return nil, fmt.Errorf("invalid time packet length %d", n)
	}
}

func FormatBinaryTime(n int, data []byte) ([]byte, error) {
	if n == 0 {
		return []byte("0000-00-00"), nil
	}

	var sign byte
	if data[0] == 1 {
		sign = byte('-')
	}

	switch n {
	case 8:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
		)), nil
	case 12:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d.%06d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
			binary.LittleEndian.Uint32(data[8:12]),
		)), nil
	default:
		return nil, fmt.Errorf("invalid time packet length %d", n)
	}
}

var (
	DoNotEscape = byte(255)
	EncodeMap   [256]byte
)

func Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))

	for i, w := 0, 0; i < len(sql); i += w {
		runeValue, width := utf8.DecodeRuneInString(sql[i:])
		if c := EncodeMap[byte(runeValue)]; c == DoNotEscape {
			dest = append(dest, sql[i:i+width]...)
		} else {
			dest = append(dest, '\\', c)
		}

		w = width
	}

	return string(dest)
}

var encodeRef = map[byte]byte{
	'\x00': '0',
	'\'':   '\'',
	'"':    '"',
	'\b':   'b',
	'\n':   'n',
	'\r':   'r',
	'\t':   't',
	26:     'Z', // ctl-Z
	'\\':   '\\',
}

func init() {
	for i := range EncodeMap {
		EncodeMap[i] = DoNotEscape
	}

	for i := range EncodeMap {
		if to, ok := encodeRef[byte(i)]; ok {
			EncodeMap[byte(i)] = to
		}
	}
}

type lengthAndDecimal struct {
	length  int
	decimal int
}

// defaultLengthAndDecimal provides default Flen and Decimal for fields
// from CREATE TABLE when they are unspecified.
var defaultLengthAndDecimal = map[byte]lengthAndDecimal{
	TypeBit:        {1, 0},
	TypeTiny:       {4, 0},
	TypeShort:      {6, 0},
	TypeInt24:      {9, 0},
	TypeLong:       {11, 0},
	TypeLonglong:   {20, 0},
	TypeDouble:     {22, -1},
	TypeFloat:      {12, -1},
	TypeNewDecimal: {11, 0},
	TypeDuration:   {10, 0},
	TypeDate:       {10, 0},
	TypeTimestamp:  {19, 0},
	TypeDatetime:   {19, 0},
	TypeYear:       {4, 0},
	TypeString:     {1, 0},
	TypeVarchar:    {5, 0},
	TypeVarString:  {5, 0},
	TypeTinyBlob:   {255, 0},
	TypeBlob:       {65535, 0},
	TypeMediumBlob: {16777215, 0},
	TypeLongBlob:   {4294967295, 0},
	TypeJSON:       {4294967295, 0},
	TypeNull:       {0, 0},
	TypeSet:        {-1, 0},
	TypeEnum:       {-1, 0},
}

// IsIntegerType indicate whether tp is an integer type.
func IsIntegerType(tp byte) bool {
	switch tp {
	case TypeTiny, TypeShort, TypeInt24, TypeLong, TypeLonglong:
		return true
	}
	return false
}

// GetDefaultFieldLengthAndDecimal returns the default display length (flen) and decimal length for column.
// Call this when no Flen assigned in ddl.
// or column value is calculated from an expression.
// For example: "select count(*) from t;", the column type is int64 and Flen in ResultField will be 21.
// See https://dev.mysql.com/doc/refman/5.7/en/storage-requirements.html
func GetDefaultFieldLengthAndDecimal(tp byte) (flen int, decimal int) {
	val, ok := defaultLengthAndDecimal[tp]
	if ok {
		return val.length, val.decimal
	}
	return -1, -1
}

// defaultLengthAndDecimal provides default Flen and Decimal for fields
// from CAST when they are unspecified.
var defaultLengthAndDecimalForCast = map[byte]lengthAndDecimal{
	TypeString:     {0, -1}, // Flen & Decimal differs.
	TypeDate:       {10, 0},
	TypeDatetime:   {19, 0},
	TypeNewDecimal: {11, 0},
	TypeDuration:   {10, 0},
	TypeLonglong:   {22, 0},
	TypeJSON:       {4194304, 0}, // Flen differs.
}

// GetDefaultFieldLengthAndDecimalForCast returns the default display length (flen) and decimal length for casted column
// when flen or decimal is not specified.
func GetDefaultFieldLengthAndDecimalForCast(tp byte) (flen int, decimal int) {
	val, ok := defaultLengthAndDecimalForCast[tp]
	if ok {
		return val.length, val.decimal
	}
	return -1, -1
}
