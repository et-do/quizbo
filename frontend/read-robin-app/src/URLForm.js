import React, { useState } from "react";
import "./URLForm.css";
import { db } from "./firebase";
import { doc, setDoc } from "firebase/firestore";

function UrlForm({ user, activePersona, setPage, setContentID, setQuizID }) {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const exampleUrls = [
    {
      text: "The World's Largest Lobster",
      url: "https://en.wikipedia.org/wiki/The_World%27s_Largest_Lobster",
    },
    {
      text: "The Death of Saas",
      url: "https://medium.com/@akiranin/the-death-of-saas-a1c5423da094",
    },
    {
      text: "How to Cook Beef Wellington",
      url: "https://www.gordonramsay.com/gr/recipes/beef-wellington/",
    },
  ];

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (!user || !activePersona || !activePersona.id) {
        throw new Error("User or active persona is not defined");
      }

      const payload = {
        url: url,
        persona: {
          id: activePersona.id,
          name: activePersona.name,
          role: activePersona.role,
          language: activePersona.language,
          difficulty: activePersona.difficulty,
        },
        content_type: "URL",
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
        throw new Error(`Error submitting URL: ${res.statusText}`);
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
        content_type: "URL",
      });

      setPage("quizPage");
      setLoading(false);
    } catch (error) {
      console.error("Error:", error);
      setError(`Error submitting URL: ${error.message}`);
      setLoading(false);
    }
  };

  const handleExampleClick = (exampleUrl) => {
    setUrl(exampleUrl);
  };

  return (
    <div className="quiz-form">
      <button
        className="quiz-form-back-button"
        onClick={() => setPage("selection")}
      >
        Back
      </button>
      <h2>Webpage Quiz</h2>
      <div className="example-urls">
        <div className="example-card">
          <h3>Try these examples:</h3>
          <ul>
            {exampleUrls.map((example, index) => (
              <li key={index}>
                <a
                  href={example.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="example-link"
                >
                  {example.text}
                </a>
                <button
                  onClick={() => handleExampleClick(example.url)}
                  className="use-url-button"
                >
                  Use this URL
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Enter your URL here..."
          required
        />
        <button type="submit">Submit</button>
      </form>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="loading-spinner"></div>}
    </div>
  );
}

export default UrlForm;
