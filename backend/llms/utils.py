from vertexai.generative_models import GenerativeModel
import vertexai.preview.generative_models as generative_models
from typing import Optional, Any
from enum import Enum


class Role(Enum):
    WEBSCRAPER = "WebScraper"
    QUESTION_CREATOR = "QuestionCreator"
    ANSWER_REVIEWER = "AnswerReviewer"


def get_response_text_from_model(
    model: Any,
    prompt: str,
    generation_config: Optional[dict] = None,
    safety_settings: Optional[dict] = None,
    stream=False,
) -> Any:

    if generation_config is None and safety_settings is None:
        response = model.generate_content(
            [prompt],
            stream=stream,
        )
        return response.text

    response = model.generate_content(
        [prompt],
        generation_config=generation_config,
        safety_settings=safety_settings,
        stream=stream,
    )

    return response.text


def create_model(model="gemini-1.5-flash-001", role: Optional[Role] = None):
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

    if role is None:
        role = Role.WEBSCRAPER

    model = GenerativeModel(
        model_name=model,
        generation_config=generation_config,
        safety_settings=safety_settings,
        system_instruction=[
            """You are an expert Webscraper. Scrape this HTML content for the Article name and its contents."""
        ],
    )

    return model
