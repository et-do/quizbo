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
    prompt = data.get("html")
    print(url)

    model = create_model(role=Role.WEBSCRAPER)

    response_text = get_response_text_from_model(prompt=prompt, model=model)

    print(response_text)

    # Generate questions and answers
    questions_answers = model.predict(
        f"Generate ten comprehension questions and answers for the following text:\n\n{extracted_text}",
        max_output_tokens=1024,
    )

    return jsonify(
        {
            "url": url,
            "summary": extracted_text,
            "questions_answers": questions_answers.text,
        }
    )


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
