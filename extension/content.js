(function() {
  if (document.getElementById('readrobin-sidebar')) {
    return;
  }

  const sidebar = document.createElement('div');
  sidebar.id = 'readrobin-sidebar';
  sidebar.style.position = 'fixed';
  sidebar.style.right = '0';
  sidebar.style.top = '0';
  sidebar.style.width = '400px'; // Increased width
  sidebar.style.height = '100%';
  sidebar.style.backgroundColor = '#f9f9f9'; // Softer background color
  sidebar.style.borderLeft = '1px solid #ccc';
  sidebar.style.zIndex = '10000';
  sidebar.style.overflowY = 'auto';
  sidebar.style.fontFamily = 'Arial, sans-serif'; // Added font
  sidebar.innerHTML = `
    <div style="position: relative; padding: 20px;">
      <button id="close-sidebar" style="position: absolute; top: 10px; right: 10px; background: none; border: none; font-size: 20px;">&times;</button>
      <h2 style="margin-top: 40px; color: #333;">Welcome to ReadRobin</h2>
      <h3 style="color: #666;">Lets test your comprehension of: Article Title</h3>
      <div id="questions">
        ${Array.from({ length: 10 }).map((_, i) => `
          <div style="margin-bottom: 20px;">
            <label for="question-${i + 1}" style="font-weight: bold; color: #333;">Question ${i + 1}:</label>
            <input type="text" id="question-${i + 1}" placeholder="Question ${i + 1}" style="width: calc(100% - 20px); padding: 5px; margin-bottom: 10px; border: 1px solid #ccc; border-radius: 4px; background-color: #fff;" disabled />
            <textarea id="answer-${i + 1}" placeholder="Your answer here..." style="width: calc(100% - 20px); height: 60px; padding: 5px; margin-bottom: 10px; border: 1px solid #ccc; border-radius: 4px; background-color: #fff;"></textarea>
            <button id="submit-${i + 1}" style="width: 100%; padding: 10px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">Submit</button>
          </div>
        `).join('')}
      </div>
    </div>
  `;

  document.body.appendChild(sidebar);

  document.getElementById('close-sidebar').addEventListener('click', () => {
    document.getElementById('readrobin-sidebar').remove();
  });
})();
