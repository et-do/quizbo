chrome.runtime.onInstalled.addListener(() => {
  chrome.identity.getAuthToken({ interactive: true }, (token) => {
    if (chrome.runtime.lastError || !token) {
      console.error('Failed to get auth token:', chrome.runtime.lastError);
      return;
    }
    console.log('Auth token:', token);

    // Get user info
    fetch('https://www.googleapis.com/oauth2/v3/userinfo', {
      headers: { Authorization: `Bearer ${token}` }
    })
    .then(response => response.json())
    .then(userInfo => {
      console.log('User info:', userInfo);
      // Store user info in local storage
      chrome.storage.local.set({ userInfo });
    })
    .catch(error => console.error('Failed to fetch user info:', error));
  });
});
