from flask import Flask, request, jsonify
import vertexai
from vertexai.generative_models import GenerativeModel
import vertexai.preview.generative_models as generative_models
import os

app = Flask(__name__)

vertexai.init(project=os.environ["GCP_PROJECT"])


@app.route("/process-page", methods=["POST"])
def process_page():
    data = request.get_json()
    url = data.get("url")
    html_content = data.get("html")
    print(url)

    generation_config = {
        "max_output_tokens": 8192,
        "temperature": 1,
        "top_p": 0.95,
    }

    safety_settings = {
        generative_models.HarmCategory.HARM_CATEGORY_HATE_SPEECH: generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
        generative_models.HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT: generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
        generative_models.HarmCategory.HARM_CATEGORY_SEXUALLY_EXPLICIT: generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
        generative_models.HarmCategory.HARM_CATEGORY_HARASSMENT: generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
    }

    # Process the HTML content with the LLM
    model = GenerativeModel(
        "gemini-1.5-flash-001",
        system_instruction=["""You are an expert Webscraper. Scrape this HTML content for the Article name and its contents."""]
    )

    response = model.generate_content(
        [f"""HTML Content: {html_content} || Your Response: """],
        generation_config=generation_config,
        safety_settings=safety_settings,
        stream=False,
    )

    # Access the text from the single response
    extracted_text = response.text
    print(extracted_text)

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
