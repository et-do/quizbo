chrome.runtime.onInstalled.addListener(() => {
    console.log('Article Ally extension installed');
  });
  
  chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.type === 'fetchComprehensionQuestions') {
      fetch('http://localhost:5000/generate-questions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ text: message.text })
      })
      .then(response => response.json())
      .then(data => {
        sendResponse({ questions: data.questions });
      })
      .catch(error => {
        console.error('Error:', error);
        sendResponse({ error: 'Failed to fetch questions' });
      });
      return true; // Will respond asynchronously
    }
  });
  