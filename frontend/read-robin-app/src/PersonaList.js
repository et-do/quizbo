import React from "react";
import { doc, updateDoc, arrayRemove } from "firebase/firestore";
import { db } from "./firebase";
import "./PersonaList.css";

const PersonaList = ({ user, personas, activePersona, setActivePersona }) => {
  const handleDelete = async (persona) => {
    if (!user) return;

    const userRef = doc(db, "users", user.uid);

    await updateDoc(userRef, {
      personas: arrayRemove(persona),
    });

    if (activePersona && activePersona.id === persona.id) {
      setActivePersona(null);
    }
  };

  return (
    <div className="persona-list-container">
      <h2>Your Personas</h2>
      <ul className="persona-list">
        {personas.map((persona) => (
          <li
            key={persona.id}
            className={`persona-item ${
              activePersona && activePersona.id === persona.id ? "active" : ""
            }`}
          >
            <div className="persona-info">
              <strong>{persona.name}</strong> - {persona.type} -{" "}
              {persona.difficulty}
            </div>
            <div className="persona-actions">
              <button
                className="persona-button"
                onClick={() => setActivePersona(persona)}
              >
                Set Active
              </button>
              <button
                className="persona-button delete-button"
                onClick={() => handleDelete(persona)}
              >
                Delete
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PersonaList;
