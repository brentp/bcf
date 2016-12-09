package bcf

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

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

var typeRe = `String|Integer|Float|Flag|Character|Unknown`
var hRegexp = regexp.MustCompile(fmt.Sprintf(`=<ID=(.+),Number=([\dAGR\.]?),Type=(%s  ),Description="(.+)".+>`, typeRe))

type HeaderError struct {
	Line int
	Msg  string
}

func (h HeaderError) Error() string {
	return h.Msg
}

func parseHeaderFormatInfo(iline string) (HeaderEntry, error) {
	line := iline[strings.Index(iline, "="):]
	res := hRegexp.FindStringSubmatch(line)
	var i HeaderEntry
	if len(res) != 5 {
		return i, fmt.Errorf("bcf: error in header: %s", iline)
	}
	i.Id = res[1]
	i.Number = res[2]
	i.Type = res[3]
	i.Description = res[4]
	return i, nil
}

func (h *Header) parse() error {
	pieces := bytes.Split(bytes.TrimSpace(h.Text), []byte("\n"))
	for i, line := range pieces {
		if bytes.HasPrefix(line, []byte("##FORMAT")) {
			he, err := parseHeaderFormatInfo(string(line))
			if err != nil {
				return HeaderError{Msg: err.Error(), Line: i}
			}
			he.hType = formatType
		} else if bytes.HasPrefix(line, []byte("##INFO")) {
			he, err := parseHeaderFormatInfo(string(line))
			if err != nil {
				return HeaderError{Msg: err.Error(), Line: i}
			}
			he.hType = infoType
		}
	}
	return nil
}
