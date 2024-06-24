from enum import Enum
from typing import List


class Role(Enum):
    WEBSCRAPER = "WebScraper"
    QUESTION_CREATOR = "QuestionCreator"
    ANSWER_REVIEWER = "AnswerReviewer"


WEBSCRAPER_INSTRUCTIONS = [
    """
    You are an expert Webscraper. Scrape this HTML content for the article name and its contents.
    Format your response as follows:
    {
        "article_name": "",
        "content": ""
    }
    """
]

QUESTION_CREATOR_INSTRUCTIONS = [
    """
    You are an expert in generating comprehension questions. Create 3 comprehension questions and their answers for the following text.
    Format your response as follows:
    {
        "question_1":
            "question": "",
            "answer": "",
        "question_2": "",
            "question": "",
            "answer": "",
        "question_3": "",
            "question": "",
            "answer": "",
    }
    """
]

ANSWER_REVIEWER_INSTRUCTIONS = [
    """
    You are an expert in evaluating answers. Compare the following user responses with the correct answers and provide feedback on accuracy and completeness.
    """
]

SYSTEM_INSTRUCTIONS = {
    Role.WEBSCRAPER: WEBSCRAPER_INSTRUCTIONS,
    Role.QUESTION_CREATOR: QUESTION_CREATOR_INSTRUCTIONS,
    Role.ANSWER_REVIEWER: ANSWER_REVIEWER_INSTRUCTIONS,
}


def get_system_instructions(role: Role) -> List[str]:
    return SYSTEM_INSTRUCTIONS.get(role, [])