package bcf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/biogo/hts/bgzf"
)

var bcfMagic = []byte{'B', 'C', 'F'}

// BCF contains the methods for reading a BCF
type BCF struct {
	Header Header
	bgz    *bgzf.Reader
	buf    []byte
}

// NewReader returns a BCF reader with the number of decompression threads.
func NewReader(r io.Reader, rd int) (*BCF, error) {
	bg, err := bgzf.NewReader(r, rd)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 4)
	if _, err := bg.Read(buf[:3]); err != nil {
		log.Println("initial read")
		return nil, err
	}
	if !bytes.Equal(buf[:3], bcfMagic) {
		return nil, fmt.Errorf("bcf: incorrect header for bcf")
	}
	h := Header{}
	if err := binary.Read(bg, binary.LittleEndian, &h.version); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(bg, buf); err != nil {
		return nil, err
	}
	h.lText = binary.LittleEndian.Uint32(buf)
	h.Text = make([]byte, int(h.lText))
	if _, err := io.ReadFull(bg, h.Text); err != nil {
		return nil, err
	}
	err = h.parse()
	return &BCF{Header: h, bgz: bg, buf: buf}, err
}

// Read returns the next variant.
func (b *BCF) Read() (*Variant, error) {

	m := make([]uint32, 8)
	if err := binary.Read(b.bgz, binary.LittleEndian, &m); err != nil {
		return nil, err
	}

	var v Variant
	v.Chrom = m[2]
	v.Pos = m[3]

	v.Qual = math.Float32frombits(m[5])
	v.nallele = uint16(m[6] >> 16)
	v.ninfo = uint16(m[6] & 0xffff)
	v.nfmt = uint8(m[7] >> 24)
	v.nsample = uint32(m[7] & 0xffffff)
	v.header = &b.Header

	// account for the 6 additional 32bit ints in m
	v.shared = make([]byte, m[0]-24)
	v.indiv = make([]byte, m[1])

	if _, err := io.ReadFull(b.bgz, v.shared); err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(b.bgz, v.indiv); err != nil {
		return nil, err
	}
	v.Id = v.readBytes()
	v.Alleles = make([][]byte, v.nallele)
	for i := uint16(0); i < v.nallele; i++ {
		v.Alleles[i] = v.readBytes()
	}
	// TODO: fix whatever is wrong here.
	v.sharedOff++
	v.filters = v.readBytes()
	v.sharedOff++

	return &v, nil
}

func (v *Variant) readBytes() []byte {
	vtype, vsize := (v.shared[v.sharedOff] & 0xf), (v.shared[v.sharedOff] >> 4)
	v.sharedOff += 1
	if vtype == 0 {
		return nil
	}
	if vtype != 7 && vtype != 1 {
		panic(fmt.Sprintf("bcf: non-char type for bytes: %d", vtype))
	}
	val := v.shared[v.sharedOff : v.sharedOff+int(vsize)]
	if vsize > 0 && val[len(val)-1] == 0 {
		val = val[:len(val)-1]
	}
	v.sharedOff += int(vsize)
	return val
}

// Variant represents a single variant object from the BCF.
type Variant struct {
	Chrom   uint32 // CHROM
	Pos     uint32
	Alleles [][]byte
	// these are the index in the header of the appropriate filter.
	filters []uint8
	Id      []byte
	Qual    float32
	// data thru info
	shared    []byte
	sharedOff int
	// data after info
	indiv   []byte
	nallele uint16
	nsample uint32
	ninfo   uint16
	nfmt    uint8
	info    *Info
	header  *Header
}

// Info returns the struct containing the INFO fields.
func (v *Variant) Info() *Info {
	if v.info == nil {
		info := infoFromBytes(v.shared)
		v.info = &info
	}
	return v.info
}
