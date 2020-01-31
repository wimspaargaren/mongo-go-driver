// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"fmt"
	"reflect"

	"github.com/wimspaargaren/mongo-go-driver/bson/bsonoptions"
	"github.com/wimspaargaren/mongo-go-driver/bson/bsonrw"
	"github.com/wimspaargaren/mongo-go-driver/bson/bsontype"
)

var defaultByteSliceCodec = NewByteSliceCodec()

// ByteSliceCodec is the Codec used for []byte values.
type ByteSliceCodec struct {
	EncodeNilAsEmpty bool
}

var _ ValueCodec = &ByteSliceCodec{}

// NewByteSliceCodec returns a StringCodec with options opts.
func NewByteSliceCodec(opts ...*bsonoptions.ByteSliceCodecOptions) *ByteSliceCodec {
	byteSliceOpt := bsonoptions.MergeByteSliceCodecOptions(opts...)
	codec := ByteSliceCodec{}
	if byteSliceOpt.EncodeNilAsEmpty != nil {
		codec.EncodeNilAsEmpty = *byteSliceOpt.EncodeNilAsEmpty
	}
	return &codec
}

// EncodeValue is the ValueEncoder for []byte.
func (bsc *ByteSliceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tByteSlice {
		return ValueEncoderError{Name: "ByteSliceEncodeValue", Types: []reflect.Type{tByteSlice}, Received: val}
	}
	if val.IsNil() && !bsc.EncodeNilAsEmpty {
		return vw.WriteNull()
	}
	return vw.WriteBinary(val.Interface().([]byte))
}

// DecodeValue is the ValueDecoder for []byte.
func (bsc *ByteSliceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tByteSlice {
		return ValueDecoderError{Name: "ByteSliceDecodeValue", Types: []reflect.Type{tByteSlice}, Received: val}
	}

	var data []byte
	var err error
	switch vr.Type() {
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		data = []byte(str)
	case bsontype.Symbol:
		sym, err := vr.ReadSymbol()
		if err != nil {
			return err
		}
		data = []byte(sym)
	case bsontype.Binary:
		var subtype byte
		data, subtype, err = vr.ReadBinary()
		if err != nil {
			return err
		}
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld {
			return fmt.Errorf("ByteSliceDecodeValue can only be used to decode subtype 0x00 or 0x02 for %s, got %v", bsontype.Binary, subtype)
		}
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a []byte", vr.Type())
	}

	val.Set(reflect.ValueOf(data))
	return nil
}