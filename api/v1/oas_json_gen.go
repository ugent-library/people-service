// Code generated by ogen, DO NOT EDIT.

package api

import (
	"math/bits"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/validate"
)

// Encode implements json.Marshaler.
func (s *Attribute) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *Attribute) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("scope")
		e.Str(s.Scope)
	}
	{
		e.FieldStart("key")
		e.Str(s.Key)
	}
	{
		e.FieldStart("value")
		e.Str(s.Value)
	}
}

var jsonFieldsNameOfAttribute = [3]string{
	0: "scope",
	1: "key",
	2: "value",
}

// Decode decodes Attribute from json.
func (s *Attribute) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode Attribute to nil")
	}
	var requiredBitSet [1]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "scope":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Str()
				s.Scope = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"scope\"")
			}
		case "key":
			requiredBitSet[0] |= 1 << 1
			if err := func() error {
				v, err := d.Str()
				s.Key = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"key\"")
			}
		case "value":
			requiredBitSet[0] |= 1 << 2
			if err := func() error {
				v, err := d.Str()
				s.Value = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"value\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode Attribute")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [1]uint8{
		0b00000111,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfAttribute) {
					name = jsonFieldsNameOfAttribute[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *Attribute) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *Attribute) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *Error) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *Error) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("code")
		e.Int64(s.Code)
	}
	{
		e.FieldStart("message")
		e.Str(s.Message)
	}
}

var jsonFieldsNameOfError = [2]string{
	0: "code",
	1: "message",
}

// Decode decodes Error from json.
func (s *Error) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode Error to nil")
	}
	var requiredBitSet [1]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "code":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Int64()
				s.Code = int64(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"code\"")
			}
		case "message":
			requiredBitSet[0] |= 1 << 1
			if err := func() error {
				v, err := d.Str()
				s.Message = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"message\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode Error")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [1]uint8{
		0b00000011,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfError) {
					name = jsonFieldsNameOfError[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *Error) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *Error) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *Identifier) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *Identifier) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("type")
		e.Str(s.Type)
	}
	{
		e.FieldStart("value")
		e.Str(s.Value)
	}
}

var jsonFieldsNameOfIdentifier = [2]string{
	0: "type",
	1: "value",
}

// Decode decodes Identifier from json.
func (s *Identifier) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode Identifier to nil")
	}
	var requiredBitSet [1]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "type":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Str()
				s.Type = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"type\"")
			}
		case "value":
			requiredBitSet[0] |= 1 << 1
			if err := func() error {
				v, err := d.Str()
				s.Value = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"value\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode Identifier")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [1]uint8{
		0b00000011,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfIdentifier) {
					name = jsonFieldsNameOfIdentifier[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *Identifier) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *Identifier) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode encodes bool as json.
func (o OptBool) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	e.Bool(bool(o.Value))
}

// Decode decodes bool from json.
func (o *OptBool) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptBool to nil")
	}
	o.Set = true
	v, err := d.Bool()
	if err != nil {
		return err
	}
	o.Value = bool(v)
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptBool) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptBool) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode encodes string as json.
func (o OptString) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	e.Str(string(o.Value))
}

// Decode decodes string from json.
func (o *OptString) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptString to nil")
	}
	o.Set = true
	v, err := d.Str()
	if err != nil {
		return err
	}
	o.Value = string(v)
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptString) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptString) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *Person) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *Person) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("name")
		e.Str(s.Name)
	}
	{
		if s.PreferredName.Set {
			e.FieldStart("preferredName")
			s.PreferredName.Encode(e)
		}
	}
	{
		if s.GivenName.Set {
			e.FieldStart("givenName")
			s.GivenName.Encode(e)
		}
	}
	{
		if s.PreferredGivenName.Set {
			e.FieldStart("preferredGivenName")
			s.PreferredGivenName.Encode(e)
		}
	}
	{
		if s.FamilyName.Set {
			e.FieldStart("familyName")
			s.FamilyName.Encode(e)
		}
	}
	{
		if s.PreferredFamilyName.Set {
			e.FieldStart("preferredFamilyName")
			s.PreferredFamilyName.Encode(e)
		}
	}
	{
		if s.HonorificPrefix.Set {
			e.FieldStart("honorificPrefix")
			s.HonorificPrefix.Encode(e)
		}
	}
	{
		if s.Email.Set {
			e.FieldStart("email")
			s.Email.Encode(e)
		}
	}
	{
		if s.Username.Set {
			e.FieldStart("username")
			s.Username.Encode(e)
		}
	}
	{
		if s.Active.Set {
			e.FieldStart("active")
			s.Active.Encode(e)
		}
	}
	{
		if s.Attributes != nil {
			e.FieldStart("attributes")
			e.ArrStart()
			for _, elem := range s.Attributes {
				elem.Encode(e)
			}
			e.ArrEnd()
		}
	}
	{
		e.FieldStart("identifiers")
		e.ArrStart()
		for _, elem := range s.Identifiers {
			elem.Encode(e)
		}
		e.ArrEnd()
	}
}

