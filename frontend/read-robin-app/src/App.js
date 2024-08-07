import React, { useState, useEffect } from "react";
import "./App.css";
import { auth } from "./firebase";
import {
  signInWithPopup,
  GoogleAuthProvider,
  onAuthStateChanged,
  signOut,
} from "firebase/auth";
import { createUserProfile } from "./UserProfile";
import {
  doc,
  getDoc,
  updateDoc,
  collection,
  getDocs,
} from "firebase/firestore";
import { db } from "./firebase";
import logo from "./logo.png";
import SelectionPage from "./SelectionPage";
import QuizPage from "./QuizPage";
import Login from "./Login";
import Sidebar from "./Sidebar";
import AttemptPage from "./AttemptPage";
import IntroScreen from "./IntroScreen";
import PersonaForm from "./PersonaForm";
import PersonaList from "./PersonaList";
import UrlForm from "./URLForm";
import PdfForm from "./PdfForm";
import AudioForm from "./AudioForm";
import VideoForm from "./VideoForm";
import TextForm from "./TextForm";
import PerformanceHistory from "./PerformanceHistory";
import ContentManagementPage from "./ContentManagementPage";

function App() {
  const [page, setPage] = useState("intro");
  const [user, setUser] = useState(null);
  const [personas, setPersonas] = useState([]);
  const [activePersona, setActivePersona] = useState(null);
  const [contentID, setContentID] = useState(null);
  const [attemptID, setAttemptID] = useState(null);
  const [quizID, setQuizID] = useState(null);
  const [showIntro, setShowIntro] = useState(null);
  const provider = new GoogleAuthProvider();

  useEffect(() => {
    const hasSeenIntro = localStorage.getItem("hasSeenIntro");
    if (hasSeenIntro) {
      setShowIntro(false);
      setPage("login");
    } else {
      setShowIntro(true);
    }

    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (user) {
        setUser(user);
        await createUserProfile(user);

        const userRef = doc(db, "users", user.uid);
        const userSnap = await getDoc(userRef);
        if (userSnap.exists()) {
          const userData = userSnap.data();
          setPersonas(userData.personas || []);
          setActivePersona(userData.activePersona || null); // Load active persona

          const personaCollection = collection(
            db,
            "users",
            user.uid,
            "personas"
          );
          const personaSnapshot = await getDocs(personaCollection);
          const personaList = personaSnapshot.docs.map((doc) => ({
            id: doc.id,
            ...doc.data(),
          }));
          setPersonas(personaList);
        }

        const hasSeenIntro = localStorage.getItem("hasSeenIntro");
        console.log("User logged in:", user);
        console.log("Has seen intro:", hasSeenIntro);
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
        await createUserProfile(result.user);
        const userRef = doc(db, "users", result.user.uid);
        const userSnap = await getDoc(userRef);
        if (userSnap.exists()) {
          const userData = userSnap.data();
          setPersonas(userData.personas || []);
          setActivePersona(userData.activePersona || null);

          const personaCollection = collection(
            db,
            "users",
            result.user.uid,
            "personas"
          );
          const personaSnapshot = await getDocs(personaCollection);
          const personaList = personaSnapshot.docs.map((doc) => ({
            id: doc.id,
            ...doc.data(),
          }));
          setPersonas(personaList);
        }
        const hasSeenIntro = localStorage.getItem("hasSeenIntro");
        console.log("User signed in:", result.user);
        console.log("Has seen intro:", hasSeenIntro);
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
    setPage(user ? "selection" : "login");
  };

  const handleSetActivePersona = async (persona) => {
    setActivePersona(persona);
    if (user) {
      const userRef = doc(db, "users", user.uid);
      await updateDoc(userRef, {
        activePersona: persona,
      });
    }
  };

  const addPersona = (newPersona) => {
    setPersonas((prevPersonas) => [...prevPersonas, newPersona]);
  };

  const renderPage = () => {
    switch (page) {
      case "login":
        return <Login signIn={signIn} />;
      case "selection":
        return <SelectionPage setPage={setPage} />;
      case "urlForm":
        return (
          <UrlForm
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "pdfForm":
        return (
          <PdfForm
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "audioForm":
        return (
          <AudioForm
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "videoForm":
        return (
          <VideoForm
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "textForm":
        return (
          <TextForm
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setContentID={setContentID}
            setQuizID={setQuizID}
          />
        );
      case "quizPage":
        return (
          <QuizPage
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            contentID={contentID}
            quizID={quizID}
          />
        );
      case "performanceHistory":
        return (
          <PerformanceHistory
            user={user}
            activePersona={activePersona}
            setPage={setPage}
            setAttemptID={setAttemptID}
            setContentID={setContentID}
          />
        );
      case "attemptPage":
        return (
          <AttemptPage
            user={user}
            activePersona={activePersona}
            contentID={contentID}
            attemptID={attemptID}
            setPage={setPage}
          />
        );
      case "intro":
        return <IntroScreen onFinish={finishIntro} />;
      case "personas":
        return (
          <>
            <PersonaForm user={user} addPersona={addPersona} />
            <PersonaList
              user={user}
              personas={personas}
              activePersona={activePersona}
              setActivePersona={handleSetActivePersona}
              setPage={setPage}
            />
            <button
              className="back-button"
              onClick={() => setPage("selection")}
            >
              Back
            </button>
          </>
        );
      case "contentManagement":
        return (
          <>
            <ContentManagementPage
              user={user}
              activePersona={activePersona}
              setPage={setPage}
              setContentID={setContentID}
              setQuizID={setQuizID}
            />
            <button
              className="back-button"
              onClick={() => setPage("selection")}
            >
              Back
            </button>
          </>
        );
      default:
        return null;
    }
  };

  if (showIntro === null) {
    // Return null or a loading indicator while determining the state
    return null;
  }

  return (
    <div className="App">
      {showIntro ? (
        <IntroScreen onFinish={finishIntro} />
      ) : (
        <>
          <header className="app-header">
            <div className="header-top-row">
              <div
                className="logo-title"
                onClick={() => setPage("selection")}
                style={{ cursor: "pointer" }}
              >
                <img src={logo} alt="Logo" className="logo" />
                <h1 className="app-title">Quizbo</h1>
              </div>
              <h2 className="tagline">
                Your AI Companion for Smarter Comprehension
              </h2>
              <div
                className="refresh-icon"
                onClick={() => window.location.reload()}
              ></div>
            </div>
            {user && (
              <div className="user-info">
                <p>Welcome, {user.displayName}</p>
                {activePersona && (
                  <div className="active-persona-card">
                    <h3
                      className="active-persona-title"
                      onClick={() => setPage("personas")}
                    >
                      Active Persona
                    </h3>
                    <div className="active-persona-details">
                      <p>
                        <strong>Name: </strong> {activePersona.name}
                      </p>
                      <p>
                        <strong>Role: </strong> {activePersona.role}
                      </p>
                      <p>
                        <strong>Language: </strong> {activePersona.language}
                      </p>
                      <p>
                        <strong>Difficulty: </strong> {activePersona.difficulty}
                      </p>
                    </div>
                  </div>
                )}
                <div className="button-container">
                  <button
                    className="button-common generate-quiz-button"
                    onClick={() => setPage("selection")}
                  >
                    Generate Quiz
                  </button>
                  <button
                    className="button-common manage-personas-button"
                    onClick={() => setPage("personas")}
                  >
                    Manage Personas
                  </button>
                  <button
                    className="button-common content-management-button"
                    onClick={() => setPage("contentManagement")}
                  >
                    Manage Content
                  </button>
                  <button
                    className="button-common performance-history-button"
                    onClick={() => setPage("performanceHistory")}
                  >
                    Performance History
                  </button>
                  <button
                    className="button-common logout-button"
                    onClick={logout}
                  >
                    Logout
                  </button>
                </div>
              </div>
            )}
          </header>
          <div className="main-content">
            {user && (
              <Sidebar
                user={user}
                activePersona={activePersona}
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
