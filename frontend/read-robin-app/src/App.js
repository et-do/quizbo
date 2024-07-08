import React, { useState, useEffect } from "react";
import "./App.css";
import { auth } from "./firebase";
import {
  signInWithPopup,
  GoogleAuthProvider,
  onAuthStateChanged,
  signOut,
} from "firebase/auth";
import logo from "./logo.png";
import SelectionPage from "./SelectionPage";
import QuizForm from "./QuizForm";
import QuizPage from "./QuizPage";

function App() {
  const [page, setPage] = useState("login");
  const [user, setUser] = useState(null);
  const [contentID, setContentID] = useState(null);
  const [quizID, setQuizID] = useState(null);
  const provider = new GoogleAuthProvider();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      if (user) {
        setUser(user);
        setPage("selection");
      } else {
        setUser(null);
        setPage("login");
      }
    });

    return () => unsubscribe();
  }, []);

  const signIn = () => {
    signInWithPopup(auth, provider)
      .then((result) => {
        setUser(result.user);
        setPage("selection");
      })
      .catch((error) => {
        console.error("Error signing in: ", error);
      });
  };

  const logout = () => {
    signOut(auth)
      .then(() => {
        setUser(null);
        setPage("login");
      })
      .catch((error) => {
        console.error("Error signing out: ", error);
      });
  };

  const renderPage = () => {
    switch (page) {
      case "login":
        return (
          <div className="login-page">
            <button onClick={signIn}>Sign in with Google</button>
          </div>
        );
      case "selection":
        return <SelectionPage setPage={setPage} />;
      case "quizForm":
        return (
          <QuizForm
            user={user}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "quizPage":
        return <QuizPage contentID={contentID} quizID={quizID} />;
      default:
        return null;
    }
  };

  return (
    <div className="App">
      <header>
        <img src={logo} alt="Logo" />
        <h1>Your AI Companion for Smarter Comprehension</h1>
        {user && (
          <div className="user-info">
            <p>Welcome, {user.displayName}</p>
            <button className="logout" onClick={logout}>
              Logout
            </button>
          </div>
        )}
      </header>
      {renderPage()}
    </div>
  );
}

export default App;
