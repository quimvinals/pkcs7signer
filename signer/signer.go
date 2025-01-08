package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/unidoc/pkcs7"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <file-to-sign> <certificate-file> <private-key-file> <output-signature-file>\n", os.Args[0])
		os.Exit(1)
	}

	fileToSign := os.Args[1]
	certFile := os.Args[2]
	keyFile := os.Args[3]
	outputSigFile := os.Args[4]

	// Read the data to sign
	data, err := os.ReadFile(fileToSign)
	if err != nil {
		log.Fatalf("Failed to read file to sign: %v", err)
	}

	wwdrCertFile := "./wwdr.pem"
	wwdrCertPEM, err := os.ReadFile(wwdrCertFile)
	if err != nil {
		log.Fatalf("Failed to read WWDR certificate file: %v", err)
	}
	wwdrCertBlock, _ := pem.Decode(wwdrCertPEM)
	if wwdrCertBlock == nil || wwdrCertBlock.Type != "CERTIFICATE" {
		log.Fatalf("Failed to decode WWDR certificate PEM")
	}
	wwdrCertificate, err := x509.ParseCertificate(wwdrCertBlock.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse WWDR certificate: %v", err)
	}

	// Load the certificate
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		log.Fatalf("Failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}

	// Load the private key
	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		log.Fatalf("Failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse private key: %v", err)
		}
	}

	// Create a PKCS7 signer
	signedData, err := pkcs7.NewSignedData(data)
	if err != nil {
		log.Fatalf("Failed to create signed data: %v", err)
	}

	signedData.SetEncryptionAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)
	signedData.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)

	// Add certificates
	signedData.AddCertificate(wwdrCertificate)

	err = signedData.AddSigner(cert, privateKey, pkcs7.SignerInfoConfig{})

	if err != nil {
		log.Fatalf("Failed to add signer: %v", err)
	}

	signedData.Detach() // Create a detached signature

	// Generate the signature
	signature, err := signedData.Finish()
	if err != nil {
		log.Fatalf("Failed to generate signature: %v", err)
	}

	// Write the signature to the output file
	err = os.WriteFile(outputSigFile, signature, 0644)
	if err != nil {
		log.Fatalf("Failed to write signature file: %v", err)
	}

	fmt.Println("Signature generated successfully with SHA-256.")
}
