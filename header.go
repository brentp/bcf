package bcf

import "bytes"

type headerType int

const (
	infoType headerType = iota
	formatType
)

// HeaderEntry holds the header Info and Format fields
type HeaderEntry struct {
	Id          string
	Description string
	Number      string // A G R . ''
	Type        string // STRING INTEGER FLOAT FLAG CHARACTER UNKONWN
	hType       headerType
}

// Header contains the type info about the BCF
type Header struct {
	version [2]uint8
	lText   uint32
	Text    []byte
	Entries []HeaderEntry
}

func (h *Header) parse() error {
	pieces := bytes.Split(bytes.TrimSpace(h.Text), []byte("\n"))
	for i, p := range pieces {
		// TODO:
		_, _ = i, p

	}
	return nil
}
