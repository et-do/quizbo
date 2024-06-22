document.getElementById('open-sidebar').addEventListener('click', () => {
  chrome.identity.getAuthToken({ interactive: true }, (token) => {
    if (chrome.runtime.lastError || !token) {
      console.error('Failed to get auth token:', chrome.runtime.lastError);
      alert('Authentication failed');
      return;
    }
    console.log('Auth token:', token);

    // Store the auth token in local storage
    chrome.storage.local.set({ authToken: token }, () => {
      console.log('Auth token stored');
      chrome.sidebarAction.open();
    });
  });
});
