// Code generated by ogen, DO NOT EDIT.

package api

import (
	"math/bits"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/validate"
)

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
		e.FieldStart("active")
		e.Bool(s.Active)
	}
	{
		e.FieldStart("name")
		e.Str(s.Name)
	}
	{
		if s.PreferredName.Set {
			e.FieldStart("preferred_name")
			s.PreferredName.Encode(e)
		}
	}
	{
		if s.GivenName.Set {
			e.FieldStart("given_name")
			s.GivenName.Encode(e)
		}
	}
	{
		if s.FamilyName.Set {
			e.FieldStart("family_name")
			s.FamilyName.Encode(e)
		}
	}
	{
		if s.PreferredGivenName.Set {
			e.FieldStart("preferred_given_name")
			s.PreferredGivenName.Encode(e)
		}
	}
	{
		if s.PreferredFamilyName.Set {
			e.FieldStart("preferred_family_name")
			s.PreferredFamilyName.Encode(e)
		}
	}
	{
		if s.HonorificPrefix.Set {
			e.FieldStart("honorific_prefix")
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
		e.FieldStart("identifiers")
		s.Identifiers.Encode(e)
	}
}

var jsonFieldsNameOfPerson = [10]string{
	0: "active",
	1: "name",
	2: "preferred_name",
	3: "given_name",
	4: "family_name",
	5: "preferred_given_name",
	6: "preferred_family_name",
	7: "honorific_prefix",
	8: "email",
	9: "identifiers",
}

// Decode decodes Person from json.
func (s *Person) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode Person to nil")
	}
	var requiredBitSet [2]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "active":
			requiredBitSet[0] |= 1 << 0
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
		case "name":
			requiredBitSet[0] |= 1 << 1
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
		case "preferred_name":
			if err := func() error {
				s.PreferredName.Reset()
				if err := s.PreferredName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferred_name\"")
			}
		case "given_name":
			if err := func() error {
				s.GivenName.Reset()
				if err := s.GivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"given_name\"")
			}
		case "family_name":
			if err := func() error {
				s.FamilyName.Reset()
				if err := s.FamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"family_name\"")
			}
		case "preferred_given_name":
			if err := func() error {
				s.PreferredGivenName.Reset()
				if err := s.PreferredGivenName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferred_given_name\"")
			}
		case "preferred_family_name":
			if err := func() error {
				s.PreferredFamilyName.Reset()
				if err := s.PreferredFamilyName.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"preferred_family_name\"")
			}
		case "honorific_prefix":
			if err := func() error {
				s.HonorificPrefix.Reset()
				if err := s.HonorificPrefix.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"honorific_prefix\"")
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
		case "identifiers":
			requiredBitSet[1] |= 1 << 1
			if err := func() error {
				if err := s.Identifiers.Decode(d); err != nil {
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
		0b00000011,
		0b00000010,
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
func (s PersonIdentifiers) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields implements json.Marshaler.
func (s PersonIdentifiers) encodeFields(e *jx.Encoder) {
	for k, elem := range s {
		e.FieldStart(k)

		e.ArrStart()
		for _, elem := range elem {
			e.Str(elem)
		}
		e.ArrEnd()
	}
}

// Decode decodes PersonIdentifiers from json.
func (s *PersonIdentifiers) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode PersonIdentifiers to nil")
	}
	m := s.init()
	var propertiesCount int
	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		propertiesCount++
		var elem []string
		if err := func() error {
			elem = make([]string, 0)
			if err := d.Arr(func(d *jx.Decoder) error {
				var elemElem string
				v, err := d.Str()
				elemElem = string(v)
				if err != nil {
					return err
				}
				elem = append(elem, elemElem)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return errors.Wrapf(err, "decode field %q", k)
		}
		m[string(k)] = elem
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode PersonIdentifiers")
	}
	// Validate properties count.
	if err := (validate.Object{
		MinProperties:    1,
		MinPropertiesSet: true,
		MaxProperties:    0,
		MaxPropertiesSet: false,
	}).ValidateProperties(propertiesCount); err != nil {
		return errors.Wrap(err, "object")
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s PersonIdentifiers) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *PersonIdentifiers) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}
