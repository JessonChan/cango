// Copyright 2020 Cango Author.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cango

import (
	"reflect"
	"strconv"
	"time"
)

type Caster func(string) reflect.Value

var (
	boolType     = reflect.Bool
	float32Type  = reflect.Float32
	float64Type  = reflect.Float64
	intType      = reflect.Int
	int8Type     = reflect.Int8
	int16Type    = reflect.Int16
	int32Type    = reflect.Int32
	int64Type    = reflect.Int64
	stringType   = reflect.String
	uintType     = reflect.Uint
	uint8Type    = reflect.Uint8
	uint16Type   = reflect.Uint16
	uint32Type   = reflect.Uint32
	uint64Type   = reflect.Uint64
	timeTypeKind = reflect.Kind(1e7)
)

var casterMap = map[reflect.Kind]Caster{
	boolType:     castBool,
	float32Type:  castFloat32,
	float64Type:  castFloat64,
	intType:      castInt,
	int8Type:     castInt8,
	int16Type:    castInt16,
	int32Type:    castInt32,
	int64Type:    castInt64,
	uintType:     castUint,
	uint8Type:    castUint8,
	uint16Type:   castUint16,
	uint32Type:   castUint32,
	uint64Type:   castUint64,
	stringType:   castString,
	timeTypeKind: castTime,
}

func castBool(value string) reflect.Value {
	if value == "on" || value == "1" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.ValueOf(false)
}

func castFloat32(value string) reflect.Value {
	var f32 float32
	if v, err := strconv.ParseFloat(value, 32); err == nil {
		return reflect.ValueOf(float32(v))
	}
	return reflect.ValueOf(f32)
}

func castFloat64(value string) reflect.Value {
	var f64 float64
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.ValueOf(f64)
}

type integerType interface {
	int | int8 | int16 | int32 | int64
}

// cast str to integer in generic method
func castInteger[T integerType](str string, t T) reflect.Value {
	bitSize := 90
	switch any(t).(type) {
	case int8:
		bitSize = 8
	case int16:
		bitSize = 16
	case int32:
		bitSize = 32
	case int64:
		bitSize = 64
	}

	if v, err := strconv.ParseInt(str, 10, bitSize); err == nil {
		return reflect.ValueOf(T(v))
	}
	return reflect.ValueOf(t)
}

type uIntegerType interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func castUInteger[T integerType](str string, t T) reflect.Value {
	if v, err := strconv.ParseUint(str, 10, 64); err == nil {
		return reflect.ValueOf(T(v))
	}
	return reflect.ValueOf(t)
}

func castInt(value string) reflect.Value {
	var i int
	if v, err := strconv.ParseInt(value, 10, 0); err == nil {
		return reflect.ValueOf(int(v))
	}
	return reflect.ValueOf(i)
}

func castInt8(value string) reflect.Value {
	var i8 int8
	if v, err := strconv.ParseInt(value, 10, 8); err == nil {
		return reflect.ValueOf(int8(v))
	}
	return reflect.ValueOf(i8)
}

func castInt16(value string) reflect.Value {
	var i16 int16
	if v, err := strconv.ParseInt(value, 10, 16); err == nil {
		return reflect.ValueOf(int16(v))
	}
	return reflect.ValueOf(i16)
}

func castInt32(value string) reflect.Value {
	var i32 int32
	if v, err := strconv.ParseInt(value, 10, 32); err == nil {
		return reflect.ValueOf(int32(v))
	}
	return reflect.ValueOf(i32)
}

func castInt64(value string) reflect.Value {
	var i64 int64
	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.ValueOf(i64)
}

func castString(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func castUint(value string) reflect.Value {
	var u uint
	if v, err := strconv.ParseUint(value, 10, 0); err == nil {
		return reflect.ValueOf(uint(v))
	}
	return reflect.ValueOf(u)
}

func castUint8(value string) reflect.Value {
	var u8 uint8
	if v, err := strconv.ParseUint(value, 10, 8); err == nil {
		return reflect.ValueOf(uint8(v))
	}
	return reflect.ValueOf(u8)
}

func castUint16(value string) reflect.Value {
	var u16 uint16
	if v, err := strconv.ParseUint(value, 10, 16); err == nil {
		return reflect.ValueOf(uint16(v))
	}
	return reflect.ValueOf(u16)
}

func castUint32(value string) reflect.Value {
	var u32 uint32
	if v, err := strconv.ParseUint(value, 10, 32); err == nil {
		return reflect.ValueOf(uint32(v))
	}
	return reflect.ValueOf(u32)
}

func castUint64(value string) reflect.Value {
	var u64 uint64
	if v, err := strconv.ParseUint(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.ValueOf(u64)
}

const (
	shortSimpleTimeFormat = "2006-01-02"
	longSimpleTimeFormat  = "2006-01-02 15:04:05"
)

func castTime(value string) reflect.Value {
	var layout string
	if len(value) == 10 {
		layout = shortSimpleTimeFormat
	}
	if len(value) == 19 {
		layout = longSimpleTimeFormat
	}
	timeTime, _ := time.ParseInLocation(layout, value, time.Local)
	return reflect.ValueOf(timeTime)
}
