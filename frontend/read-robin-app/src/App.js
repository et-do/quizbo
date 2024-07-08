import React, { useState, useEffect } from "react";
import "./App.css";
import { auth } from "./firebase";
import { onAuthStateChanged, signOut } from "firebase/auth";
import logo from "./logo.png";
import Login from "./Login";
import SelectionPage from "./SelectionPage";
import QuizForm from "./QuizForm";

function App() {
  const [user, setUser] = useState(null);
  const [page, setPage] = useState("selection");

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

  const logout = () => {
    signOut(auth)
      .then(() => {
        setUser(null);
      })
      .catch((error) => {
        console.error("Error signing out: ", error);
      });
  };

  return (
    <div className="App">
      <header>
        <img src={logo} alt="Logo" />
        <h1>Your AI Companion for Smarter Comprehension</h1>
      </header>
      {user ? (
        <div>
          <p>Welcome, {user.displayName}</p>
          <button className="logout" onClick={logout}>
            Logout
          </button>
          {page === "selection" && <SelectionPage setPage={setPage} />}
          {page === "quizForm" && <QuizForm user={user} setPage={setPage} />}
        </div>
      ) : (
        <Login />
      )}
    </div>
  );
}

export default App;
