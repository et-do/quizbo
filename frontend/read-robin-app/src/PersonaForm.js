import React, { useState } from "react";
import { doc, updateDoc, arrayUnion } from "firebase/firestore";
import { db } from "./firebase";

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
    <form onSubmit={handleSubmit}>
      <h2>Create a Persona</h2>
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
      <div>
        <label>
          User Type:
          <input
            type="text"
            value={userType}
            onChange={(e) => setUserType(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Difficulty Level:
          <input
            type="text"
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value)}
            required
          />
        </label>
      </div>
      <button type="submit">Add Persona</button>
    </form>
  );
};

export default PersonaForm;
