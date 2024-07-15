// App.js
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
import Login from "./Login";
import Sidebar from "./Sidebar";
import AttemptPage from "./AttemptPage";

function App() {
  const [page, setPage] = useState("login");
  const [user, setUser] = useState(null);
  const [contentID, setContentID] = useState(null);
  const [attemptID, setAttemptID] = useState(null);
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
        return <Login signIn={signIn} />;
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
        return (
          <QuizPage
            user={user}
            setPage={setPage}
            contentID={contentID}
            quizID={quizID}
          />
        );
      case "attemptPage":
        return (
          <AttemptPage
            user={user}
            contentID={contentID}
            attemptID={attemptID}
            setPage={setPage}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className="App">
      <header className="app-header">
        <div className="header-top-row">
          <div className="logo-title">
            <img src={logo} alt="Logo" className="logo" />
            <h1 className="app-title">ReadRobin</h1>
          </div>
          <h2 className="tagline">
            Your AI Companion for Smarter Comprehension
          </h2>
        </div>
        {user && (
          <div className="user-info">
            <p>Welcome, {user.displayName}</p>
            <button className="logout" onClick={logout}>
              Logout
            </button>
          </div>
        )}
      </header>
      <div className="main-content">
        <Sidebar
          user={user}
          setContentID={setContentID}
          setAttemptID={setAttemptID}
          setPage={setPage}
        />
        <div className="page-content">{renderPage()}</div>
      </div>
    </div>
  );
}

export default App;
