from flask import Flask, request, jsonify
from flask_cors import CORS
from backend.llms.utils import create_model, Role
import vertexai
import os
import json

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})

vertexai.init(project=os.environ["GCP_PROJECT"])


@app.route("/process-page", methods=["POST"])
def process_page():
    print("request received")
    data = request.get_json()
    url = data.get("url")
    html_content = data.get("html")

    # Web scraping
    scraper_model = create_model(role=Role.WEBSCRAPER)
    response = scraper_model.generate_content(
        [f"HTML Content: {html_content} || Your Response: "], stream=False
    )
    scraped_content = response.text  # Directly get the text from the response
    print(scraped_content)
    
    # Generating questions and answers
    question_creator_model = create_model(role=Role.QUESTION_CREATOR)
    qa_response = question_creator_model.generate_content(
        [f"Content:\n\n{scraped_content}"], stream=False
    )
    qa_text = qa_response.text  # Directly get the text from the response
    print(qa_text)
    
    # Assuming the response is a JSON string
    try:
        questions_answers_dict = json.loads(qa_text)
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return jsonify({"error": "Failed to decode the questions and answers response."}), 500

    print(questions_answers_dict)

    return jsonify(
        {
            "url": url,
            "scraped_content": scraped_content,
            "questions_answers": questions_answers_dict,
        }
    )


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
