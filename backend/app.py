from flask import Flask, request, jsonify
from google.auth.transport import requests
from google.oauth2 import id_token

app = Flask(__name__)

@app.route('/generate-questions', methods=['POST'])
def generate_questions():
    data = request.get_json()
    token = request.headers.get('Authorization').split('Bearer ')[1]
    try:
        # Verify the token
        id_info = id_token.verify_oauth2_token(token, requests.Request())
        user_id = id_info['sub']
        
        text = data['text']
        
        # Mocking LLM response for illustration purposes
        questions = ["Question 1", "Question 2", "Question 3", "Question 4", "Question 5"]
        
        return jsonify({'user_id': user_id, 'questions': questions})
    except ValueError:
        # Invalid token
        return jsonify({'error': 'Invalid token'}), 401

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
