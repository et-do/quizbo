import React, { useState } from "react";
import { doc, updateDoc, arrayUnion } from "firebase/firestore";
import { db } from "./firebase";
import { v4 as uuidv4 } from "uuid"; // Add UUID for unique IDs
import "./PersonaForm.css";

const PersonaForm = ({ user }) => {
  const [personaName, setPersonaName] = useState("");
  const [userType, setUserType] = useState("");
  const [difficulty, setDifficulty] = useState("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!user) return;

    const userRef = doc(db, "users", user.uid);

    const newPersona = {
      id: uuidv4(), // Add a unique ID
      name: personaName,
      type: userType,
      difficulty: difficulty,
    };

    await updateDoc(userRef, {
      personas: arrayUnion(newPersona),
    });

    setPersonaName("");
    setUserType("");
    setDifficulty("");
  };

  return (
    <div className="persona-form-container">
      <form onSubmit={handleSubmit} className="persona-form">
        <h2>Create a Persona</h2>
        <p className="persona-subtext">
          Make your personalized quiz persona to tailor questions and difficulty
          to your goals!
        </p>
        <p className="persona-placeholder">
          You are a{" "}
          <input
            type="text"
            value={userType}
            onChange={(e) => setUserType(e.target.value)}
            placeholder="student, CEO, researcher"
            required
          />{" "}
          looking for quizzes of{" "}
          <input
            type="text"
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value)}
            placeholder="ex: easy, medium, expert"
            required
          />{" "}
          difficulty.
        </p>
        <div>
          <input
            type="text"
            value={personaName}
            onChange={(e) => setPersonaName(e.target.value)}
            placeholder="Persona Name..."
            required
          />
        </div>
        <button type="submit" className="persona-submit-button">
          Add Persona
        </button>
      </form>
    </div>
  );
};

export default PersonaForm;
