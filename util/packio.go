package util

import (
	"fmt"
	"io"
)

// This file contains functions for endian-netural network I/O.

var MaxLength int64 = 4096 // Maximum length for ReadString/ReadBytes

// errTooLong is retured by ReadString/ReadBytes if the length exceeds MaxLength.
type errTooLong int64

var _ error = errTooLong(0)

// Error implements error.Error.
func (e errTooLong) Error() string {
	return fmt.Sprintf("read string/bytes too long (%d, max %d)", int64(e), MaxLength)
}

// Reading

// ReadInt16 reads a Int16 from the reader.
func ReadInt16(reader io.Reader) (result int16, err error) {
	var buf [2]byte
	n, err := io.ReadFull(reader, buf[:])
	if err != nil || n == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	return int16(ntoh16(buf)), nil
}

// ReadInt32 reads a Int16 from the reader.
func ReadInt32(reader io.Reader) (result int32, err error) {
	var buf [4]byte
	n, err := io.ReadFull(reader, buf[:])
	if err != nil || n == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	return int32(ntoh32(buf)), nil
}

// ReadInt64 reads a Int16 from the reader.
func ReadInt64(reader io.Reader) (result int64, err error) {
	var buf [8]byte
	n, err := io.ReadFull(reader, buf[:])
	if err != nil || n == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	return int64(ntoh64(buf)), nil
}

// ReadString reads a String from the reader.
//
// Don't use this when you can use ReadBytes.
// It essentialy is a shortcut to string(ReadBytes(...)).
func ReadString(reader io.Reader) (result string, err error) {
	buf, err := ReadBytes(reader)
	return string(buf), err
}

// ReadBytes reads a byte slice from the reader.
func ReadBytes(reader io.Reader) (result []byte, err error) {
	n, err := ReadInt64(reader)
	if err != nil {
		return
	}

	if n > MaxLength {
		return nil, errTooLong(n)
	}
	if n == 0 {
		return nil, nil
	}

	result = make([]byte, n)
	nread, err := io.ReadFull(reader, result)
	if err != nil || nread == 0 {
		return nil, io.ErrUnexpectedEOF
	}

	return
}

// Writing

// WriteInt16 writes a Int16 to the writer.
func WriteInt16(writer io.Writer, val int16) (err error) {
	b := hton16(uint16(val))
	_, err = writer.Write(b[:])
	return
}

// WriteInt32 writes a Int16 to the writer.
func WriteInt32(writer io.Writer, val int32) (err error) {
	b := hton32(uint32(val))
	_, err = writer.Write(b[:])
	return
}

// WriteInt64 writes a Int16 to the writer.
func WriteInt64(writer io.Writer, val int64) (err error) {
	b := hton64(uint64(val))
	_, err = writer.Write(b[:])
	return
}

// WriteBytes writes a slice of bytes to the writer.
//
// It writes the length of the slice before the data
// so that it can be retrieved by ReadBytes/ReadString.
func WriteBytes(writer io.Writer, data []byte) (err error) {
	err = WriteInt64(writer, int64(len(data)))
	if err != nil || len(data) == 0 {
		return
	}
	_, err = writer.Write(data)
	return
}

// WriteString writes a string to the writer.
//
// It writes the length of the slice before the data
// so that it can be retrieved by ReadBytes/ReadString.
func WriteString(writer io.Writer, data string) (err error) {
	err = WriteInt64(writer, int64(len(data)))
	if err != nil || len(data) == 0 {
		return
	}
	_, err = io.WriteString(writer, data)
	return
}
