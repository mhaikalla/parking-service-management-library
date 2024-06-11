package crypts

import "errors"

// Pad add padding to `buf` based on `size`.
// For compability with PKCS #5, `size` is size of AES block cipher.
// Size of AES block cipher is 16 byte.
func Pad(buf []byte, size int) ([]byte, error) {
	bufLen := len(buf)                    // length of the byte slice target
	padLen := size - bufLen%size          // length of pad we need to add
	padded := make([]byte, bufLen+padLen) // create a new byte slice with length of byte slice + length of pad
	copy(padded, buf)                     // add byte slice target to a new created byte slice above
	for i := 0; i < padLen; i++ {         // let fill the rest of data
		padded[bufLen+i] = byte(padLen)
	}
	return padded, nil // return the result
}

// Unpad remove padding to `buf` based on `size`.
// For compability with PKCS #5, `size` is size of AES block cipher.
// Size of AES block cipher is 16 byte.
func Unpad(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 { // check if length of byte slice is divided by the size
		return nil, errors.New("pkcs7: Padded value wasn't in correct size")
	}

	bufLen := len(padded) - int(padded[len(padded)-1]) // calculate the size of padded char

	buf := make([]byte, bufLen) // create a new byte slice
	copy(buf, padded[:bufLen])  // copy the result so we're not changing original byte slice
	return buf, nil             // return result
}
