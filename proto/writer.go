package proto

import (
	"encoding/binary"
	"io"

	"github.com/valyala/bytebufferpool"
)

type Writer struct {
	w io.Writer
}

var bbp bytebufferpool.Pool

func NewWriter(w io.Writer) Writer {
	return Writer{w: w}
}

func between(n, lo, hi int) bool {
	return lo <= n && n <= hi
}

func encodeLen(n int) []byte {
	if between(n, 0, 0x7f) {
		return []byte{byte(n)}
	} else if between(n, 0x80, 0x3fff) {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(n)|0x8000)
		return buf
	} else if between(n, 0x4000, 0x1fffff) {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(n)|0xc00000)
		return buf[1:]
	} else if between(n, 0x200000, 0xfffffff) {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(n)|0xe0000000)
		return buf
	} else if n >= 0x10000000 {
		buf := make([]byte, 5)
		buf[0] = 0xf0
		binary.BigEndian.PutUint32(buf[1:], uint32(n))
		return buf
	} else {
		panic("negative len")
	}
}

func writeWord(bb *bytebufferpool.ByteBuffer, str string) {
	bb.Write(encodeLen(len(str)))
	bb.WriteString(str)
}

func (wr *Writer) WriteWords(str ...string) error {
	bb := bbp.Get()
	defer bbp.Put(bb)
	for _, s := range str {
		writeWord(bb, s)
	}
	bb.WriteByte(0)
	_, err := wr.w.Write(bb.Bytes())
	return err
}
