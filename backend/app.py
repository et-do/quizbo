from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

# Replace with your Gemini API details
GEMINI_API_URL = 'https://api.gemini.com/v1/endpoint'
GEMINI_API_KEY = 'YOUR_GEMINI_API_KEY'
GEMINI_API_SECRET = 'YOUR_GEMINI_API_SECRET'

@app.route('/generate-questions', methods=['POST'])
def generate_questions():
    data = request.get_json()
    text = data['text']
    
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {GEMINI_API_KEY}'
    }
    
    payload = {
        'prompt': f"Generate five comprehension questions for the following text:\n\n{text}",
        'max_tokens': 150
    }
    
    response = requests.post(GEMINI_API_URL, headers=headers, json=payload)
    
    if response.status_code == 200:
        questions = response.json().get('choices', [])[0].get('text', '').strip().split('\n')
        return jsonify({'questions': questions})
    else:
        return jsonify({'error': 'Failed to generate questions'}), response.status_code

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
