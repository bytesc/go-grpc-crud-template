package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func RsaDecode(encryptedData string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	block, _ := pem.Decode([]byte(KeyForPwd.PrivateKey))
	priKey, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		return "", errors.New("解析私钥失败")
	}
	originalData, encryptErr := rsa.DecryptPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), encryptedDecodeBytes)
	if encryptErr != nil {
		fmt.Println(encryptErr.Error())
	}
	return string(originalData), encryptErr
}

func RsaEncode(plainData string) (string, error) {
	// 解析公钥
	block, _ := pem.Decode([]byte(KeyForPwd.PublicKey))
	if block == nil {
		return "", errors.New("解析公钥失败")
	}

	// 尝试使用ParsePKIXPublicKey解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 对数据进行RSA加密
	encryptedData, encryptErr := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(plainData))
	if encryptErr != nil {
		return "", encryptErr
	}

	// 对加密后的数据进行Base64编码
	encryptedDataBase64 := base64.StdEncoding.EncodeToString(encryptedData)
	return encryptedDataBase64, nil
}
