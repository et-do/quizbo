import React, { useState, useEffect } from "react";
import "./App.css";
import { auth } from "./firebase";
import {
  signInWithPopup,
  GoogleAuthProvider,
  onAuthStateChanged,
  signOut,
} from "firebase/auth";
import { createUserProfile } from "./UserProfile"; // Import the function
import logo from "./logo.png";
import SelectionPage from "./SelectionPage";
import QuizForm from "./QuizForm";
import QuizPage from "./QuizPage";
import Login from "./Login";
import Sidebar from "./Sidebar";
import AttemptPage from "./AttemptPage";
import IntroScreen from "./IntroScreen";

function App() {
  const [page, setPage] = useState("login");
  const [user, setUser] = useState(null);
  const [contentID, setContentID] = useState(null);
  const [attemptID, setAttemptID] = useState(null);
  const [quizID, setQuizID] = useState(null);
  const [showIntro, setShowIntro] = useState(false);
  const provider = new GoogleAuthProvider();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (user) {
        setUser(user);
        await createUserProfile(user); // Create/update the user profile
        const hasSeenIntro = localStorage.getItem("hasSeenIntro");
        console.log("User logged in:", user); // Debugging log
        console.log("Has seen intro:", hasSeenIntro); // Debugging log
        if (!hasSeenIntro) {
          setShowIntro(true);
        } else {
          setPage("selection");
        }
      } else {
        setUser(null);
        setPage("login");
      }
    });

    return () => unsubscribe();
  }, []);

  const signIn = () => {
    signInWithPopup(auth, provider)
      .then(async (result) => {
        setUser(result.user);
        await createUserProfile(result.user); // Create/update the user profile
        const hasSeenIntro = localStorage.getItem("hasSeenIntro");
        console.log("User signed in:", result.user); // Debugging log
        console.log("Has seen intro:", hasSeenIntro); // Debugging log
        if (!hasSeenIntro) {
          setShowIntro(true);
        } else {
          setPage("selection");
        }
      })
      .catch((error) => {
        console.error("Error signing in:", error);
      });
  };

  const logout = () => {
    signOut(auth)
      .then(() => {
        setUser(null);
        setPage("login");
      })
      .catch((error) => {
        console.error("Error signing out:", error);
      });
  };

  const finishIntro = () => {
    localStorage.setItem("hasSeenIntro", "true");
    setShowIntro(false);
    setPage("selection");
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
      {showIntro && <IntroScreen onFinish={finishIntro} />}
      {!showIntro && (
        <>
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
                <button className="logout-button" onClick={logout}>
                  Logout
                </button>
              </div>
            )}
          </header>
          <div className="main-content">
            {user && (
              <Sidebar
                user={user}
                setContentID={setContentID}
                setAttemptID={setAttemptID}
                setPage={setPage}
              />
            )}
            <div className="page-content">{renderPage()}</div>
          </div>
        </>
      )}
    </div>
  );
}

export default App;
