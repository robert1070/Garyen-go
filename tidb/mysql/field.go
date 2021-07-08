package mysql

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeValue mysql time value
type TimeValue struct {
	IsNegative  bool
	Day         int
	Hour        int
	Minute      int
	Second      int
	Microsecond int
}

// IsNull check TimeValue if null
func (m *TimeValue) IsNull() bool {
	return m.Day == 0 && m.Hour == 0 && m.Minute == 0 && m.Second == 0 && m.Microsecond == 0
}

type FieldData []byte

type Field struct {
	Data         FieldData
	Schema       []byte
	Table        []byte
	OrgTable     []byte
	Name         []byte
	OrgName      []byte
	Charset      uint16
	ColumnLength uint32
	Type         uint8
	Flag         uint16
	Decimal      uint8

	DefaultValueLength uint64
	DefaultValue       []byte
}

func (p FieldData) Parse() (f *Field, err error) {
	f = new(Field)

	data := make([]byte, len(p))
	copy(data, p)
	f.Data = data

	var n int
	pos := 0
	// skip catelog, always def
	n, err = SkipLengthEncodedString(p)
	if err != nil {
		return
	}
	pos += n

	// schema
	f.Schema, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// table
	f.Table, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// org_table
	f.OrgTable, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// name
	f.Name, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// org_name
	f.OrgName, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// skip oc
	pos += 1

	// charset
	f.Charset = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	// column length
	f.ColumnLength = binary.LittleEndian.Uint32(p[pos:])
	pos += 4

	// type
	f.Type = p[pos]
	pos++

	// flag
	f.Flag = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	// decimals 1
	f.Decimal = p[pos]
	pos++

	// filter [0x00][0x00]
	pos += 2

	f.DefaultValue = nil
	// if more data, command was field list
	if len(p) > pos {
		// length of default value lenenc-int
		f.DefaultValueLength, _, n = LengthEncodedInt(p[pos:])
		pos += n

		if pos+int(f.DefaultValueLength) > len(p) {
			err = ErrMalformPacket
			return
		}

		// default value string[$len]
		f.DefaultValue = p[pos:(pos + int(f.DefaultValueLength))]
	}

	return
}

func (f *Field) Dump() []byte {
	if f.Data != nil {
		return f.Data
	}

	l := len(f.Schema) + len(f.Table) + len(f.OrgTable) + len(f.Name) + len(f.OrgName) + len(f.DefaultValue) + 48

	data := make([]byte, 0, l)

	data = append(data, PutLengthEncodedString([]byte("def"))...)

	data = append(data, PutLengthEncodedString(f.Schema)...)

	data = append(data, PutLengthEncodedString(f.Table)...)
	data = append(data, PutLengthEncodedString(f.OrgTable)...)

	data = append(data, PutLengthEncodedString(f.Name)...)
	data = append(data, PutLengthEncodedString(f.OrgName)...)

	data = append(data, 0x0c)

	data = append(data, Uint16ToBytes(f.Charset)...)
	data = append(data, Uint32ToBytes(f.ColumnLength)...)
	data = append(data, f.Type)
	data = append(data, Uint16ToBytes(f.Flag)...)
	data = append(data, f.Decimal)
	data = append(data, 0, 0)

	if f.DefaultValue != nil {
		data = append(data, Uint64ToBytes(f.DefaultValueLength)...)
		data = append(data, f.DefaultValue...)
	}

	return data
}

func stringToMysqlTime(s string) (TimeValue, error) {
	var v TimeValue

	timeFields := strings.SplitN(s, ":", 2)
	if len(timeFields) != 2 {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	hour, err := strconv.ParseInt(timeFields[0], 10, 64)
	if err != nil {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	if strings.HasPrefix(timeFields[0], "-") {
		v.IsNegative = true
		hour = Abs(hour)
	}

	day := int(hour / 24)
	hourRest := int(hour % 24)

	timeRest := strconv.Itoa(hourRest) + ":" + timeFields[1]

	ts, err := time.Parse("15:04:05", timeRest)
	if err != nil {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	if ts.Nanosecond()%1000 != 0 {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	v.Day = day
	v.Hour = ts.Hour()
	v.Minute = ts.Minute()
	v.Second = ts.Second()
	v.Microsecond = ts.Nanosecond() / 1000

	return v, nil
}

func mysqlTimeToBinaryResult(v TimeValue) []byte {
	var t []byte
	var length uint8

	if v.IsNull() {
		length = 0
		t = append(t, length)
	} else {
		if v.Microsecond == 0 {
			length = 8
		} else {
			length = 12
		}

		t = append(t, length)
		if v.IsNegative {
			t = append(t, 1)
		} else {
			t = append(t, 0)
		}

		t = AppendUint32(t, uint32(v.Day))
		t = append(t, uint8(v.Hour))
		t = append(t, uint8(v.Minute))
		t = append(t, uint8(v.Second))

		if v.Microsecond != 0 {
			t = AppendUint32(t, uint32(v.Microsecond))
		}
	}

	return t
}
