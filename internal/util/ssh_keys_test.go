// Copyright (C) 2024
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.

package util

import (
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateSSHKeyPairInMemory(t *testing.T) {
	keyPair, err := GenerateSSHKeyPairInMemory()
	if err != nil {
		t.Fatalf("Failed to generate SSH key pair: %v", err)
	}

	// Verify private key is not empty
	if keyPair.PrivateKey == "" {
		t.Error("Private key is empty")
	}

	// Verify public key is not empty
	if keyPair.PublicKey == "" {
		t.Error("Public key is empty")
	}

	// Verify private key has correct format
	if !strings.Contains(keyPair.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----") ||
		!strings.Contains(keyPair.PrivateKey, "-----END RSA PRIVATE KEY-----") {
		t.Error("Private key doesn't have correct PEM format")
	}

	// Verify public key has correct format
	if !strings.HasPrefix(keyPair.PublicKey, "ssh-rsa ") {
		t.Error("Public key doesn't have correct SSH format")
	}

	// Verify public key can be parsed by SSH library
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(keyPair.PublicKey))
	if err != nil {
		t.Errorf("Generated public key is not valid SSH key: %v", err)
	}

	// Test multiple generations produce different keys
	keyPair2, err := GenerateSSHKeyPairInMemory()
	if err != nil {
		t.Fatalf("Failed to generate second SSH key pair: %v", err)
	}

	if keyPair.PrivateKey == keyPair2.PrivateKey {
		t.Error("Two key generations produced identical private keys")
	}

	if keyPair.PublicKey == keyPair2.PublicKey {
		t.Error("Two key generations produced identical public keys")
	}
}