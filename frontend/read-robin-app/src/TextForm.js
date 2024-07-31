import React, { useState } from "react";
import "./TextForm.css";
import { db } from "./firebase";
import { doc, setDoc } from "firebase/firestore";

function TextForm({ user, activePersona, setPage, setContentID, setQuizID }) {
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (!user || !activePersona || !activePersona.id) {
        throw new Error("User or active persona is not defined");
      }

      const payload = {
        url: title,
        content_text: text,
        persona: {
          id: activePersona.id,
          name: activePersona.name,
          role: activePersona.role,
          language: activePersona.language,
          difficulty: activePersona.difficulty,
        },
        content_type: "Text",
      };
      console.log("Payload being sent to backend:", payload);

      const idToken = await user.getIdToken();
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/submit`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${idToken}`,
          },
          body: JSON.stringify(payload),
        }
      );

      if (!res.ok) {
        throw new Error(`Error submitting text: ${res.statusText}`);
      }

      const data = await res.json();
      setContentID(data.content_id);
      setQuizID(data.quiz_id);

      const quizRef = doc(
        db,
        "users",
        user.uid,
        "personas",
        activePersona.id,
        "quizzes",
        data.content_id
      );
      await setDoc(quizRef, {
        contentID: data.content_id,
        url: data.url,
        title: data.title,
        content_text: data.content_text,
        content_type: "Text",
      });

      setPage("quizPage");
      setLoading(false);
    } catch (error) {
      console.error("Error:", error);
      setError(`Error submitting text: ${error.message}`);
      setLoading(false);
    }
  };

  return (
    <div className="quiz-form">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Text Quiz</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter your title here..."
          required
        />
        <textarea
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Enter your text here..."
          rows="10"
          required
        />
        <button type="submit">Submit</button>
      </form>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="loading-spinner"></div>}
    </div>
  );
}

export default TextForm;