var jsonFieldsNameOfPerson = [12]string{
	0:  "name",
	1:  "preferredName",
	2:  "givenName",
	3:  "preferredGivenName",
	4:  "familyName",
	5:  "preferredFamilyName",
	6:  "honorificPrefix",
	7:  "email",
	8:  "username",
	9:  "active",
	10: "attributes",
	11: "identifiers",
}

// Decode decodes Person from json.
func (s *Person) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode Person to nil")
	}
	var requiredBitSet [2]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "name":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Str()
				s.Name = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"name\"")
			}
		case "preferredName":
			if err := func() error {
				s.PreferredName.Reset()
				if err := s.PreferredName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredName\"")
			}
		case "givenName":
			if err := func() error {
				s.GivenName.Reset()
				if err := s.GivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"givenName\"")
			}
		case "preferredGivenName":
			if err := func() error {
				s.PreferredGivenName.Reset()
				if err := s.PreferredGivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredGivenName\"")
			}
		case "familyName":
			if err := func() error {
				s.FamilyName.Reset()
				if err := s.FamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"familyName\"")
			}
		case "preferredFamilyName":
			if err := func() error {
				s.PreferredFamilyName.Reset()
				if err := s.PreferredFamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredFamilyName\"")
			}
		case "honorificPrefix":
			if err := func() error {
				s.HonorificPrefix.Reset()
				if err := s.HonorificPrefix.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"honorificPrefix\"")
			}
		case "email":
			if err := func() error {
				s.Email.Reset()
				if err := s.Email.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"email\"")
			}
		case "username":
			if err := func() error {
				s.Username.Reset()
				if err := s.Username.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"username\"")
			}
		case "active":
			if err := func() error {
				s.Active.Reset()
				if err := s.Active.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"active\"")
			}
		case "attributes":
			if err := func() error {
				s.Attributes = make([]Attribute, 0)
				if err := d.Arr(func(d *jx.Decoder) error {
					var elem Attribute
					if err := elem.Decode(d); err != nil {
						return err
					}
					s.Attributes = append(s.Attributes, elem)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"attributes\"")
			}
		case "identifiers":
			requiredBitSet[1] |= 1 << 3
			if err := func() error {
				s.Identifiers = make([]Identifier, 0)
				if err := d.Arr(func(d *jx.Decoder) error {
					var elem Identifier
					if err := elem.Decode(d); err != nil {
						return err
					}
					s.Identifiers = append(s.Identifiers, elem)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"identifiers\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode Person")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [2]uint8{
		0b00000001,
		0b00001000,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfPerson) {
					name = jsonFieldsNameOfPerson[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *Person) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *Person) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *PersonRecord) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *PersonRecord) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("name")
		e.Str(s.Name)
	}
	{
		if s.PreferredName.Set {
			e.FieldStart("preferredName")
			s.PreferredName.Encode(e)
		}
	}
	{
		if s.GivenName.Set {
			e.FieldStart("givenName")
			s.GivenName.Encode(e)
		}
	}
	{
		if s.PreferredGivenName.Set {
			e.FieldStart("preferredGivenName")
			s.PreferredGivenName.Encode(e)
		}
	}
	{
		if s.FamilyName.Set {
			e.FieldStart("familyName")
			s.FamilyName.Encode(e)
		}
	}
	{
		if s.PreferredFamilyName.Set {
			e.FieldStart("preferredFamilyName")
			s.PreferredFamilyName.Encode(e)
		}
	}
	{
		if s.HonorificPrefix.Set {
			e.FieldStart("honorificPrefix")
			s.HonorificPrefix.Encode(e)
		}
	}
	{
		if s.Email.Set {
			e.FieldStart("email")
			s.Email.Encode(e)
		}
	}
	{
		if s.Username.Set {
			e.FieldStart("username")
			s.Username.Encode(e)
		}
	}
	{
		e.FieldStart("active")
		e.Bool(s.Active)
	}
	{
		if s.Attributes != nil {
			e.FieldStart("attributes")
			e.ArrStart()
			for _, elem := range s.Attributes {
				elem.Encode(e)
			}
			e.ArrEnd()
		}
	}
	{
		e.FieldStart("identifiers")
		e.ArrStart()
		for _, elem := range s.Identifiers {
			elem.Encode(e)
		}
		e.ArrEnd()
	}
	{
		e.FieldStart("createdAt")
		json.EncodeDateTime(e, s.CreatedAt)
	}
	{
		e.FieldStart("updatedAt")
		json.EncodeDateTime(e, s.UpdatedAt)
	}
}

