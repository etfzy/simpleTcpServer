package proto

import (
	"bytes"
	"encoding/binary"
	"errors"

	mem "github.com/etfzy/mempool/base"
)

type Tlv struct {
	flag      uint64
	flagLen   uint64
	lengthLen uint64
	maxLength uint64
	border    binary.ByteOrder
}

func NewTlv() *Tlv {
	return &Tlv{}
}

func (r *Tlv) write(length uint64, value uint64, out *mem.Buffer[byte]) error {
	buffer := bytes.NewBuffer(*out.Buf())
	switch length {
	case 2:
		return binary.Write(buffer, r.border, uint16(value))
	case 4:
		return binary.Write(buffer, r.border, uint32(value))
	case 8:
		return binary.Write(buffer, r.border, uint64(value))
	}

	return errors.New("length is unfair!")
}

func (r *Tlv) read(length uint64, input *mem.Buffer[byte]) uint64 {
	switch length {
	case 2:
		return uint64(r.border.Uint16(*input.Buf()))
	case 4:
		return uint64(r.border.Uint32(*input.Buf()))
	case 8:
		return uint64(r.border.Uint64(*input.Buf()))
	}

	return 0
}

func (r *Tlv) WriteFlag(out *mem.Buffer[byte]) error {
	return r.write(r.flagLen, r.flag, out)
}

func (r *Tlv) WriteLength(value uint64, out *mem.Buffer[byte]) error {
	return r.write(r.lengthLen, value, out)
}

func (r *Tlv) ReadFlag(bs *mem.Buffer[byte]) error {
	flag := r.read(r.flagLen, bs)
	if flag != r.flag {
		return errors.New("receive flag is not match!")
	}

	return nil
}

func (r *Tlv) ReadLength(bs *mem.Buffer[byte]) (uint64, error) {
	length := r.read(r.lengthLen, bs)
	if length <= r.maxLength && length != 0 {
		return length, nil
	}

	return length, errors.New("content length is too long!")
}

func (r *Tlv) GetFlag() uint64 {
	return r.flag
}

func (r *Tlv) SetFlag(val uint64) *Tlv {
	r.flag = val
	return r
}

func (r *Tlv) GetFlagLen() uint64 {
	return r.flagLen
}

func (r *Tlv) SetFlagLen(val uint64) *Tlv {
	r.flagLen = val
	return r
}

func (r *Tlv) GetLengthLen() uint64 {
	return r.lengthLen
}

func (r *Tlv) SetLengthLen(val uint64) *Tlv {
	r.lengthLen = val
	return r
}

func (r *Tlv) GetMaxLength() uint64 {
	return r.maxLength
}

func (r *Tlv) SetMaxLength(val uint64) *Tlv {
	r.maxLength = val
	return r
}

func (r *Tlv) GetBorder() binary.ByteOrder {
	return r.border
}

func (r *Tlv) SetBorder(val binary.ByteOrder) *Tlv {
	r.border = val
	return r
}

func (r *Tlv) CheckTlv() error {
	if r.GetFlag() == 0 {
		return errors.New("proto config flag can not be zero!")
	}

	if r.GetFlagLen() == 0 || (r.GetFlagLen() != 2 && r.GetFlagLen() != 4 && r.GetFlagLen() != 8) {
		return errors.New("proto config flag length must be 2|4|8!")
	}

	if r.GetLengthLen() == 0 || (r.GetLengthLen() != 2 && r.GetLengthLen() != 4 && r.GetLengthLen() != 8) {
		return errors.New("proto config flag length must be 2|4|8!")
	}

	if r.GetMaxLength() == 0 {
		return errors.New("proto config max length can not be zero!")
	}

	if r.GetBorder() == nil {
		return errors.New("proto config border can not be nil!")
	}
	return nil
}
