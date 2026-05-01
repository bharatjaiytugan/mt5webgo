//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf16"
)

// webAPIWord is the constant used in password hashing.
const webAPIWord = "WebAPI"

// GetFromHex parses a hex string into a byte slice.
func GetFromHex(hexStr string) []byte {
	if len(hexStr) == 0 {
		return nil
	}
	if len(hexStr)%2 != 0 {
		return nil
	}
	result, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil
	}
	return result
}

// GetRandomHex generates a byte slice of the given length with random values.
func GetRandomHex(length int) []byte {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = byte(r.Intn(255))
	}
	return result
}

// GetMD5String computes the MD5 hash of a string (encoded as UTF-16LE, matching C# Encoding.Unicode)
// and returns the hex string.
func GetMD5String(input string) string {
	// C# Encoding.Unicode is UTF-16LE
	utf16Codes := utf16.Encode([]rune(input))
	buf := make([]byte, len(utf16Codes)*2)
	for i, code := range utf16Codes {
		binary.LittleEndian.PutUint16(buf[i*2:], code)
	}
	hash := md5.Sum(buf)
	return hex.EncodeToString(hash[:])
}

// GetMD5Bytes computes the MD5 hash of a byte slice and returns the hex string.
func GetMD5Bytes(input []byte) string {
	hash := md5.Sum(input)
	return hex.EncodeToString(hash[:])
}

// GetMD5ByteResult computes the MD5 hash of a byte slice and returns the raw hash bytes.
func GetMD5ByteResult(input []byte) []byte {
	if input == nil {
		return nil
	}
	hash := md5.Sum(input)
	result := make([]byte, len(hash))
	copy(result, hash[:])
	return result
}

// GetMD5ByteRange computes the MD5 hash of a sub-range of a byte slice.
func GetMD5ByteRange(input []byte, offset, count int) []byte {
	if input == nil {
		return nil
	}
	hash := md5.Sum(input[offset : offset+count])
	result := make([]byte, len(hash))
	copy(result, hash[:])
	return result
}

// GetHashFromPassword computes the password hash used for MT5 authentication.
func GetHashFromPassword(password string, randCode []byte) []byte {
	if randCode == nil {
		return nil
	}
	// hash of password (UTF-16LE encoded)
	utf16Codes := utf16.Encode([]rune(password))
	pwdBytes := make([]byte, len(utf16Codes)*2)
	for i, code := range utf16Codes {
		binary.LittleEndian.PutUint16(pwdBytes[i*2:], code)
	}
	hash := GetMD5ByteResult(pwdBytes)
	if hash == nil {
		return nil
	}
	// concatenate hash + "WebAPI" (UTF-8)
	apiWord := []byte(webAPIWord)
	hashContains := CopyBuffer(hash, apiWord)
	if hashContains == nil {
		return nil
	}
	hash = GetMD5ByteResult(hashContains)
	if hash == nil {
		return nil
	}
	// concatenate hash + randCode
	hashContains = CopyBuffer(hash, randCode)
	if hashContains == nil {
		return nil
	}
	// final hash
	return GetMD5ByteResult(hashContains)
}

// GetHashFromPasswordCryptRand computes the hash used for setting up AES encryption.
func GetHashFromPasswordCryptRand(password string, randCode []byte) []byte {
	if randCode == nil {
		return nil
	}
	// hash of password (UTF-16LE encoded)
	utf16Codes := utf16.Encode([]rune(password))
	pwdBytes := make([]byte, len(utf16Codes)*2)
	for i, code := range utf16Codes {
		binary.LittleEndian.PutUint16(pwdBytes[i*2:], code)
	}
	hash := GetMD5ByteResult(pwdBytes)
	if hash == nil {
		return nil
	}
	// concatenate hash + "WebAPI" (UTF-8)
	apiWord := []byte(webAPIWord)
	hashContains := CopyBuffer(hash, apiWord)
	if hashContains == nil {
		return nil
	}
	return GetMD5ByteResult(hashContains)
}

// CopyBuffer concatenates two byte slices into a new byte slice.
func CopyBuffer(ar1, ar2 []byte) []byte {
	result := make([]byte, len(ar1)+len(ar2))
	copy(result, ar1)
	copy(result[len(ar1):], ar2)
	return result
}

// CopyBufferOffset concatenates a sub-range of ar1 with all of ar2 into a new byte slice.
func CopyBufferOffset(ar1 []byte, offset, count int, ar2 []byte) []byte {
	if ar2 == nil {
		return ar1
	}
	result := make([]byte, count+len(ar2))
	copy(result, ar1[offset:offset+count])
	copy(result[count:], ar2)
	return result
}

// CopyBufferList concatenates a list of byte slices into a single byte slice of the given length.
func CopyBufferList(list [][]byte, length int64) []byte {
	if len(list) == 0 {
		return nil
	}
	result := make([]byte, length)
	currentPos := 0
	for _, ar := range list {
		copy(result[currentPos:], ar)
		currentPos += len(ar)
	}
	return result
}

// CompareBytes compares two byte slices for equality.
func CompareBytes(ar1, ar2 []byte) bool {
	if ar1 == nil {
		return ar2 == nil
	}
	if ar2 == nil {
		return false
	}
	if len(ar1) != len(ar2) {
		return false
	}
	for i := range ar1 {
		if ar1[i] != ar2[i] {
			return false
		}
	}
	return true
}

// Quotes escapes special symbols: =, |, \n, \ by prepending a backslash.
// The order matches the C# implementation exactly: =, |, \n, \ .
func Quotes(str string) string {
	if len(str) == 0 {
		return str
	}
	str = strings.ReplaceAll(str, "=", "\\=")
	str = strings.ReplaceAll(str, "|", "\\|")
	str = strings.ReplaceAll(str, "\n", "\\\n")
	str = strings.ReplaceAll(str, "\\", "\\\\")
	return str
}

// GetString converts a UTF-16LE encoded byte slice to a Go string.
func GetString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if len(b)%2 != 0 {
		logWrite(LogTypeError, fmt.Sprintf("get string from bytes failed: odd byte count %d", len(b)))
		return ""
	}
	u16s := make([]uint16, len(b)/2)
	for i := range u16s {
		u16s[i] = binary.LittleEndian.Uint16(b[i*2:])
	}
	return string(utf16.Decode(u16s))
}

// GetHex returns the lowercase hex string representation of a byte slice.
func GetHex(hash []byte) string {
	return hex.EncodeToString(hash)
}

// GetHexRange returns the lowercase hex string for a sub-range of a byte slice.
func GetHexRange(hash []byte, begin, count int) string {
	end := begin + count
	if end > len(hash) {
		end = len(hash)
	}
	return hex.EncodeToString(hash[begin:end])
}

// ConvertToUnixTimestamp converts a time.Time to Unix timestamp (seconds since 1970).
func ConvertToUnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

// ConvertFromUnixTime converts a Unix timestamp (seconds since 1970) to time.Time.
func ConvertFromUnixTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}

// ConvertToOldVolume converts a new 8-digit volume to old 4-digit volume by dividing by 10000.
func ConvertToOldVolume(newVolume uint64) uint64 {
	return newVolume / 10000
}

// ConvertToNewVolume converts an old 4-digit volume to new 8-digit volume by multiplying by 10000.
func ConvertToNewVolume(oldVolume uint64) uint64 {
	return oldVolume * 10000
}
