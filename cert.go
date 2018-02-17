package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

func FromPemFile(filename string, password string) (tls.Certificate, error) {

	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}

	return decodePem(tls.Certificate{}, pemBytes, password)
}

func FromPemBytes(cert []byte, key []byte) (tls.Certificate, error) {

	certi, err := tls.X509KeyPair(cert, key)

	if err != nil {
		log.Println("CERTIFICATE ERROR:",err)

		return certi , err;
	}

	return certi , nil;

}

func decodePem(cert tls.Certificate, bytes []byte, password string) (tls.Certificate, error) {
	block, rest := pem.Decode(bytes)

	log.Println("type of the block", block , rest)

	if block == nil {
		return cert, nil
	}
	if x509.IsEncryptedPEMBlock(block) {
		_, err := x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return cert, errors.New("Error decrypting certificate")
		}
	}
	switch block.Type {
	case "CERTIFICATE":
		cert.Certificate = append(cert.Certificate, block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Println("PRIVATE KEY ERROR ", err)

			return cert, errors.New("Error parsing RSA PRIVATE KEY")
		}
		cert.PrivateKey = key
	default:
		return cert, errors.New("Cert block wasn't CERTIFICATE or PRIVATE KEY")
	}
	return decodePem(cert, rest, password)
}
