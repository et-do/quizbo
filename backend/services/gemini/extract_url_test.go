package gemini

import (
	"context"
	"os"
	"testing"
)

const testHTML string = `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>`

func TestExtractContentFromHtml(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	geminiClient, err := NewGeminiClient(ctx)
	if err != nil {
		t.Fatalf("NewGeminiClient: expected no error, got %v", err)
	}

	contentMap, fullHTML, err := geminiClient.ExtractContentFromHtml(ctx, testHTML)
	if err != nil {
		t.Fatalf("ExtractContent: expected no error, got %v", err)
	}

	if contentMap["content"] == "" {
		t.Errorf("ExtractContent: expected extracted content, got an empty string")
	} else {
		t.Logf("Extracted Content: %s", contentMap["content"])
	}

	if contentMap["title"] == "" {
		t.Errorf("ExtractContent: expected generated title, got an empty string")
	} else {
		t.Logf("Generated Title: %s", contentMap["title"])
	}

	t.Logf("Full HTML Response: %s", fullHTML)
}
