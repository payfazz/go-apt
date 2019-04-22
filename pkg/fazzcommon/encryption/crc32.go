package encryption

import (
	"fmt"
	"hash/crc32"
)

// HashCRC32 is a function that used for hashing the message into CRC32
func HashCRC32(message string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x\n", crc32.Checksum([]byte(message), crc32q))
}
