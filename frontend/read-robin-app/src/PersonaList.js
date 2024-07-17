import React from "react";
import { doc, updateDoc, arrayRemove } from "firebase/firestore";
import { db } from "./firebase";

const PersonaList = ({ user, personas, activePersona, setActivePersona }) => {
  const handleDelete = async (persona) => {
    if (!user) return;

    const userRef = doc(db, "users", user.uid);

    await updateDoc(userRef, {
      personas: arrayRemove(persona),
    });

    if (activePersona && activePersona.name === persona.name) {
      setActivePersona(null);
    }
  };

  return (
    <div>
      <h2>Your Personas</h2>
      <ul>
        {personas.map((persona, index) => (
          <li key={index}>
            {persona.name} - {persona.type} - {persona.difficulty}
            <button onClick={() => setActivePersona(persona)}>
              Set Active
            </button>
            <button onClick={() => handleDelete(persona)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PersonaList;
