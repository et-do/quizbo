import React, { useState } from "react";
import { doc, setDoc } from "firebase/firestore";
import { db } from "./firebase";
import { v4 as uuidv4 } from "uuid"; // Add UUID for unique IDs
import "./PersonaForm.css";

const PersonaForm = ({ user, addPersona }) => {
  // Add addPersona prop
  const [personaName, setPersonaName] = useState("");
  const [userType, setUserType] = useState("");
  const [difficulty, setDifficulty] = useState("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!user) return;

    const personaId = uuidv4(); // Create a unique ID
    const newPersona = {
      id: personaId,
      name: personaName,
      type: userType,
      difficulty: difficulty,
    };

    const personaRef = doc(db, "users", user.uid, "personas", personaId);
    await setDoc(personaRef, newPersona);

    addPersona(newPersona); // Update the state in the parent component

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
        <div>
          <label>
            Describe yourself (e.g., student, CEO, researcher):
            <input
              type="text"
              value={userType}
              onChange={(e) => setUserType(e.target.value)}
              placeholder="e.g., student, CEO, researcher"
              required
            />
          </label>
        </div>
        <div>
          <label>
            What level of difficulty are you looking for? (e.g., easy, medium,
            expert):
            <input
              type="text"
              value={difficulty}
              onChange={(e) => setDifficulty(e.target.value)}
              placeholder="e.g., easy, medium, expert"
              required
            />
          </label>
        </div>
        <div>
          <label>
            Give a name to this persona:
            <input
              type="text"
              value={personaName}
              onChange={(e) => setPersonaName(e.target.value)}
              placeholder="Persona Name..."
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