var jsonFieldsNameOfPersonRecord = [14]string{
	0:  "name",
	1:  "preferredName",
	2:  "givenName",
	3:  "preferredGivenName",
	4:  "familyName",
	5:  "preferredFamilyName",
	6:  "honorificPrefix",
	7:  "email",
	8:  "username",
	9:  "active",
	10: "attributes",
	11: "identifiers",
	12: "createdAt",
	13: "updatedAt",
}

// Decode decodes PersonRecord from json.
func (s *PersonRecord) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode PersonRecord to nil")
	}
	var requiredBitSet [2]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "name":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Str()
				s.Name = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"name\"")
			}
		case "preferredName":
			if err := func() error {
				s.PreferredName.Reset()
				if err := s.PreferredName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredName\"")
			}
		case "givenName":
			if err := func() error {
				s.GivenName.Reset()
				if err := s.GivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"givenName\"")
			}
		case "preferredGivenName":
			if err := func() error {
				s.PreferredGivenName.Reset()
				if err := s.PreferredGivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredGivenName\"")
			}
		case "familyName":
			if err := func() error {
				s.FamilyName.Reset()
				if err := s.FamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"familyName\"")
			}
		case "preferredFamilyName":
			if err := func() error {
				s.PreferredFamilyName.Reset()
				if err := s.PreferredFamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferredFamilyName\"")
			}
		case "honorificPrefix":
			if err := func() error {
				s.HonorificPrefix.Reset()
				if err := s.HonorificPrefix.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"honorificPrefix\"")
			}
		case "email":
			if err := func() error {
				s.Email.Reset()
				if err := s.Email.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"email\"")
			}
		case "username":
			if err := func() error {
				s.Username.Reset()
				if err := s.Username.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"username\"")
			}
		case "active":
			requiredBitSet[1] |= 1 << 1
			if err := func() error {
				v, err := d.Bool()
				s.Active = bool(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"active\"")
			}
		case "attributes":
			if err := func() error {
				s.Attributes = make([]Attribute, 0)
				if err := d.Arr(func(d *jx.Decoder) error {
					var elem Attribute
					if err := elem.Decode(d); err != nil {
						return err
					}
					s.Attributes = append(s.Attributes, elem)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"attributes\"")
			}
		case "identifiers":
			requiredBitSet[1] |= 1 << 3
			if err := func() error {
				s.Identifiers = make([]Identifier, 0)
				if err := d.Arr(func(d *jx.Decoder) error {
					var elem Identifier
					if err := elem.Decode(d); err != nil {
						return err
					}
					s.Identifiers = append(s.Identifiers, elem)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"identifiers\"")
			}
		case "createdAt":
			requiredBitSet[1] |= 1 << 4
			if err := func() error {
				v, err := json.DecodeDateTime(d)
				s.CreatedAt = v
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"createdAt\"")
			}
		case "updatedAt":
			requiredBitSet[1] |= 1 << 5
			if err := func() error {
				v, err := json.DecodeDateTime(d)
				s.UpdatedAt = v
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"updatedAt\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode PersonRecord")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [2]uint8{
		0b00000001,
		0b00111010,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfPersonRecord) {
					name = jsonFieldsNameOfPersonRecord[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *PersonRecord) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *PersonRecord) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}
