from flask import Flask, request, jsonify
from flask_cors import CORS
from backend.llms.utils import create_model, Role, get_response_text_from_model
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

    try:
        # Web scraping
        scraper_model = create_model(role=Role.WEBSCRAPER)
        scraped_content = get_response_text_from_model(
            scraper_model,
            f"HTML Content: {html_content} || Your Response: ",
            stream=False,
        )
        print("Scraped content:", scraped_content)

        # Generating questions and answers
        question_creator_model = create_model(role=Role.QUESTION_CREATOR)
        qa_response = get_response_text_from_model(
            question_creator_model,
            f"Content:\n\n{scraped_content}",
            stream=False
        )
        print("Questions and answers response:", qa_response)

        # Assuming the response is a JSON string
        questions_answers_dict = json.loads(qa_response)
        print("Questions and answers dictionary:", questions_answers_dict)

        return jsonify(
            {
                "url": url,
                "scraped_content": scraped_content,
                "questions_answers": questions_answers_dict,
            }
        ), 200

    except Exception as e:
        print(f"Error processing request: {e}")
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
