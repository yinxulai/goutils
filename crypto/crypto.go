package crypto

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// RSAEncryptWithSha1Base64 （1）加密：采用sha1算法加密后转base64格式
func RSAEncryptWithSha1Base64(originalData, publicKey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(publicKey)                                                  // 转码
	pubKey, _ := x509.ParsePKIXPublicKey(key)                                                             // 生成 x509 公钥
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData)) // 加密
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

// RSADecryptWithSha1Base64 （2）解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
func RSADecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	key, _ := base64.StdEncoding.DecodeString(privateKey)
	prvKey, _ := x509.ParsePKCS1PrivateKey(key)
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	return string(originalData), err
}

// RSASignWithSha1Hex （3）签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func RSASignWithSha1Hex(data string, prvKey string) (string, error) {
	keyByts, err := hex.DecodeString(prvKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyByts)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}

// RSAVerySignWithSha1Base64 （4）验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func RSAVerySignWithSha1Base64(originalData, signData, pubKey string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

// MD5Encrypt 加密
func MD5Encrypt(originalData string, pubKey string) string {
	h := md5.New()
	h.Write([]byte(originalData))
	h.Write([]byte(pubKey))
	return hex.EncodeToString(h.Sum(nil))
}
