import React, { useState } from "react";
import { doc, deleteDoc } from "firebase/firestore";
import { db, auth } from "./firebase"; // Import auth from firebase
import "./App.css";

function HowTo({ user, setPage }) {
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleDeleteProfile = async () => {
    setLoading(true);
    setError(null);
    try {
      console.log("Starting profile deletion process...");

      // Delete user profile from Firestore
      const userDocRef = doc(db, "users", user.uid);
      console.log("Deleting user document from Firestore with ID:", user.uid);
      await deleteDoc(userDocRef);
      console.log("User document deleted from Firestore.");

      // Sign out the user
      console.log("Signing out the user...");
      await auth.signOut();
      console.log("User signed out.");

      // Optionally, delete the user's Firebase Authentication profile
      console.log("Deleting the user's Firebase Authentication profile...");
      await user.delete();
      console.log("User's Firebase Authentication profile deleted.");

      // Redirect to login page or another appropriate page
      console.log("Redirecting to login page...");
      setPage("login");
    } catch (error) {
      console.error("Error during profile deletion process:", error);
      setError("Error deleting profile: " + error.message);
    } finally {
      setLoading(false);
      setShowConfirmation(false);
    }
  };

  return (
    <div className="howto-page">
      <h1>How to Use Quizbo</h1>
      <p>
        Quizbo is your AI companion for smarter comprehension. Follow the steps
        below to get the most out of Quizbo!
      </p>
      <h2>1. Manage Personas</h2>
      <p>
        Go to{" "}
        <strong
          onClick={() => setPage("personas")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Personas
        </strong>{" "}
        to create a persona. Creating a persona will allow you to take tests
        however you'd like! Define who you are, your difficulty choice, language
        choice, and give that persona a name. Now all your quizzes under that
        active persona will be tailored for you! To change or add personas, go
        to the{" "}
        <strong
          onClick={() => setPage("personas")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Personas
        </strong>{" "}
        tab.
      </p>
      <h2>2. Generate a Quiz</h2>
      <p>
        Once you have a persona, go and{" "}
        <strong
          onClick={() => setPage("selection")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          generate a quiz
        </strong>
        ! You can generate quizzes from website HTML through URLs, PDFs, audio,
        or video, as well as plain text. Once you've generated a quiz, you can
        answer as few questions as you'd like, and get instant feedback from
        Quizbo on how you were right or wrong, as well as learning tips!
      </p>
      <h2>3. Performance History</h2>
      <p>
        Once you've completed an attempt, your performance history will be saved
        on the sidebar and in the{" "}
        <strong
          onClick={() => setPage("performanceHistory")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Performance History
        </strong>{" "}
        tab. Clicking on any of the attempts in the sidebar will bring you to
        the attempts page where you can see your answers, the correct answers,
        and references from the text itself. Check out the{" "}
        <strong
          onClick={() => setPage("performanceHistory")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Performance History
        </strong>{" "}
        tab to get historic trends on your learning performance.
      </p>
      <h2>4. Manage Content</h2>
      <p>
        Want to quiz yourself again on the same content? Navigate to the{" "}
        <strong
          onClick={() => setPage("contentManagement")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Content
        </strong>{" "}
        tab to regenerate quizzes based on content you've already submitted. You
        can also see your latest scores to ensure you are up to date with all
        your important content.
      </p>

      <h2>Delete Profile</h2>
      <p>
        If you wish to delete your profile, you can do so by clicking the button
        below. Please note that this action is irreversible. While you'll be
        able to make another account, you will lose all of your history.
      </p>
      <button
        onClick={() => setShowConfirmation(true)}
        className="delete-button"
      >
        Delete Profile
      </button>

      {showConfirmation && (
        <div className="confirmation-modal">
          <div className="confirmation-content">
            <h3>Are you sure you want to delete your profile?</h3>
            <p>This action cannot be undone.</p>
            <button
              onClick={handleDeleteProfile}
              disabled={loading}
              className="confirm-button"
            >
              {loading ? "Deleting..." : "Yes, Delete"}
            </button>
            <button
              onClick={() => setShowConfirmation(false)}
              className="cancel-button"
            >
              Cancel
            </button>
            {error && <p style={{ color: "red" }}>{error}</p>}
          </div>
        </div>
      )}
    </div>
  );
}

export default HowTo;
