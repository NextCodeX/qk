package core

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

// 这里包含一些常用的编码函数

// base64 编码
func (this *InternalFunctionSet) Base64Encode(arg interface{}) string {
	if raw, ok := arg.([]byte); ok {
		return base64.StdEncoding.EncodeToString(raw)
	} else if raw, ok := arg.(string); ok {
		return base64.StdEncoding.EncodeToString([]byte(raw))
	} else {
		fmt.Println("base64Encode() the first parameter type must be ByteArray or String")
	}
	return ""
}
func (this *InternalFunctionSet) Base64(arg interface{}) string {
	return this.Base64Encode(arg)
}

// base64 解码
func (this *InternalFunctionSet) Base64Decode(raw string) []byte {
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		runtimeExcption(err)
	}
	return data
}
func (this *InternalFunctionSet) Debase64(raw string) []byte {
	return this.Base64Decode(raw)
}

// gzip 编码 （压缩）
func (this *InternalFunctionSet) GzipEncode(src interface{}) []byte {
	var raw []byte
	if bs, ok := src.([]byte); ok {
		raw = bs
	} else if str, ok := src.(string); ok {
		raw = []byte(str)
	} else {
		runtimeExcption("GzipEncode() the first parameter type must be ByteArray or String")
		return nil
	}

	return gzipEncode(raw)
}
func (this *InternalFunctionSet) Gzip(src interface{}) []byte {
	return this.GzipEncode(src)
}
func gzipEncode(raw []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(raw); err != nil {
		runtimeExcption(err)
	}
	if err := gz.Flush(); err != nil {
		runtimeExcption(err)
	}
	if err := gz.Close(); err != nil {
		runtimeExcption(err)
	}
	return buf.Bytes()
}

// gzip 解码 （解压缩）
func (this *InternalFunctionSet) GzipDecode(data []byte) []byte {
	return gzipDecode(data)
}
func (this *InternalFunctionSet) Degzip(data []byte) []byte {
	return gzipDecode(data)
}
func gzipDecode(data []byte) []byte {
	bytesReader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		runtimeExcption(err)
	}
	res, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		runtimeExcption(err)
	}
	return res
}

// md5 编码
func (this *InternalFunctionSet) Md5(raw interface{}) string {
	bs := toBytes(raw)
	return md5Encode(bs)
}

func md5Encode(bs []byte) string {
	m := md5.New()
	m.Write(bs)
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func toBytes(arg interface{}) []byte {
	if bs, ok := arg.([]byte); ok {
		return bs
	} else if str, ok := arg.(string); ok {
		return []byte(str)
	} else if ival, ok := arg.(int64); ok {
		return intToBytes(ival)
	} else if fval, ok := arg.(float64); ok {
		return floatToBytes(fval)
	} else {
		return make([]byte, 0, 0)
	}
}

// 对称加密
func (this *InternalFunctionSet) AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		runtimeExcption("AesEncrypt() failed to parse key, key must be base64 form", err)
	}
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

// 对称加密
func (this *InternalFunctionSet) Aes(orig string, key string) string {
	return this.AesEncrypt(orig, key)
}

// 对称解密
func (this *InternalFunctionSet) AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		runtimeExcption("AesDecrypt() failed to parse data, data must be base64 form", err)
	}
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		runtimeExcption("AesDecrypt() failed to parse key, key must be base64 form", err)
	}
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}
func (this *InternalFunctionSet) Deaes(cryted string, key string) string {
	return this.AesDecrypt(cryted, key)
}

// 补码
// AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
