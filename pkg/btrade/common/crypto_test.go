package common

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	key = "uy4ymckgirj3nverpa67vsqp7gbf1yg7"
	iv  = "ewfrq37gka4w7pf3"
)

func TestAesCBC(t *testing.T) {
	plain := "131a64d0a4065f570c0740572b6901a6965eb1622849e62f4b60cebd24345667bb8017dd60716aa59549af8a103e1c72b9f945f8753120b8cabef8b3099d765a"

	encBytes, err := AesCBCEncrypt([]byte(plain), []byte(key), []byte(iv))
	assert.Nil(t, err)
	t.Log(hex.EncodeToString(encBytes))

	decBytes, err := AesCBCDecrypt(encBytes, []byte(key), []byte(iv))
	assert.Nil(t, err)
	assert.Equal(t, plain, string(decBytes))
	t.Log(string(decBytes))
}

func TestZHaobiEncrypt(t *testing.T) {
	ed255Key := "6e5e4e87a9eefafd850257ba61e81902ccc7fa634330743d57c0407ffe8c4b4633a3048cbc1ba6ac6360ecb1aa3b9d7451374087ef434847ca33afa1c3c35ce9"
	priv := ed255Key[:len(ed255Key)/2]

	encBytes, err := AesCBCEncrypt([]byte(priv), []byte(key), []byte(iv))
	assert.Nil(t, err)

	dbStr := hex.EncodeToString(encBytes)
	t.Log("db priv2:", dbStr)

	bytes, _ := hex.DecodeString(dbStr)
	decBytes, err := AesCBCDecrypt(bytes, []byte(key), []byte(iv))
	assert.Nil(t, err)
	assert.Equal(t, string(decBytes), priv)
}

func TestDecPrivate2(t *testing.T) {
	priv2 := "4132c02d8fbe7c549e43cc2f0ca3bd197169f98307433b9801d654d1f73768779abacc39731730044c13f7181210150fb79890c19e5f2a60a0930b1fe2edd3c40e3a09051553f7f89dbda139e764f161"
	bytes, _ := hex.DecodeString(priv2)
	decBytes, err := AesCBCDecrypt(bytes, []byte(key), []byte(iv))
	assert.Nil(t, err)
	t.Log(string(decBytes))
}
