package proto

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) Reader {
	return Reader{r: r}
}

func (rd *Reader) readBuf(buf []byte) error {
	_, err := io.ReadFull(rd.r, buf)
	return err
}

func (rd *Reader) readLen() (int, error) {
	buf := make([]byte, 4)
	err := rd.readBuf(buf[:1])
	if err != nil {
		return 0, err
	}
	c := buf[0]
	switch {
	case (c & 0x80) == 0:
		return int(c), nil
	case (c & 0xc0) == 0x80:
		buf[0] = 0
		buf[2] = c & 0x7f // clear one bit
		err = rd.readBuf(buf[3:])
	case (c & 0xe0) == 0xc0:
		buf[0] = 0
		buf[1] = c & 0x3f // clear two bits
		err = rd.readBuf(buf[2:])
	case (c & 0xf0) == 0xe0:
		buf[0] = c & 0x1f // clear three bits
		err = rd.readBuf(buf[1:])
	default:
		err = rd.readBuf(buf)
	}
	if err != nil {
		return 0, err
	} else if (buf[0] & 0x80) != 0 {
		return 0, fmt.Errorf("invalid len")
	}
	return int(binary.BigEndian.Uint32(buf)), nil
}

func (rd *Reader) ReadWord() (string, error) {
	n, err := rd.readLen()
	if err != nil {
		return "", err
	} else if n == 0 {
		return "", nil
	}
	buf := make([]byte, n)
	if err = rd.readBuf(buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
