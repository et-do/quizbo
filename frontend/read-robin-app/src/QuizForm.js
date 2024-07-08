import React, { useState } from "react";

function QuizForm({ user, setPage }) {
  const [url, setUrl] = useState("");
  const [contentID, setContentID] = useState(null);
  const [quizID, setQuizID] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [responses, setResponses] = useState({});
  const [status, setStatus] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setContentID(null);
    setQuizID(null);
    setQuestions([]);
    setLoading(true);

    if (!user) {
      setError("You must be logged in to submit a URL");
      setLoading(false);
      return;
    }

    try {
      const idToken = await user.getIdToken();
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/submit`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${idToken}`,
          },
          body: JSON.stringify({ url }),
        }
      );
      const data = await res.json();
      setContentID(data.content_id);
      setQuizID(data.quiz_id);
      fetchQuestions(data.content_id, data.quiz_id);
    } catch (error) {
      setError("Error submitting URL");
      setLoading(false);
    }
  };

  const fetchQuestions = async (contentID, quizID) => {
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
    <div>
      <form onSubmit={handleSubmit}>
        <label>
          URL:
          <input
            type="text"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
        </label>
        <button type="submit">Submit</button>
        <button type="button" onClick={() => setPage("selection")}>
          Back
        </button>
      </form>
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
              <input
                type="text"
                value={responses[index] || ""}
                onChange={(e) => handleResponseChange(e, index)}
              />
              <button onClick={() => handleSubmitResponse(index, item.answer)}>
                Submit
              </button>
              {status[index] && (
                <span
                  className={
                    status[index] === "Correct" ? "correct" : "incorrect"
                  }
                >
                  {status[index]}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default QuizForm;
