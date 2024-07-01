import React, { useState, useEffect } from "react";
import "./App.css";
import { auth } from "./firebase";
import {
  signInWithPopup,
  GoogleAuthProvider,
  onAuthStateChanged,
  signOut,
} from "firebase/auth";
import logo from "./logo.png"; // Ensure you have a logo.png in the src directory

function App() {
  const [url, setUrl] = useState("");
  const [quizID, setQuizID] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [error, setError] = useState(null);
  const [user, setUser] = useState(null);

  const provider = new GoogleAuthProvider();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      if (user) {
        setUser(user);
      } else {
        setUser(null);
      }
    });

    return () => unsubscribe();
  }, []);

  const signIn = () => {
    signInWithPopup(auth, provider)
      .then((result) => {
        setUser(result.user);
      })
      .catch((error) => {
        setError(error.message);
      });
  };

  const logout = () => {
    signOut(auth)
      .then(() => {
        setUser(null);
      })
      .catch((error) => {
        setError(error.message);
      });
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setQuizID(null);
    setQuestions([]);

    if (!user) {
      setError("You must be logged in to submit a URL");
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
      setQuizID(data.quiz_id);
      fetchQuestions(data.quiz_id);
    } catch (error) {
      setError("Error submitting URL");
    }
  };

  const fetchQuestions = async (quizID) => {
    try {
      const res = await fetch(
        `https://read-robin-6yudia4zva-nn.a.run.app/quiz/${quizID}`
      );
      const data = await res.json();
      setQuestions(data.questions);
    } catch (error) {
      setError("Error fetching questions");
    }
  };

  return (
    <div className="App">
      <header>
        <img src={logo} alt="Logo" />
        <h1>Submit a URL</h1>
      </header>
      {user ? (
        <div>
          <p>Welcome, {user.displayName}</p>
          <button className="logout" onClick={logout}>
            Logout
          </button>
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
        </div>
      ) : (
        <button onClick={signIn}>Sign in with Google</button>
      )}
      {error && <div style={{ color: "red" }}>{error}</div>}
      {questions.length > 0 && (
        <div>
          <h2>Questions</h2>
          {questions.map((item, index) => (
            <div key={index} className="quiz-item">
              <h3>Question {index + 1}</h3>
              <p>{item}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
