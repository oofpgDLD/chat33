package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/33cn/chat33/pkg/btrade/crypto"
	_ "github.com/33cn/chat33/pkg/btrade/crypto/ed25519"
	"github.com/tendermint/ed25519"
)

var (
	ErrSign = errors.New("ErrSign")
)

func RandInstructionId() int64 {
	//return int64(RandUint64() >> 1)
	return RandInt64()
}

func HexToPrivkey(priv string) *[64]byte {
	bytes, err := hex.DecodeString(priv)
	if err != nil {
		panic(err)
	}
	if len(bytes) != 64 {
		panic("priv key len is not 64")
	}
	var data [64]byte
	copy(data[:], bytes)
	return &data
}

// =======================================================

// BECAREFUL: The following 4 func only appiled to "ed25519"

func GetPub(key *[64]byte) string {
	/*
		uid := hex.EncodeToString(key[32:])
		return uid
	*/

	c, err := crypto.New("ed25519")
	if err != nil {
		panic(err.Error())
	}

	pub, err := c.PubKeyFromBytes(key[32:])
	if err != nil {
		panic(err.Error())
	}

	return pub.KeyString()
}

func GenKey() (priv *[64]byte, err error) {
	/*
		_, priv, err = ed25519.GenerateKey(crand.Reader)
		if err != nil {
			return nil, err
		}
		return priv, nil
	*/
	c, err := crypto.New("ed25519")
	if err != nil {
		return nil, err
	}

	privKey, err := c.GenKey()
	if err != nil {
		return nil, err
	}

	b := privKey.Bytes()
	priv = &[64]byte{}
	copy(priv[:], b)
	return
}

func GenKeyWithPassword(password string) (*[64]byte, error) {
	h := sha256.New()

	_, err := h.Write([]byte(password))
	if nil != err {
		return nil, err
	}

	privatekey := h.Sum(nil)

	userkey := new([64]byte)
	copy(userkey[:32], privatekey)

	ed25519.MakePublicKey(userkey)

	return userkey, nil
}

func Signdata(priv *[64]byte, data []byte) []byte {
	/*
		var fakeCrypto = false
		if fakeCrypto {
			sign := make([]byte, 64)
			return sign[:]
		}
		sign := ed25519.Sign(priv, data)
		return sign[:]
	*/
	c, err := crypto.New("ed25519")
	if err != nil {
		panic(err.Error())
	}

	privKey, err := c.PrivKeyFromBytes(priv[:])
	if err != nil {
		panic(err.Error())
	}

	sign := privKey.Sign(data)
	if sign == nil {
		panic("Sign ed25519 error")
	}

	return sign.Bytes()
}

func CheckSign(data []byte, uid []byte, sign []byte) error {
	c, err := crypto.New("ed25519")
	if err != nil {
		return err
	}
	sig, err := c.SignatureFromBytes(sign)
	if err != nil {
		return err
	}
	pub, err := c.PubKeyFromBytes(uid[:])
	if err != nil {
		return err
	}
	if !pub.VerifyBytes(data, sig) {
		return ErrSign
	}
	return nil
}

func CheckPub(bytes []byte, crptoType string) error {
	c, err := crypto.New(crptoType)
	if err != nil {
		return err
	}
	if _, err := c.PubKeyFromBytes(bytes); err != nil {
		return err
	}
	return nil
}

func AesCBCEncrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS5Padding(origData, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	mode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCBCDecrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	encryptData = PKCS5UnPadding(encryptData)
	return encryptData, nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
