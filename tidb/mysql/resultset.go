package mysql

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"Garyen-go/tidb/util"
)

type RowData []byte

func (p RowData) Parse(f []*Field, binary bool) ([]interface{}, error) {
	if binary {
		return p.ParseBinary(f)
	}

	return p.ParseText(f)
}

func (p RowData) ParseText(f []*Field) ([]interface{}, error) {
	data := make([]interface{}, len(f))

	var err error
	var v []byte
	var isNull, isUnsigned bool
	var pos = 0
	var n int

	for i := range f {
		v, isNull, n, err = LengthEncodedString(p[pos:])
		if err != nil {
			return nil, err
		}

		pos += n

		if isNull {
			data[i] = nil
		} else {
			isUnsigned = f[i].Flag&uint16(UnsignedFlag) > 0
			switch f[i].Type {
			case TypeTiny, TypeShort, TypeLong, TypeInt24,
				TypeLonglong, TypeYear:
				if isUnsigned {
					data[i], err = strconv.ParseUint(string(v), 10, 64)
				} else {
					data[i], err = strconv.ParseInt(string(v), 10, 64)
				}
			case TypeFloat, TypeDouble, TypeNewDecimal:
				data[i], err = strconv.ParseFloat(string(v), 64)
			case TypeVarchar, TypeVarString,
				TypeString, TypeDatetime,
				TypeDate, TypeDuration, TypeTimestamp, TypeJSON:
				data[i] = string(v)
			default:
				data[i] = v
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return data, nil
}

func (p RowData) ParseBinary(f []*Field) ([]interface{}, error) {
	data := make([]interface{}, len(f))

	if p[0] != OKHeader {
		return nil, ErrMalformPacket
	}

	pos := 1 + ((len(f) + 7 + 2) >> 3)

	nullBitmap := p[1:pos]

	var isUnsigned bool
	var isNull bool
	var n int
	var err error
	var v []byte

	for i := range data {
		if nullBitmap[(i+2)/8]&(1<<(uint(i+2)%8)) > 0 {
			data[i] = nil
			continue
		}

		isUnsigned = f[i].Flag&uint16(UnsignedFlag) > 0

		switch f[i].Type {
		case TypeNull:
			data[i] = nil
			continue

		case TypeTiny:
			if isUnsigned {
				data[i] = uint64(p[pos])
			} else {
				data[i] = int64(p[pos])
			}
			pos++
			continue

		case TypeShort, TypeYear:
			if isUnsigned {
				data[i] = uint64(binary.LittleEndian.Uint16(p[pos : pos+2]))
			} else {
				var n int16
				err = binary.Read(bytes.NewBuffer(p[pos:pos+2]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = int64(n)
			}
			pos += 2
			continue

		case TypeInt24, TypeLong:
			if isUnsigned {
				data[i] = uint64(binary.LittleEndian.Uint32(p[pos : pos+4]))
			} else {
				var n int32
				err = binary.Read(bytes.NewBuffer(p[pos:pos+4]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = int64(n)
			}
			pos += 4
			continue

		case TypeLonglong:
			if isUnsigned {
				data[i] = binary.LittleEndian.Uint64(p[pos : pos+8])
			} else {
				var n int64
				err = binary.Read(bytes.NewBuffer(p[pos:pos+8]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = n
			}
			pos += 8
			continue

		case TypeFloat:
			// data[i] = float64(math.Float32frombits(binary.LittleEndian.Uint32(p[pos : pos+4])))
			var n float32
			err = binary.Read(bytes.NewBuffer(p[pos:pos+4]), binary.LittleEndian, &n)
			if err != nil {
				return nil, err
			}
			data[i] = float64(n)
			pos += 4
			continue

		case TypeDouble:
			var n float64
			err = binary.Read(bytes.NewBuffer(p[pos:pos+8]), binary.LittleEndian, &n)
			if err != nil {
				return nil, err
			}
			data[i] = n
			pos += 8
			continue

		case TypeDecimal, TypeNewDecimal, TypeVarchar,
			TypeBit, TypeEnum, TypeSet, TypeTinyBlob,
			TypeMediumBlob, TypeLongBlob, TypeBlob,
			TypeVarString, TypeString, TypeGeometry, TypeJSON:
			v, isNull, n, err = LengthEncodedString(p[pos:])
			pos += n
			if err != nil {
				return nil, err
			}

			if !isNull {
				data[i] = v
				continue
			} else {
				data[i] = nil
				continue
			}
		case TypeDate, TypeNewDate:
			var num uint64
			num, isNull, n = LengthEncodedInt(p[pos:])

			pos += n

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryDate(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		case TypeTimestamp, TypeDatetime:
			var num uint64
			num, isNull, n = LengthEncodedInt(p[pos:])

			pos += n

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryDateTime(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		case TypeDuration:
			var num uint64
			num, isNull, n = LengthEncodedInt(p[pos:])

			pos += n

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryTime(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("Stmt Unknown FieldType %d %s", f[i].Type, f[i].Name)
		}
	}

	return data, nil
}

type Result struct {
	Status uint16

	InsertId     uint64
	AffectedRows uint64

	*Resultset

	Info []byte
}

type Resultset struct {
	Fields     []*Field
	FieldNames map[string]int
	Values     [][]interface{}

	RowDatas []RowData
}

func (r *Resultset) RowNumber() int {
	return len(r.Values)
}

func (r *Resultset) ColumnNumber() int {
	return len(r.Fields)
}

func (r *Resultset) GetValue(row, column int) (interface{}, error) {
	if row >= len(r.Values) || row < 0 {
		return nil, fmt.Errorf("invalid row index %d", row)
	}

	if column >= len(r.Fields) || column < 0 {
		return nil, fmt.Errorf("invalid column index %d", column)
	}

	return r.Values[row][column], nil
}

func (r *Resultset) NameIndex(name string) (int, error) {
	if column, ok := r.FieldNames[name]; ok {
		return column, nil
	}

	return 0, fmt.Errorf("invalid field name %s", name)
}

func (r *Resultset) GetValueByName(row int, name string) (interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetValue(row, column)
}

func (r *Resultset) IsNull(row, column int) (bool, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return false, err
	}

	return d == nil, nil
}

func (r *Resultset) IsNullByName(row int, name string) (bool, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return false, err
	}

	return r.IsNull(row, column)
}

func (r *Resultset) GetUint(row, column int) (uint64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case uint64:
		return v, nil
	case int64:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case []byte:
		return strconv.ParseUint(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

func (r *Resultset) GetUintByName(row int, name string) (uint64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}

	return r.GetUint(row, column)
}

func (r *Resultset) GetIntByName(row int, name string) (int64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}

	return r.GetInt(row, column)
}

func (r *Resultset) GetInt(row, column int) (int64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case uint64:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

func (r *Resultset) GetFloat(row, column int) (float64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case float64:
		return v, nil
	case uint64:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

func (r *Resultset) GetFloatByName(row int, name string) (float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}

	return r.GetFloat(row, column)
}

func (r *Resultset) GetString(row, column int) (string, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return "", err
	}

	switch v := d.(type) {
	case string:
		return v, nil
	case []byte:
		return util.String(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("data type is %T", v)
	}
}

func (r *Resultset) GetStringByName(row int, name string) (string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return "", err
	}

	return r.GetString(row, column)
}

// BuildBinaryResultset build binary resultset
// https://dev.mysql.com/doc/internals/en/binary-protocol-resultset.html
func BuildBinaryResultset(fields []*Field, values [][]interface{}) (*Resultset, error) {
	r := new(Resultset)
	r.Fields = make([]*Field, len(fields))
	for i := range fields {
		r.Fields[i] = fields[i]
	}

	bitmapLen := (len(fields) + 7 + 2) >> 3
	for i, v := range values {
		if len(v) != len(r.Fields) {
			return nil, fmt.Errorf("row %d has %d columns not equal %d", i, len(v), len(r.Fields))
		}

		var row []byte
		nullBitMap := make([]byte, bitmapLen)
		row = append(row, 0)
		row = append(row, nullBitMap...)
		for j, rowVal := range v {
			if rowVal == nil {
				bytePos := (j + 2) / 8
				bitPos := byte((j + 2) % 8)
				nullBitMap[bytePos] |= 1 << bitPos
				continue
			}

			var err error
			row, err = AppendBinaryValue(row, r.Fields[j].Type, rowVal)
			if err != nil {
				return nil, err
			}
		}
		copy(row[1:], nullBitMap)
		r.RowDatas = append(r.RowDatas, row)
	}

	return r, nil
}
