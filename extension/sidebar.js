document.addEventListener('DOMContentLoaded', () => {
    console.log('Sidebar loaded');
  
    chrome.storage.local.get(['authToken', 'userInfo'], (result) => {
      if (result.authToken && result.userInfo) {
        const authToken = result.authToken;
        const userId = result.userInfo.sub; // Using the 'sub' field as user ID
        console.log('User ID:', userId);
  
        // Example of sending a request with the user ID
        fetch('http://localhost:5000/generate-questions', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${authToken}`
          },
          body: JSON.stringify({
            text: 'Your text content here',
            user_id: userId
          })
        })
        .then(response => response.json())
        .then(data => {
          console.log('Questions:', data.questions);
          document.getElementById('content').innerHTML += `<p>User ID: ${userId}</p>`;
          document.getElementById('content').innerHTML += `<ul>${data.questions.map(q => `<li>${q}</li>`).join('')}</ul>`;
        })
        .catch(error => console.error('Error:', error));
      } else {
        console.log('No auth token or user info found');
      }
    });
  
    // Add event listener for the close button
    document.getElementById('close-sidebar').addEventListener('click', () => {
      document.getElementById('readrobin-sidebar').remove();
    });
  });
  