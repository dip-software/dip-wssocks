package tools

import (
	"testing"

	"github.com/loafoe/caddy-token/keys"
)

func TestGetGatewayFromAPIKey(t *testing.T) {
	// Shared secret for testing
	sharedSecret := "test-secret"

	// Test cases for different API key scenarios
	testCases := []struct {
		name          string
		scopes        []string
		expectedGW    string
		expectErr     bool
		invalidFormat bool
	}{
		{
			name:       "Valid API key with gateway scope",
			scopes:     []string{"gw:test-gateway", "other-scope"},
			expectedGW: "wss://test-gateway",
			expectErr:  false,
		},
		{
			name:       "Valid API key with multiple gateway scopes (should use first)",
			scopes:     []string{"other-scope", "gw:primary-gateway", "gw:secondary-gateway"},
			expectedGW: "wss://primary-gateway",
			expectErr:  false,
		},
		{
			name:       "Valid API key without gateway scope",
			scopes:     []string{"scope1", "scope2"},
			expectedGW: "",
			expectErr:  false,
		},
		{
			name:          "Invalid API key format",
			scopes:        []string{},
			expectedGW:    "",
			expectErr:     true,
			invalidFormat: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiKey string
			var err error

			if tc.invalidFormat {
				// Use an invalid API key format
				apiKey = "invalid-format"
			} else {
				// Generate a valid API key with the test scopes
				// Using version "2" as required by the caddy-token/keys package
				apiKey, _, err = keys.GenerateDeterministicAPIKey("2", sharedSecret,
					keys.WithScopes(tc.scopes),
					keys.WithToken("test-token"),
					keys.WithOrganization("test-org"),
					keys.WithRegion("test-region"))
				if err != nil {
					t.Fatalf("Failed to generate test API key: %v", err)
				}

				// Print the API key format for debugging only in verbose mode
				if testing.Verbose() {
					t.Logf("Generated API key: %s", apiKey)
				}
			}

			// Call the function under test
			gateway, err := GetGatewayFromAPIKey(apiKey)

			// Check error expectations
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check gateway expectations
			if gateway != tc.expectedGW {
				t.Errorf("Expected gateway '%s', got '%s'", tc.expectedGW, gateway)
			}
		})
	}
}

func TestGetEndpointFromAPIKey(t *testing.T) {
	// Shared secret for testing
	sharedSecret := "test-secret"

	// Test cases for different API key scenarios
	testCases := []struct {
		name          string
		scopes        []string
		expectedEP    string
		expectErr     bool
		invalidFormat bool
	}{
		{
			name:       "Valid API key with endpoint scope",
			scopes:     []string{"ep:test-endpoint", "other-scope"},
			expectedEP: "test-endpoint",
			expectErr:  false,
		},
		{
			name:       "Valid API key with multiple endpoint scopes (should use first)",
			scopes:     []string{"other-scope", "ep:primary-endpoint", "ep:secondary-endpoint"},
			expectedEP: "primary-endpoint",
			expectErr:  false,
		},
		{
			name:       "Valid API key without endpoint scope",
			scopes:     []string{"scope1", "scope2"},
			expectedEP: "",
			expectErr:  false,
		},
		{
			name:          "Invalid API key format",
			scopes:        []string{},
			expectedEP:    "",
			expectErr:     true,
			invalidFormat: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiKey string
			var err error

			if tc.invalidFormat {
				// Use an invalid API key format
				apiKey = "invalid-format"
			} else {
				// Generate a valid API key with the test scopes
				// Using version "2" as required by the caddy-token/keys package
				apiKey, _, err = keys.GenerateDeterministicAPIKey("2", sharedSecret,
					keys.WithScopes(tc.scopes),
					keys.WithToken("test-token"),
					keys.WithOrganization("test-org"),
					keys.WithRegion("test-region"))
				if err != nil {
					t.Fatalf("Failed to generate test API key: %v", err)
				}

				// Print the API key format for debugging only in verbose mode
				if testing.Verbose() {
					t.Logf("Generated API key: %s", apiKey)
				}
			}

			// Call the function under test
			endpoint, err := GetEndpointFromAPIKey(apiKey)

			// Check error expectations
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check endpoint expectations
			if endpoint != tc.expectedEP {
				t.Errorf("Expected endpoint '%s', got '%s'", tc.expectedEP, endpoint)
			}
		})
	}
}

func TestGetAllEndpointsFromAPIKey(t *testing.T) {
	// Shared secret for testing
	sharedSecret := "test-secret"

	// Test cases for different API key scenarios
	testCases := []struct {
		name          string
		scopes        []string
		expectedEPs   []string
		expectErr     bool
		invalidFormat bool
	}{
		{
			name:        "Valid API key with single endpoint scope",
			scopes:      []string{"ep:test-endpoint", "other-scope"},
			expectedEPs: []string{"test-endpoint"},
			expectErr:   false,
		},
		{
			name:        "Valid API key with multiple endpoint scopes",
			scopes:      []string{"other-scope", "ep:endpoint1", "ep:endpoint2", "gw:gateway"},
			expectedEPs: []string{"endpoint1", "endpoint2"},
			expectErr:   false,
		},
		{
			name:        "Valid API key without endpoint scope",
			scopes:      []string{"scope1", "gw:gateway", "scope2"},
			expectedEPs: nil,
			expectErr:   false,
		},
		{
			name:        "Valid API key with only prefix-like scopes",
			scopes:      []string{"ep", "endpoint:test", "epx:test"},
			expectedEPs: nil,
			expectErr:   false,
		},
		{
			name:          "Invalid API key format",
			scopes:        []string{},
			expectedEPs:   nil,
			expectErr:     true,
			invalidFormat: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiKey string
			var err error

			if tc.invalidFormat {
				// Use an invalid API key format
				apiKey = "invalid-format"
			} else {
				// Generate a valid API key with the test scopes
				// Using version "2" as required by the caddy-token/keys package
				apiKey, _, err = keys.GenerateDeterministicAPIKey("2", sharedSecret,
					keys.WithScopes(tc.scopes),
					keys.WithToken("test-token"),
					keys.WithOrganization("test-org"),
					keys.WithRegion("test-region"))
				if err != nil {
					t.Fatalf("Failed to generate test API key: %v", err)
				}

				// Print the API key format for debugging only in verbose mode
				if testing.Verbose() {
					t.Logf("Generated API key: %s", apiKey)
				}
			}

			// Call the function under test
			endpoints, err := GetAllEndpointsFromAPIKey(apiKey)

			// Check error expectations
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check endpoints expectations
			if len(endpoints) != len(tc.expectedEPs) {
				t.Errorf("Expected %d endpoints, got %d", len(tc.expectedEPs), len(endpoints))
			} else {
				// Check each endpoint
				for i, expected := range tc.expectedEPs {
					if i >= len(endpoints) || endpoints[i] != expected {
						t.Errorf("Expected endpoint[%d] '%s', got '%s'", i, expected, endpoints[i])
					}
				}
			}
		})
	}
}
