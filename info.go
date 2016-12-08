package bcf

import (
	"encoding/binary"
	"fmt"
	"math"
)

type typed uint8

const (
	sFlag    typed = 0
	sInt8    typed = 1
	sInt16   typed = 2
	sInt32   typed = 3
	sFloat32 typed = 5
	sChar    typed = 7
	sArray   typed = 15
)

func (s typed) nBytes() int {
	switch s {
	case sFlag:
		return 0
	case sInt8:
		return 1
	case sInt16:
		return 2
	case sInt32:
		return 4
	case sFloat32:
		return 4
	case sChar:
		return 1
	}
	panic(fmt.Sprintf("unknown typed: %d", s))
}

func (t typed) String() string {
	switch t {
	case sFlag:
		return "flag"
	case sInt8:
		return "int8"
	case sInt16:
		return "int16"
	case sInt32:
		return "int32"
	case sFloat32:
		return "float32"
	case sChar:
		return "char"
	case sArray:
		return "array"
	}
	return "uknown"
}

type infoentry struct {
	count     uint32
	headerKey uint16
	etype     typed
	buf       []byte
}

func (i infoentry) val() interface{} {
	if i.etype == sChar {
		return string(i.buf)
	}
	if i.count == 1 {
		switch i.etype {
		case sInt8:
			return int(i.buf[0])
		case sInt16:
			return int(binary.LittleEndian.Uint16(i.buf))
		case sInt32:
			return int(binary.LittleEndian.Uint32(i.buf))
		case sFloat32:
			return math.Float32frombits(binary.LittleEndian.Uint32(i.buf))
		default:
			panic("bcf: unknown type")
		}
	}
	// flag
	if i.count == 0 {
		return true
	}

	vals := make([]interface{}, i.count)
	for k := range vals {
		switch i.etype {
		case sInt32:
			vals[k] = int(binary.LittleEndian.Uint32(i.buf[k*4 : k*4+4]))
		case sFloat32:
			vals[k] = math.Float32frombits(binary.LittleEndian.Uint32(i.buf[k*4 : k*4+4]))
		case sInt16:
			vals[k] = int(binary.LittleEndian.Uint16(i.buf[k*2 : k*2+2]))
		case sInt8:
			vals[k] = int(i.buf[k])
		default:
			panic("bcf: unknown type")
		}
	}
	return vals
}

// Info contains the entries from the INFO field for a given Record
// TODO: add a link to the header so we can search by key.
type Info struct {
	buffer
	entries []infoentry
}

func infoFromBytes(buf []byte) Info {
	info := Info{buffer: buffer{buf: buf}, entries: make([]infoentry, 0, 8)}

	for info.off < uint32(len(info.buf)) {
		// key is the index of this entry in the header dict.
		ktype, _ := info.typed()
		key := info.int(ktype.nBytes())

		itype, isize := info.typed()
		if isize == 15 {
			atype, _ := info.typed()
			// update isize to be the length of the array.
			isize = info.int(atype.nBytes())

		}
		off := info.off
		info.off += uint32(itype.nBytes() * isize)
		e := infoentry{count: uint32(isize), headerKey: uint16(key), etype: itype, buf: info.buf[off:info.off]}
		//fmt.Println("val, type, count:", e.val(), e.etype, e.count)
		info.entries = append(info.entries, e)
	}
	return info
}

/*
In BCF2, a typed value consists of a typing byte and the actual value with type mandated by the typing
byte. In the typing byte, the lowest four bits give the atomic type. If the number represented by the higher
4 bits is smaller than 15, it is the size of the following vector; if the number equals 15, the following typed
integer is the array size. The highest 4 bits of a Flag type equals 0 and in this case, n
*/

// buffer allows us to share code between info and format fields
type buffer struct {
	buf []byte
	off uint32
}

func (b *buffer) typed() (typed, int) {
	t, s := typed(b.buf[b.off]&0xf), int(b.buf[b.off]>>4)
	b.off++
	return t, s
}

func (b *buffer) int(nbytes int) int {
	off := b.off
	b.off += uint32(nbytes)
	switch nbytes {
	case 1:
		return int(b.buf[off])
	case 2:
		return int(binary.LittleEndian.Uint16(b.buf[off:b.off]))
	case 4:
		return int(binary.LittleEndian.Uint32(b.buf[off:b.off]))
	default:
		panic(fmt.Sprintf("bcf: unknown int size: %d", nbytes))
	}
}
