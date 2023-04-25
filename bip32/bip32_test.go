package bip32

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

var goodSeed = []byte("seedlen seedlen seedlen seedlen seedlen seedlen seedlen seedlen ")

func TestGoodPath(t *testing.T) {
	path := "m/44'/540'/2'/0'/0'"
	_, err := Derive(path, goodSeed)
	require.NoError(t, err)
}

func TestBadPath(t *testing.T) {
	path := "bad path"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestEmptyPath(t *testing.T) {
	path := ""
	key, err := Derive(path, goodSeed)
	require.Equal(t, pathErr, err)
	require.Nil(t, key)
}

func TestEmptySeed(t *testing.T) {
	path := "hello world"
	var seed []byte
	key, err := Derive(path, seed)
	require.Equal(t, badSeedLen, err)
	require.Nil(t, key)
}

func TestSeedLength(t *testing.T) {
	path := "hello world"
	seedTooShort := goodSeed[:len(goodSeed)-1]
	key, err := Derive(path, seedTooShort)
	require.Equal(t, badSeedLen, err)
	require.Nil(t, key)
	seedTooLong := append(goodSeed, byte(1))
	key, err = Derive(path, seedTooLong)
	require.Equal(t, badSeedLen, err)
	require.Nil(t, key)
}

func TestNonHardenedPath(t *testing.T) {
	path := "m/44'/540'/2/0'/0'"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathMalformed(t *testing.T) {
	path := "m/44'/540'/2'/0'/"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
	path = "m/44'/540'/2'/0'/1'/"
	key, err = Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathShort(t *testing.T) {
	path := "m/44'"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
	path = "m"
	key, err = Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathLong(t *testing.T) {
	path := "m/44'/540'/2'/0'/1'/2'"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathBadPurpose(t *testing.T) {
	path := "m/41'/540'/2'/0'/1'"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathBadCoinType(t *testing.T) {
	path := "m/44'/542'/2'/0'/1'"
	key, err := Derive(path, goodSeed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestSimple(t *testing.T) {
	seed1 := "q9ONfCHWlAwY5Ea4AIkePeA0dUSYnyqcNVElvREhYJ0OU7kvwg9MroDKTcmSxTKt"
	seed2 := "gk3YIKmq0j5rqjLBrJVaQfC6k7rIHLp96G2ZGEg1vl0cwv7LNsjOaPOTAJEzmWVi"
	seed3 := "wZtfpYa1lzcdHTp9UwdbzTLlMnFTZee5Y7dLN7e795Kfr4U1E7AQaUcqGnMNEraR"
	// Test vectors
	vectors := []struct{ path, seed, key string }{
		{
			"m/44'/540'/2'/0'/0'",
			seed1,
			"02c6197bee36ee341f1bad9bbb26b71b9b51c78f4a6aebd722d6ccda9dbfc7a6cabdd9afe0274ba6eb6b4ea6b04f5110e13321ae3c1316482aace2a8bff26c89",
		},
		{
			"m/44'/540'/1'/0'/0'",
			seed1,
			"856761857b51d7ab4c8e86fbabe1fbdce2eb2851cc70389ee914a375617daaa86cb278b215158b31b6f96f0516f05d52d416f179758858f329ae9a08995ffef9",
		},
		{
			"m/44'/540'/1'/1'/1'",
			seed1,
			"417f8cf19c44bc3398fa789bb1d0697c56c872c4702bfacaffa6b4132f39119b4562f9f55ea6836f451100ceac42be3cb9f0022b5677d12460c09651b05d160c",
		},
		{
			"m/44'/540'/100'/0'/0'",
			seed1,
			"39a84a1eed038a9560529646615a512847a09de0ba52b5b388c5aa949a8a77c0d08480429cc3c12c15b6abf14276967069eec285a264beaac23acf3789c25932",
		},
		{
			"m/44'/540'/2'/100'/200'",
			seed1,
			"bf04d7fcc2cd70904b119e67a20c33039cf5355b2837ddabe28a4a53d0df6b830046e8f9a1150560bdba11cd8e3ae652cd4f77de6603ab8ad964a775ed55f9da",
		},
		{
			"m/44'/540'/0'/0'/0'",
			seed2,
			"477179d9a8a3571aa068d7a10da60cae617f3eb232b493bfd079842c3f0e69662978db84488e3644a6a986dd52622440a0441206e4c24a24f115336819fecaa7",
		},
		{
			"m/44'/540'/1'/0'/0'",
			seed2,
			"14d2c7daefe19e395328efcbcefbbb854d73890a77bb300ee310a71d0ea76826c106af1346f27fbb9f4623e9cc362fd3c9ad6e139652685e30f35794d52d190d",
		},
		{
			"m/44'/540'/1'/1'/2'",
			seed3,
			"81e65d2c6806d32a223a64200c22217408a071896cab6e657e729d724295013e47f73f6b47a70373961a41cdc8130ac7a5887ec18ecf32237ca8cb6e7ba069b9",
		},
		{
			"m/44'/540'/1'/0'",
			seed1,
			"03d65083ada328fa74637cc8588ddc80ef6fd732e60e84715c551221c2313b88fb53139aa915339ec0250a47725345b11c79cf9a39437e8dcd107cf17ab7416d",
		},
		{
			"m/44'/540'/1'",
			seed2,
			"1831d61e76a8706f20614b96b74b50b2214a2bbf398735dd10c720d214b78c0c2e318aefd286a646ec9c6313bb036132e4893f4985c1ed99d54a555f1a7b72b0",
		},
	}
	for i, vec := range vectors {
		key, err := Derive(vec.path, []byte(vec.seed))
		require.NoError(t, err)
		assert.Equalf(t, vec.key, hex.EncodeToString(key[:]), "test case %d path %s seed %s got unexpected result", i, vec.path, vec.seed)
	}
}

//func TestSimpleChild(t *testing.T) {
//	seed := []byte("hello world")
//	key, err := DeriveChild(seed, 0)
//	require.NoError(t, err)
//	require.Equal(t, "00b4b57f20439c858cd66edaa77fa6b56fd4b8617a614e97829aedcc31aa82a0ee75ecabd996fc882b1b7b13c4bcabc256d123ade062ee767b41ac1a489d04d9", hex.EncodeToString(key[:]))
//}
//
//// Test that deriving a full-path key in one step produces the same result from deriving one in two steps
//func TestIncrementalDerivation(t *testing.T) {
//	seed := []byte("hello world")
//	key, err := FromSeed(seed)
//	require.NoError(t, err)
//	require.Equal(t, "309649976bcd6a2f5e8247ca5cf72c566d8d6d2211eb471ca65b542c2635106f4ab69140806a86f8ac278e437dd1703cbc088a35528f39a0115c19a9e1f9ac87", hex.EncodeToString(key[:]))
//
//	// now derive a child key
//	path := "m/44'/540'/0'/0'/0'"
//	key, err = DeriveChild()
//}
