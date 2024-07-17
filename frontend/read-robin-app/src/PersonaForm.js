import React, { useState } from "react";
import { doc, updateDoc, arrayUnion } from "firebase/firestore";
import { db } from "./firebase";
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
        <p className="persona-placeholder">
          You are a{" "}
          <input
            type="text"
            value={userType}
            onChange={(e) => setUserType(e.target.value)}
            placeholder="(profession, age, etc)"
            required
          />{" "}
          looking for quizzes of{" "}
          <input
            type="text"
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value)}
            placeholder="(beginner, intermediate, expert)"
            required
          />{" "}
          difficulty.
        </p>
        <div>
          <label>
            Persona Name:
            <input
              type="text"
              value={personaName}
              onChange={(e) => setPersonaName(e.target.value)}
              required
            />
          </label>
        </div>
        <button type="submit" className="persona-submit-button">
          Add Persona
        </button>
      </form>
    </div>
  );
};

export default PersonaForm;
