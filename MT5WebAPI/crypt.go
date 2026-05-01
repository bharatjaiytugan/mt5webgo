//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// MTCrypt handles AES-256 OFB mode encryption/decryption for MT5 protocol packets.
type MTCrypt struct {
	cryptIV    [16][]byte
	aesOut     []byte
	aesIn      []byte
	cryptAes   *MTCryptAES
	deCryptAes *MTCryptAES
}

// SetCryptAESRand initializes the AES encryption keys from the server's random
// challenge and the user password.
func (c *MTCrypt) SetCryptAESRand(cryptRand []byte, password string) {
	tempRand := GetHashFromPasswordCryptRand(password, cryptRand)
	for i := 0; i < 16; i++ {
		tempRand = GetMD5ByteResult(CopyBufferOffset(cryptRand, i*16, 16, tempRand))
		c.cryptIV[i] = tempRand
	}
}

// DeCryptPacket decrypts a protocol packet body using AES-256 in OFB mode.
func (c *MTCrypt) DeCryptPacket(packetBody []byte) []byte {
	if packetBody == nil {
		return nil
	}
	// lazy-init decrypt AES
	if c.deCryptAes == nil {
		keyTemp := CopyBuffer(c.cryptIV[0], c.cryptIV[1])
		if keyTemp == nil {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, "packet did not decrypt, decrypt key is empty")
			}
			return nil
		}
		c.deCryptAes = NewMTCryptAES(keyTemp)
		c.aesIn = c.cryptIV[3]
	}
	if c.aesIn == nil || len(c.aesIn) == 0 {
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "packet did not decrypt, aes in is empty")
		}
		return nil
	}
	result := make([]byte, len(packetBody))
	key := 16
	for i := 0; i < len(packetBody); i++ {
		if key >= 16 {
			c.aesIn = c.deCryptAes.EncryptBlock(c.aesIn)
			key = 0
		}
		result[i] = packetBody[i] ^ c.aesIn[key]
		key++
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeError, fmt.Sprintf("decrypt: '%s' to '%s'", GetHex(packetBody), GetString(result)))
	}
	return result
}

// CryptPacket encrypts a protocol packet body using AES-256 in OFB mode.
func (c *MTCrypt) CryptPacket(packetBody []byte) []byte {
	if packetBody == nil {
		return nil
	}
	result := make([]byte, len(packetBody))
	// lazy-init encrypt AES
	if c.cryptAes == nil {
		keyTemp := CopyBuffer(c.cryptIV[0], c.cryptIV[1])
		if keyTemp == nil {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, "packet did not crypt, crypt key is empty")
			}
			return nil
		}
		c.cryptAes = NewMTCryptAES(keyTemp)
		c.aesOut = c.cryptIV[2]
	}
	if c.aesOut == nil || len(c.aesOut) == 0 {
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "packet did not crypt, aes out is empty")
		}
		return nil
	}
	key := 16
	for i := 0; i < len(packetBody); i++ {
		if key >= 16 {
			c.aesOut = c.cryptAes.EncryptBlock(c.aesOut)
			key = 0
		}
		result[i] = packetBody[i] ^ c.aesOut[key]
		key++
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeError, fmt.Sprintf("crypt: '%s' to '%s'", GetString(packetBody), GetHex(result)))
	}
	return result
}
