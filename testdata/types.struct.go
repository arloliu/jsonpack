package testdata

type NativeTypes struct {
	FieldBool       bool     `json:"field_bool"`
	FieldBoolPtr    *bool    `json:"field_bool_ptr"`
	FieldInt8       int8     `json:"field_int8"`
	FieldInt8Ptr    *int8    `json:"field_int8_ptr"`
	FieldInt16      int16    `json:"field_int16"`
	FieldInt16Ptr   *int16   `json:"field_int16_ptr"`
	FieldInt32      int32    `json:"field_int32"`
	FieldInt32Ptr   *int32   `json:"field_int32_ptr"`
	FieldInt64      int64    `json:"field_int64"`
	FieldInt64Ptr   *int64   `json:"field_int64_ptr"`
	FieldUint8      uint8    `json:"field_uint8"`
	FieldUint8Ptr   *uint8   `json:"field_uint8_ptr"`
	FieldUint16     uint16   `json:"field_uint16"`
	FieldUint16Ptr  *uint16  `json:"field_uint16_ptr"`
	FieldUint32     uint32   `json:"field_uint32"`
	FieldUint32Ptr  *uint32  `json:"field_uint32_ptr"`
	FieldUint64     uint64   `json:"field_uint64"`
	FieldUint64Ptr  *uint64  `json:"field_uint64_ptr"`
	FieldFloat32    float32  `json:"field_float32"`
	FieldFloat32Ptr *float32 `json:"field_float32_ptr"`
	FieldFloat64    float64  `json:"field_float64"`
	FieldFloat64Ptr *float64 `json:"field_float64_ptr"`
	FieldString     string   `json:"field_string"`
	FieldStringPtr  *string  `json:"field_string_ptr"`
	OmitField       string   `json:"-"`
}

type Types struct {
	NativeTypes
	SliceTypes []NativeTypes `json:"slice_types"`
}
