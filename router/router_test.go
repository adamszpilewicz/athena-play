package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketAndKeyFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		bucket         string
		key            string
		expectedErr    error
		expectedBucket string
		expectedKey    string
	}{
		{
			name:           "Both bucket and key are present",
			bucket:         "my-bucket",
			key:            "my-key",
			expectedErr:    nil,
			expectedBucket: "my-bucket",
			expectedKey:    "my-key",
		},
		{
			name:           "Missing bucket",
			bucket:         "",
			key:            "my-key",
			expectedErr:    errors.New("missing required parameters"),
			expectedBucket: "",
			expectedKey:    "",
		},
		{
			name:           "Missing key",
			bucket:         "my-bucket",
			key:            "",
			expectedErr:    errors.New("missing required parameters"),
			expectedBucket: "",
			expectedKey:    "",
		},
		{
			name:           "Both bucket and key are missing",
			bucket:         "",
			key:            "",
			expectedErr:    errors.New("missing required parameters"),
			expectedBucket: "",
			expectedKey:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			q := req.URL.Query()
			q.Add("bucket", tt.bucket)
			q.Add("key", tt.key)
			req.URL.RawQuery = q.Encode()

			bucket, key, err := getBucketAndKeyFromRequest(req)
			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.expectedBucket, bucket)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}
