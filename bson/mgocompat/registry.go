// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mgocompat

import (
	"errors"
	"reflect"
	"time"

	"github.com/wimspaargaren/mongo-go-driver/bson"
	"github.com/wimspaargaren/mongo-go-driver/bson/bsoncodec"
	"github.com/wimspaargaren/mongo-go-driver/bson/bsonoptions"
	"github.com/wimspaargaren/mongo-go-driver/bson/bsontype"
)

var (
	// ErrSetZero may be returned from a SetBSON method to have the value set to its respective zero value.
	ErrSetZero = errors.New("set to zero")

	tInt            = reflect.TypeOf(int(0))
	tUint           = reflect.TypeOf(uint(0))
	tTime           = reflect.TypeOf(time.Time{})
	tM              = reflect.TypeOf(bson.M{})
	tInterfaceSlice = reflect.TypeOf([]interface{}{})
	tByteSlice      = reflect.TypeOf([]byte{})
	tEmpty          = reflect.TypeOf((*interface{})(nil)).Elem()
	tGetter         = reflect.TypeOf((*Getter)(nil)).Elem()
	tSetter         = reflect.TypeOf((*Setter)(nil)).Elem()
)

// mgoRegistry is the default bsoncodec.Registry. It contains the default codecs and the
// primitive codecs.
var mgoRegistry = newRegistryBuilder().Build()

// newRegistryBuilder creates a new RegistryBuilder configured with the default encoders and
// deocders from the bsoncodec.DefaultValueEncoders and bsoncodec.DefaultValueDecoders types and the
// PrimitiveCodecs type in this package.
func newRegistryBuilder() *bsoncodec.RegistryBuilder {
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.PrimitiveCodecs{}.RegisterPrimitiveCodecs(rb)

	structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.DefaultStructTagParser,
		bsonoptions.StructCodec().
			SetDecodeZeroStruct(true).
			SetDecodeDeepZeroInline(true).
			SetEncodeOmitDefaultStruct(true).
			SetAllowUnexportedFields(true))
	emptyInterCodec := bsoncodec.NewEmptyInterfaceCodec(
		bsonoptions.EmptyInterfaceCodec().
			SetDecodeAsMap(true).
			SetDecodeBinaryAsSlice(true))
	mapCodec := bsoncodec.NewMapCodec(
		bsonoptions.MapCodec().
			SetDecodeZerosMap(true).
			SetEncodeNilAsEmpty(true))

	rb.RegisterDecoder(tEmpty, emptyInterCodec).
		RegisterDecoder(tSetter, bsoncodec.ValueDecoderFunc(SetterDecodeValue)).
		RegisterDefaultDecoder(reflect.String, bsoncodec.NewStringCodec(bsonoptions.StringCodec().SetDecodeObjectIDAsHex(false))).
		RegisterDefaultDecoder(reflect.Struct, structcodec).
		RegisterDefaultDecoder(reflect.Map, mapCodec).
		RegisterEncoder(tByteSlice, bsoncodec.NewByteSliceCodec(bsonoptions.ByteSliceCodec().SetEncodeNilAsEmpty(true))).
		RegisterEncoder(tGetter, bsoncodec.ValueEncoderFunc(GetterEncodeValue)).
		RegisterDefaultEncoder(reflect.Struct, structcodec).
		RegisterDefaultEncoder(reflect.Slice, bsoncodec.NewSliceCodec(bsonoptions.SliceCodec().SetEncodeNilAsEmpty(true))).
		RegisterDefaultEncoder(reflect.Map, mapCodec).
		RegisterTypeMapEntry(bsontype.Int32, tInt).
		RegisterTypeMapEntry(bsontype.Type(0), tM).
		RegisterTypeMapEntry(bsontype.DateTime, tTime).
		RegisterTypeMapEntry(bsontype.Array, tInterfaceSlice)

	return rb
}
