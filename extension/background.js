chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'generateQuiz') {
    const { url, html, userId } = message.data;
    console.log('Sending request to backend with:', { url, html, userId });

    fetch('http://127.0.0.1:5000/process-page', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        url: url,
        html: html,
        user_id: userId
      })
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(data => {
      console.log('Received response from backend:', data);
      sendResponse(data);
    })
    .catch(error => {
      console.error('Error:', error);
      sendResponse({ error: 'Failed to fetch data. Please try again later.' });
    });

    return true;  // Keep the message channel open for sendResponse
  }
});
