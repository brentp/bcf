package bcf

var bcfTypeShift = []uint8{0, 0, 1, 2, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Format contains the BCF format/sample entries.
type Format struct {
	buffer
	entries []infoentry
}

func formatFromBytes(buf []byte, nSample uint32) Format {
	f := Format{buffer: buffer{buf: buf}, entries: make([]infoentry, 0, 3)}
	for f.off < uint32(len(f.buffer.buf)) {
		// key is the index of this entry in the header dict.
		ktype, _ := f.typed()
		key := f.int(ktype.nBytes())
		itype, isize := f.typed()

		/*
			type: 1, id: 4, n: 2, size: 2
			type: 5, id: 6, n: 3, size: 12
			type: 2, id: 1, n: 2, size: 4
			type: 2, id: 2, n: 1, size: 2
			type: 5, id: 3, n: 1, size: 4
			type: 1, id: 5, n: 3, size: 3

		*/
		n := isize << bcfTypeShift[itype]
		if isize == 15 {
			atype, _ := f.typed()
			//fmt.Println("sized")
			// update isize to be the length of the array.
			isize = f.int(atype.nBytes())
		}

		off := f.off
		f.off += nSample * uint32(n)
		e := infoentry{count: nSample * uint32(isize), headerKey: uint16(key),
			etype: itype, buf: f.buf[off:f.off]}
		f.entries = append(f.entries, e)
		//fmt.Println(" after:", f.buf[f.off], "moved:", nSample*uint32(isize))
	}
	return f
}
