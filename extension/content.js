(function() {
  if (document.getElementById('readrobin-sidebar')) {
    return;
  }

  const sidebar = document.createElement('div');
  sidebar.id = 'readrobin-sidebar';
  sidebar.style.position = 'fixed';
  sidebar.style.right = '0';
  sidebar.style.top = '0';
  sidebar.style.width = '300px';
  sidebar.style.height = '100%';
  sidebar.style.backgroundColor = '#f1f1f1';
  sidebar.style.borderLeft = '1px solid #ccc';
  sidebar.style.zIndex = '10000';
  sidebar.style.overflowY = 'auto';
  sidebar.innerHTML = `
    <div style="padding: 10px;">
      <h2>ReadRobin Sidebar</h2>
      <p>Content goes here...</p>
      <button id="close-sidebar" style="position: absolute; top: 10px; right: 10px;">Close</button>
    </div>
  `;

  document.body.appendChild(sidebar);

  document.getElementById('close-sidebar').addEventListener('click', () => {
    document.getElementById('readrobin-sidebar').remove();
  });
})();
