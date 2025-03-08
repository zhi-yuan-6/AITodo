package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// 生成pem私钥
func GeneratePEM() {
	//生成ECDSA私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	//将生成的私钥（privateKey）序列化为PKCS#8 格式
	bytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)

	//将私钥转换为PEM格式
	var privateKeyPEM = &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: bytes,
	}

	//输出PEM格式的私钥
	fmt.Println(string(pem.EncodeToMemory(privateKeyPEM)))

	//使用私钥作为jwtSecret
	JwtSecret = privateKey
}

func pemToECDSAPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	// 解析 PEM 数据块
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 数据")
	}

	// 解析私钥
	var key interface{}
	var err error

	// 尝试作为 PKCS#8 私钥解析
	if key, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		// 尝试作为 PKCS#1 私钥解析
		if key, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
			return nil, fmt.Errorf("无法解析私钥: %v", err)
		}
	}

	// 将解析后的密钥断言为 *ecdsa.PrivateKey
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("密钥类型不匹配，无法断言为 *ecdsa.PrivateKey")
	}

	return ecdsaKey, nil
}

func ReadPEM(fileName string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// 从文件中读取 PEM 数据
	pemFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("读取 PEM 文件失败: %v\n", err)
	}

	// 将 PEM 数据转换为 ECDSA 私钥
	privateKey, err := pemToECDSAPrivateKey(pemFile)
	if err != nil {
		return nil, nil, fmt.Errorf("转换私钥失败: %v\n", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}
