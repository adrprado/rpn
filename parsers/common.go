package parsers

import (
	"hash/fnv"
)

var fnvHash = fnv.New32a()

//
// getHash returns the FNV-1 non-cryptographic hash
//
func getHash(s string) uint32 {
	fnvHash.Write([]byte(s))
	defer fnvHash.Reset()

	return fnvHash.Sum32()
}
