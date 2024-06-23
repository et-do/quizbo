chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'sendPageInfo') {
    const { url, html } = message;

    // Replace with your backend URL
    const backendUrl = 'http://localhost:5000/process-page';

    fetch(backendUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        url: url,
        html: html,
      }),
    })
    .then(response => response.json())
    .then(data => {
      console.log('Page processed:', data);
      sendResponse({ status: 'success', data: data });
    })
    .catch(error => {
      console.error('Error processing page:', error);
      sendResponse({ status: 'error', error: error });
    });

    // Keep the message channel open until sendResponse is called
    return true;
  }
});
