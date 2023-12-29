package cryptor

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

/*
GetCertsAndKeys создает и возвращает сертификат, публичный и приватные ключи в формате PEM для переданной организации.
Сертификат содержит информацию о владельце, разрешенные IP-адреса (localhost для ТЗ), даты начала и окончания
действия сертификата, а также использование ключа для цифровой подписи и авторизации клиента и сервера.
*/
func GetCertsAndKeys(org string) (certtificate, publicKey, private []byte, err error) {
	const (
		// Длина приватного ключа
		lenPrevateKey = 4096
		// Номер сертификата
		numCert = 2674
	)

	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(numCert),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{org},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — год
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ
	privateKey, err := rsa.GenerateKey(rand.Reader, lenPrevateKey)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("encode certPEM: %w", err)
	}
	// Закрытый ключ
	var privateKeyPEM bytes.Buffer

	if err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("encode privateKeyPEM: %w", err)
	}

	publicKeyDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("x509.MarshalPKIXPublicKey: %w", err)
	}

	// Открытый ключ
	var publicKeyPEM bytes.Buffer

	if err := pem.Encode(&publicKeyPEM, &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("encode publicKeyPEM: %w", err)
	}

	return certPEM.Bytes(), publicKeyPEM.Bytes(), privateKeyPEM.Bytes(), nil
}
