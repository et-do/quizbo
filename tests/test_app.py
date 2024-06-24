import pytest
import json
from backend.app import app  # Adjust the import according to your app structure

@pytest.fixture
def client():
    app.config['TESTING'] = True
    client = app.test_client()

    yield client

def test_process_page(client):
    # Define the payload
    payload = {
        "url": "http://example.com",
        "html": "<html><body><h1>Example Article</h1><p>This is an example.</p></body></html>"
    }

    # Send POST request
    response = client.post("/process-page", data=json.dumps(payload), content_type='application/json')

    # Check the response
    assert response.status_code == 200
    data = response.get_json()
    assert 'url' in data
    assert 'scraped_content' in data
    assert 'questions_answers' in data

    # Print the output for verification
    print(f"URL: {data['url']}")
    print(f"Scraped Content: {data['scraped_content']}")
    print(f"Questions and Answers: {data['questions_answers']}")
