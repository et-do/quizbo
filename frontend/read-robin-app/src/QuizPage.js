import React, { useState, useEffect } from "react";
import "./App.css";

function QuizPage({ contentID, quizID }) {
  const [questions, setQuestions] = useState([]);
  const [responses, setResponses] = useState({});
  const [status, setStatus] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchQuestions = async () => {
      try {
        const res = await fetch(
          `https://read-robin-dev-6yudia4zva-nn.a.run.app/quiz/${contentID}/${quizID}`
        );
        const data = await res.json();
        setQuestions(data.questions);
        setLoading(false);
      } catch (error) {
        setError("Error fetching questions");
        setLoading(false);
      }
    };
    fetchQuestions();
  }, [contentID, quizID]);

  const handleResponseChange = (e, index) => {
    const newResponses = { ...responses, [index]: e.target.value };
    setResponses(newResponses);
  };

  const handleSubmitResponse = (index, correctAnswer) => {
    const isCorrect = responses[index] === correctAnswer;
    const newStatus = {
      ...status,
      [index]: isCorrect ? "Correct" : "Incorrect",
    };
    setStatus(newStatus);
  };

  return (
    <div className="quiz-page">
      <h2>Quiz</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && (
        <div className="loading-spinner">Generating your Quiz...</div>
      )}
      {questions.length > 0 && (
        <div>
          <h2>Questions</h2>
          {questions.map((item, index) => (
            <div key={index} className="quiz-item">
              <h3>Question {index + 1}</h3>
              <p>{item.question}</p>
              <div className="response-container">
                <input
                  type="text"
                  value={responses[index] || ""}
                  onChange={(e) => handleResponseChange(e, index)}
                />
                <button
                  onClick={() => handleSubmitResponse(index, item.answer)}
                >
                  Submit
                </button>
              </div>
              {status[index] && (
                <div
                  className={
                    status[index] === "Correct" ? "correct" : "incorrect"
                  }
                >
                  {status[index]}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default QuizPage;
