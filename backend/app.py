from flask import Flask, request, jsonify
from backend.llms.utils import create_model, get_response_text_from_model, Role
import vertexai
import os

app = Flask(__name__)

vertexai.init(project=os.environ["GCP_PROJECT"])

@app.route("/process-page", methods=["POST"])
def process_page():
    data = request.get_json()
    url = data.get("url")
    html_content = data.get("html")

    # Web scraping
    scraper_model = create_model(role=Role.WEBSCRAPER)
    response = scraper_model.generate_content([f"HTML Content: {html_content} || Your Response: "], stream=False)
    scraped_content = response.text

    # Generating questions and answers
    question_creator_model = create_model(role=Role.QUESTION_CREATOR)
    qa_response = question_creator_model.generate_content([f"Generate ten comprehension questions and answers for the following text:\n\n{scraped_content}"], stream=False)
    questions_answers = qa_response.text

    return jsonify({
        "url": url,
        "scraped_content": scraped_content,
        "questions_answers": questions_answers
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)