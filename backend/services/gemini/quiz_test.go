package gemini

import (
	"context"
	"testing"

	"read-robin/models"

	"github.com/stretchr/testify/assert"
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

func TestGenerateQuiz(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	summarizedContent := "Test summarized content"
	personaName := "Test Persona"
	personaRole := "Test Role"
	personaLanguage := "English"
	personaDifficulty := "Easy"

	quizContent, response, err := client.GenerateQuiz(ctx, summarizedContent, personaName, personaRole, personaLanguage, personaDifficulty)
	assert.NoError(t, err)
	assert.NotEmpty(t, quizContent)
	assert.NotEmpty(t, response)
}

func TestExtractAndGenerateQuiz(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	htmlContent := testHTML
	persona := models.Persona{
		Name:       "Test Persona",
		Role:       "Test Role",
		Language:   "English",
		Difficulty: "Easy",
	}

	quizContentMap, title, err := client.ExtractAndGenerateQuiz(ctx, htmlContent, persona)
	assert.NoError(t, err)
	assert.NotNil(t, quizContentMap)
	assert.NotEmpty(t, title)
}

func TestReviewResponse(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	reviewData := "Test review data"

	reviewResult, err := client.ReviewResponse(ctx, reviewData)
	assert.NoError(t, err)
	assert.NotEmpty(t, reviewResult)
}
