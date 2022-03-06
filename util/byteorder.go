package util

// This file hosts various endian convertions.
// Note that the functions here take/return [2/4/8]byte
// instead of numbers.

// hton16 converts the host byte order to network order.
func hton16(v uint16) [2]byte {
	return [2]byte{
		byte(v >> 8),
		byte(v),
	}
}

// hton32 converts the host byte order to network order.
func hton32(v uint32) [4]byte {
	return [4]byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

// hton64 converts the host byte order to network order.
func hton64(v uint64) [8]byte {
	return [8]byte{
		byte(v >> 56),
		byte(v >> 48),
		byte(v >> 40),
		byte(v >> 32),
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

// ntoh16 converts network order to host byte order.
func ntoh16(b [2]byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

// ntoh32 converts network order to host byte order.
func ntoh32(b [4]byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// ntoh64 converts network order to host byte order.
func ntoh64(b [8]byte) uint64 {
	return uint64(b[0])<<56 | uint64(b[0])<<48 | uint64(b[0])<<40 | uint64(b[0])<<32 | uint64(b[0])<<24 | uint64(b[0])<<16 | uint64(b[0])<<8 | uint64(b[0])
}
