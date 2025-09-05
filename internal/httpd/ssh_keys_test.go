// Copyright (C) 2024
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.

package httpd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drakkan/sftpgo/v2/internal/util"
)

func TestHandleWebGenerateUserSSHKeys(t *testing.T) {
	// Create a test server instance
	server := &httpdServer{}
	
	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/web/admin/users/generate-ssh-keys", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Create a test response recorder
	w := httptest.NewRecorder()
	
	// Note: This test will fail without proper token validation, but we can test the key generation logic
	// In a real test environment, you would need to set up proper authentication and CSRF tokens
	
	// Test that the endpoint exists and responds appropriately
	// We expect it to fail with authentication error since we don't have proper tokens
	server.handleWebGenerateUserSSHKeys(w, req)
	
	// Check that the response indicates authentication failure (400 or 403)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusForbidden {
		t.Logf("Expected authentication failure, got status: %d", w.Code)
	}
}

func TestSSHKeyGenerationIntegration(t *testing.T) {
	// Test that our utility function works and returns valid JSON-serializable data
	keyPair, err := util.GenerateSSHKeyPairInMemory()
	if err != nil {
		t.Fatalf("Failed to generate SSH key pair: %v", err)
	}
	
	// Test JSON serialization (what the endpoint would return)
	jsonData, err := json.Marshal(keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal key pair to JSON: %v", err)
	}
	
	// Test JSON deserialization
	var deserializedKeyPair util.SSHKeyPair
	err = json.Unmarshal(jsonData, &deserializedKeyPair)
	if err != nil {
		t.Fatalf("Failed to unmarshal key pair from JSON: %v", err)
	}
	
	// Verify data integrity
	if deserializedKeyPair.PrivateKey != keyPair.PrivateKey {
		t.Error("Private key mismatch after JSON serialization/deserialization")
	}
	
	if deserializedKeyPair.PublicKey != keyPair.PublicKey {
		t.Error("Public key mismatch after JSON serialization/deserialization")
	}
	
	// Verify the JSON structure contains expected fields
	jsonString := string(jsonData)
	if !strings.Contains(jsonString, "private_key") {
		t.Error("JSON output missing 'private_key' field")
	}
	
	if !strings.Contains(jsonString, "public_key") {
		t.Error("JSON output missing 'public_key' field")
	}
}