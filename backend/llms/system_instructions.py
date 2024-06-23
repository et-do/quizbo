from enum import Enum
from typing import List


class Role(Enum):
    WEBSCRAPER = "WebScraper"
    QUESTION_CREATOR = "QuestionCreator"
    ANSWER_REVIEWER = "AnswerReviewer"


WEBSCRAPER_INSTRUCTIONS = [
    """
    You are an expert Webscraper. Scrape this HTML content for the Article name and its contents.
    """
]

QUESTION_CREATOR_INSTRUCTIONS = [
    """
    You are an expert in generating comprehension questions. Create ten comprehension questions for the following text. Ensure the questions vary in difficulty and cover different aspects of the text.
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


# Example usage:
role = Role.QUESTION_CREATOR
instructions = get_system_instructions(role)
print(instructions)
