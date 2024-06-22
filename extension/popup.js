document.addEventListener('DOMContentLoaded', () => {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      chrome.tabs.sendMessage(tabs[0].id, { type: 'fetchComprehensionQuestions' }, (response) => {
        if (response.error) {
          document.getElementById('questions').innerText = response.error;
        } else {
          const questions = response.questions;
          const questionsDiv = document.getElementById('questions');
          questionsDiv.innerHTML = '';
          questions.forEach((question, index) => {
            const questionElement = document.createElement('p');
            questionElement.innerText = `${index + 1}. ${question}`;
            questionsDiv.appendChild(questionElement);
          });
        }
      });
    });
  });
  