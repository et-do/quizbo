import React, { useState } from "react";
import "./App.css";

function QuizForm({ user, setPage }) {
  const [url, setUrl] = useState("");
  const [contentID, setContentID] = useState(null);
  const [quizID, setQuizID] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setContentID(null);
    setQuizID(null);
    setLoading(true);

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
      setPage("quizPage", { contentID: data.content_id, quizID: data.quiz_id });
      setLoading(false);
    } catch (error) {
      setError("Error submitting URL");
      setLoading(false);
    }
  };

  return (
    <div className="quiz-form">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Webpage Quiz</h2>
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
      </form>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && (
        <div className="loading-spinner">Generating your Quiz...</div>
      )}
    </div>
  );
}

export default QuizForm;
