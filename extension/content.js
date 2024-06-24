(function() {
  if (document.getElementById('readrobin-sidebar')) {
    return;
  }

  const sidebar = document.createElement('div');
  sidebar.id = 'readrobin-sidebar';
  sidebar.style.position = 'fixed';
  sidebar.style.right = '0';
  sidebar.style.top = '0';
  sidebar.style.width = '400px';  // Increased width
  sidebar.style.height = '100%';
  sidebar.style.backgroundColor = '#f1f1f1';
  sidebar.style.borderLeft = '1px solid #ccc';
  sidebar.style.zIndex = '10000';
  sidebar.style.overflowY = 'auto';
  sidebar.innerHTML = `
      <div id="close-sidebar" style="position: absolute; top: 10px; right: 10px; cursor: pointer;">X</div>
      <div style="padding: 10px;">
          <h2>Welcome to ReadRobin</h2>
          <h3>Article Title</h3>
          <button id="generate-quiz">Generate Quiz</button>
          <div id="quiz-content"></div>
      </div>
  `;

  document.body.appendChild(sidebar);

  // Debug statement to ensure the sidebar is added
  console.log('Sidebar added to the DOM.');

  document.getElementById('close-sidebar').addEventListener('click', () => {
    sidebar.remove();
  });

  const generateQuizButton = document.getElementById('generate-quiz');

  // Debug statement to ensure the button is found
  console.log('Generate Quiz button found:', generateQuizButton);

  generateQuizButton.addEventListener('click', () => {
    // Debug statement to ensure the click event is triggered
    console.log('Generate Quiz button clicked.');
    generateQuiz();
  });

  function generateQuiz() {
    chrome.storage.local.get(['userId'], (result) => {
      if (result.userId) {
        const userId = result.userId;
        const url = window.location.href;
        const html = document.documentElement.innerHTML;

        console.log('Sending request to backend with:', { url, html, userId });

        chrome.runtime.sendMessage({
          action: 'generateQuiz',
          data: { url, html, userId }
        }, (response) => {
          if (response.error) {
            console.error('Error from backend:', response.error);
            alert('Failed to fetch data. Please try again later.');
            return;
          }
          console.log('Received response from backend:', response);
          const quizContent = document.getElementById('quiz-content');
          if (!response.questions_answers || !Array.isArray(response.questions_answers)) {
            console.error('Invalid questions and answers format:', response.questions_answers);
            alert('Failed to load quiz questions.');
            return;
          }
          quizContent.innerHTML = `
              <h3>Quiz Questions</h3>
              ${response.questions_answers.map((qa, index) => `
                  <div class="question-container">
                      <p>${index + 1}. ${qa.question}</p>
                      <input type="text" id="answer-${index}" placeholder="Type your answer here">
                      <button onclick="submitAnswer(${index})">Submit</button>
                  </div>
              `).join('')}
          `;
        });
      } else {
        console.log('No user ID found');
      }
    });
  }

  window.submitAnswer = function(index) {
    const answerInput = document.getElementById(`answer-${index}`);
    const userAnswer = answerInput.value;
    console.log(`User answer for question ${index + 1}: ${userAnswer}`);
    // Handle answer submission logic here
  };
})();
