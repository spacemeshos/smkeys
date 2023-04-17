package bip32

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoodPath(t *testing.T) {
	path := "m/44'/540'/2'/0'/0'"
	seed := []byte("hello world")
	_, err := Derive(path, seed)
	require.NoError(t, err)
}

func TestBadPath(t *testing.T) {
	path := "bad path"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestEmptyPath(t *testing.T) {
	path := ""
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pathErr, err)
	require.Nil(t, key)
}

func TestEmptySeed(t *testing.T) {
	path := "hello world"
	var seed []byte
	key, err := Derive(path, seed)
	require.Equal(t, seedErr, err)
	require.Nil(t, key)
}

func TestNonHardenedPath(t *testing.T) {
	path := "m/44'/540'/2/0'/0'"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathMalformed(t *testing.T) {
	path := "m/44'/540'/2'/0'/"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
	path = "m/44'/540'/2'/0'/1'/"
	seed = []byte("hello world")
	key, err = Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathShort(t *testing.T) {
	path := "m/44'/540'/2'/0'"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
	path = "m/44'/540'/2'"
	key, err = Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathLong(t *testing.T) {
	path := "m/44'/540'/2'/0'/1'/2'"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathBadPurpose(t *testing.T) {
	path := "m/41'/540'/2'/0'/1'"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestPathBadCoinType(t *testing.T) {
	path := "m/44'/542'/2'/0'/1'"
	seed := []byte("hello world")
	key, err := Derive(path, seed)
	require.Equal(t, pointerErr, err)
	require.Nil(t, key)
}

func TestSimple(t *testing.T) {
	// Test vectors
	vectors := []struct{ path, seed, key string }{
		{
			"m/44'/540'/2'/0'/0'",
			"hello world",
			"7158d5bcf4a5811ff0a81bbad619e140ea2c0f46224d5ee54bc257faa4e34785a9dde42835b420a93956b8cfc26945a97d99f676a4c0b5ced0081c2e39b2c5e5",
		},
		{
			"m/44'/540'/1'/0'/0'",
			"hello world",
			"ea9a4a7d7ced401b6f0bef20c4cc9dc726723191ae96487a473fda91c3f988189c2139abee768f643bcb2dae38d970ae32b682e8e328c8d34d894cd22373da15",
		},
		{
			"m/44'/540'/1'/1'/1'",
			"hello world",
			"77d785e700e30ccc28ae59de74dd31d5a14674f4ded114263caebdaebad67cc5b1dc7675098476bdddcfa7c02d16daf6e8c9a748c3c2141e040931ec7f9b6ec6",
		},
		{
			"m/44'/540'/100'/0'/0'",
			"hello world",
			"d1b593331f9c8c8ece8db6e755a3eeab1e5989a349aab80e136590c23ad606ed22b9c319998614a0869ddd49f31b4d2f56035cc935b0e38201e470467346a0ad",
		},
		{
			"m/44'/540'/2'/100'/200'",
			"hello world",
			"a22ccc3e60f5bb17173a9326beead0e9607989a0ac03e93f5dd378eab34382e9703aa6a75c658dff490ebd34cfac9cd71bf75ad87adec4407fcd66458bf1bc5f",
		},
		{
			"m/44'/540'/0'/0'/0'",
			"goodbye world",
			"f06132d2c4bf8b8c46d1848ef14d99df5444bd194b8f894fc9373e7b4f2f4fe7fe6bed30f5839c2a49d588a82696bee203244f88549b0d088a28a04d1ba829ab",
		},
		{
			"m/44'/540'/1'/0'/0'",
			"goodbye world",
			"304ccb96364ca1c6f56b9f1ee645ea0c0d6d31d29723c5ab3ec673242cdcb995e9cc5f5a18130d3208cd56079f4250da97dcc3c8cf60c9cc463321d92c83056f",
		},
		//{
		//	"m/44'/540'/1'/0'/0'/0'",
		//	"goodbye world",
		//	"0bbc325b65ff280134dd961f001803bcf9bb935a1dfd0b3aaa636fb3e77b3576d5b021bd2424066d2af34d78003c3d780d89742cb7ce697d3fae6c901b58ba81",
		//},
		{
			"m/44'/540'/1'/1'/2'",
			"what a wonderful world",
			"2e72f645bf1a24aee0dc467988f8bf0e284a205144919645351f3dc461bebba52e98256e26223678c5c2ea065f92abc3f9821032ed42286570536da38bce4862",
		},
		//{
		//	"m/44'/540'/1'/0'",
		//	"what a wonderful world",
		//	"f2b16f208262a3fb6476b30727f3c96b5f3517ac0404f88fe8d5e74ed7796f5918a42b368e4ba042886186abd6b67e2eab1d6c54ffe47cf200658a1e06777d1d",
		//},
		//{
		//	"m/44'/540'/1'",
		//	"what a wonderful world",
		//	"06deb1d30a0621f963059b1178bc9fbac84ee9ab9ccb24463f541b023011ef81aae8b6ae91c5fe6701151bb571f39f3161decb1535e61a5b4bfe908319e5348c",
		//},
	}
	for i, vec := range vectors {
		key, err := Derive(vec.path, []byte(vec.seed))
		require.NoError(t, err)
		require.Equalf(t, vec.key, hex.EncodeToString(key[:]), "test case %d path %s seed %s got unexpected result", i, vec.path, vec.seed)
	}
}

func TestSimpleChild(t *testing.T) {
	seed := []byte("hello world")
	key, err := DeriveChild(seed, 0)
	require.NoError(t, err)
	require.Equal(t, "00b4b57f20439c858cd66edaa77fa6b56fd4b8617a614e97829aedcc31aa82a0ee75ecabd996fc882b1b7b13c4bcabc256d123ade062ee767b41ac1a489d04d9", hex.EncodeToString(key[:]))
}
