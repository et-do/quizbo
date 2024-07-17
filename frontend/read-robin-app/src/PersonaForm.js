import React, { useState } from "react";
import { doc, setDoc } from "firebase/firestore";
import { db } from "./firebase";
import { v4 as uuidv4 } from "uuid";
import "./PersonaForm.css";

const PersonaForm = ({ user, addPersona }) => {
  const [personaName, setPersonaName] = useState("");
  const [userRole, setUserRole] = useState("");
  const [language, setLanguage] = useState("");
  const [difficulty, setDifficulty] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false); // Add submitting state

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!user) return;

    setIsSubmitting(true); // Start animation

    const personaId = uuidv4();
    const newPersona = {
      id: personaId,
      name: personaName,
      role: userRole,
      language: language,
      difficulty: difficulty,
    };

    const personaRef = doc(db, "users", user.uid, "personas", personaId);
    await setDoc(personaRef, newPersona);

    addPersona(newPersona);

    setPersonaName("");
    setUserRole("");
    setDifficulty("");
    setLanguage("");

    setTimeout(() => setIsSubmitting(false), 1000); // End animation after 1 second
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
            Describe yourself:
            <input
              type="text"
              value={userRole}
              onChange={(e) => setUserRole(e.target.value)}
              placeholder="e.g., student, CEO, researcher"
              required
            />
          </label>
        </div>
        <div>
          <label>
            What level of difficulty are you looking for?:
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
            What language do you want the questions to be in?:
            <input
              type="text"
              value={language}
              onChange={(e) => setLanguage(e.target.value)}
              placeholder="e.g., english, japanese, spanish"
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
        <button
          type="submit"
          className={`persona-submit-button ${
            isSubmitting ? "submitting" : ""
          }`}
        >
          {isSubmitting ? "Adding..." : "Add Persona"}
        </button>
      </form>
    </div>
  );
};

export default PersonaForm;
