package main
import (
	"testing"
	"net/url"
	"reflect"
	"fmt"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLUsesH2AsFallback(t *testing.T) {
	inputBody := "<html><body><h2>Fallback Title</h2></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Fallback Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLPrefersH1OverH2(t *testing.T) {
	inputBody := "<html><body><h2>Secondary Title</h2><h1>Primary Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Primary Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLReturnsEmptyStringWhenNoHeadings(t *testing.T) {
	inputBody := "<html><body><p>No headings here</p></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLReturnsEmptyStringForEmptyHTML(t *testing.T) {
	inputBody := ""
	actual := getHeadingFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLHandlesNestedContent(t *testing.T) {
	inputBody := "<html><body><h1><span>Nested</span> Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Nested Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLMultipleH1Tags(t *testing.T) {
	inputBody := "<html><body><h1>First Title</h1><h1>Second Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "First Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoMain(t *testing.T) {
	inputBody := `<html><body>
		<p>First paragraph.</p>
		<p>Second paragraph.</p>
	</body></html>`

	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoParagraph(t *testing.T) {
	inputBody := `<html><body>
		<main>
			<h1>Title</h1>
			<div>No paragraphs here.</div>
		</main>
	</body></html>`

	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLFirstParagraphInMain(t *testing.T) {
	inputBody := `<html><body>
		<main>
			<p>First main paragraph.</p>
			<p>Second main paragraph.</p>
		</main>
	</body></html>`

	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		body     string
		expected []string
		wantErr  bool
		img bool
	}{
		{
			name:     "absolute URL",
			baseURL:  "https://crawler-test.com",
			body:     `<html><body><a href="https://crawler-test.com">Boot.dev</a></body></html>`,
			expected: []string{"https://crawler-test.com"},
		},
		{
			name:     "relative URL",
			baseURL:  "https://crawler-test.com",
			body:     `<a href="/about">About</a>`,
			expected: []string{"https://crawler-test.com/about"},
		},
		{
			name:     "no links",
			baseURL:  "https://crawler-test.com",
			body:     `<html><body>Hello!</body></html>`,
			expected: []string{},
		},
		{
			name:	"relative URL - img",
			baseURL: "https://crawler-test.com",
			body: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected: []string{"https://crawler-test.com/logo.png"},
			img: true,
		},
		{
			name: 	"absolute URL - img",
			baseURL: "https://crawler-test.com",
			body: `<html><body><img src="https://cdn.example.com/image.jpg"></body></html>`,
			expected: []string{"https://cdn.example.com/image.jpg"},
			img: true,
		},
		{	
			name:	"relative path without leading slash",
			baseURL: "https://crawler-test.com/blog/",
			body: `<html><body><img src="images/photo.jpg"></body></html>`,
			expected : []string{"https://crawler-test.com/blog/images/photo.jpg"},
			img: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("couldn't parse base URL: %v", err)
			}

			var actual []string
			if tt.img {
				actual, err = getImagesFromHTML(tt.body, baseURL)
			} else {
				actual, err = getURLsFromHTML(tt.body, baseURL)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
				
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				fmt.Printf("expected type %T, got type %T", tt.expected, actual)
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}


func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:             "https://crawler-test.com",
		Heading:         "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
