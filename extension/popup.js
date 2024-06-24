document.addEventListener('DOMContentLoaded', () => {
  const loginButton = document.getElementById('login');

  if (loginButton) {
    loginButton.addEventListener('click', () => {
      const userId = 'user_' + Math.random().toString(36).substr(2, 9);
      chrome.storage.local.set({ userId }, () => {
        console.log('User ID generated and saved:', userId);
        alert('Logged in with User ID: ' + userId);
      });

      chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
        const activeTab = tabs[0];
        if (!activeTab.url.startsWith('chrome://')) {
          chrome.scripting.executeScript({
            target: { tabId: activeTab.id },
            files: ['content.js']
          }, () => {
            chrome.scripting.insertCSS({
              target: { tabId: activeTab.id },
              files: ['sidebar.css']
            });
          });
        } else {
          console.error('Cannot run script on chrome:// URLs');
        }
      });
    });
  } else {
    console.error('Login button not found');
  }
});
