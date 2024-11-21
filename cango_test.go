package cango

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCango(t *testing.T) {
}

// TestExtractRoutePaths tests the extractRoutePaths function with various input types
func TestExtractRoutePaths(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		wantPrefixes []string
		wantIsCango  bool
		wantErr      bool
	}{
		{
			name: "Basic struct with URI field",
			input: struct {
				URI `value:"/test"`
			}{},
			wantPrefixes: []string{"/test"},
			wantIsCango:  true,
		},
		{
			name: "Struct with multiple paths",
			input: struct {
				URI `value:"/test;/test2"`
			}{},
			wantPrefixes: []string{"/test", "/test2"},
			wantIsCango:  true,
		},
		{
			name: "Struct with parameter in path",
			input: struct {
				URI `value:"/users/{id}"`
			}{},
			wantPrefixes: []string{"/users/:id"},
			wantIsCango:  true,
		},
		{
			name: "Struct without URI field",
			input: struct {
				Field string
			}{},
			wantPrefixes: []string{""},
			wantIsCango:  false,
		},
		{
			name: "Pointer to struct",
			input: &struct {
				URI `value:"/pointer-test"`
			}{},
			wantPrefixes: []string{"/pointer-test"},
			wantIsCango:  true,
		},
		{
			name: "Empty path value",
			input: struct {
				URI `value:""`
			}{},
			wantPrefixes: []string{""},
			wantIsCango:  true,
		},
		{
			name: "Path with special characters",
			input: struct {
				URI `value:"/api/v1/{id}/details"`
			}{},
			wantPrefixes: []string{"/api/v1/:id/details"},
			wantIsCango:  true,
		},
		{
			name: "Multiple parameters in path",
			input: struct {
				URI `value:"/users/{userId}/posts/{postId}"`
			}{},
			wantPrefixes: []string{"/users/:userId/posts/:postId"},
			wantIsCango:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefixes, gotIsCango := extractRoutePaths(reflect.TypeOf(tt.input))

			// Check if prefixes match
			if !reflect.DeepEqual(gotPrefixes, tt.wantPrefixes) {
				t.Errorf("extractRoutePaths() prefixes = %v, want %v", gotPrefixes, tt.wantPrefixes)
			}

			// Check if isCango matches
			if gotIsCango != tt.wantIsCango {
				t.Errorf("extractRoutePaths() isCango = %v, want %v", gotIsCango, tt.wantIsCango)
			}
		})
	}
}

// TestExtractRoutePathsEdgeCases tests edge cases and potential error conditions
func TestExtractRoutePathsEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		shouldPanic bool
	}{
		{
			name:        "Nil interface",
			input:       nil,
			shouldPanic: false,
		},
		{
			name:        "Non-struct type",
			input:       "string",
			shouldPanic: false,
		},
		{
			name: "Embedded struct",
			input: struct {
				Embedded struct {
					URI `value:"/embedded"`
				}
			}{},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.shouldPanic {
					t.Errorf("extractRoutePaths() panic = %v, shouldPanic = %v", r != nil, tt.shouldPanic)
				}
			}()

			prefixes, isCango := extractRoutePaths(reflect.TypeOf(tt.input))
			// For non-panic cases, just ensure the function returns
			if !tt.shouldPanic {
				t.Logf("extractRoutePaths() returned prefixes=%v, isCango=%v", prefixes, isCango)
			}
		})
	}
}

// TestExtractRoutePathsWithCustomTypes tests the function with custom types implementing URI
type CustomURI struct{}

func (c CustomURI) ginContext() *gin.Context { return nil }
func (c CustomURI) Context() *gin.Context    { return nil }

func TestExtractRoutePathsWithCustomTypes(t *testing.T) {
	type CustomType struct {
		CustomURI `value:"/custom"`
	}

	prefixes, isCango := extractRoutePaths(reflect.TypeOf(CustomType{}))

	if !isCango {
		t.Error("Expected CustomType to be recognized as Cango type")
	}

	expectedPrefix := []string{"/custom"}
	if !reflect.DeepEqual(prefixes, expectedPrefix) {
		t.Errorf("Expected prefix %v, got %v", expectedPrefix, prefixes)
	}
}
