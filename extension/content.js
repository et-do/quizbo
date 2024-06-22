console.log('Article Ally content script loaded');

// Example: Sending the entire article text to the background script
const articleText = document.body.innerText;

chrome.runtime.sendMessage({ type: 'fetchComprehensionQuestions', text: articleText }, (response) => {
  if (response.error) {
    console.error(response.error);
  } else {
    console.log('Comprehension questions:', response.questions);
    // You can now display the questions in your extension's popup or directly on the page
  }
});
