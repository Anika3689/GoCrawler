package main 
import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
        	{
			name:	  "strip trailing slash",
			inputURL: "http://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{	
			name:	  "strip fragment",
			inputURL: "https://www.example.com/users?id=42#profile",
			expected: "www.example.com/users?id=42",
		},
		{	
			name: "remove duplicate slashes",
			inputURL: "http://docs.example.com/a//b///c",
			expected: "docs.example.com/a/b/c",
		},
		{	name: "lowercase host",
			inputURL: "https://WWW.ExAmPlE.CoM/users?id=42",
			expected: "www.example.com/users?id=42",
		},
		{
			name: "remove default ports",
			inputURL: "http://www.example.com:80/users",
			expected: "www.example.com/users",
		},
		{
                        name: "combined",
                        inputURL: "http://API.EXAMPLE.COM:80///v1//search?q=go#top",
                        expected: "api.example.com/v1/search?q=go",
 		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s (%s) FAIL: expected URL: %v, actual: %v", i, tc.name, tc.inputURL, tc.expected, actual)
			}
		})
	}
}
