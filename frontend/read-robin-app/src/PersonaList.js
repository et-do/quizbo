import React, { useEffect, useState } from "react";
import { doc, deleteDoc, collection, getDocs } from "firebase/firestore";
import { db } from "./firebase";
import "./PersonaList.css";

const PersonaList = ({ user, activePersona, setActivePersona }) => {
  const [personas, setPersonas] = useState([]);

  useEffect(() => {
    const fetchPersonas = async () => {
      if (!user) return;

      const personaCollection = collection(db, "users", user.uid, "personas");
      const personaSnapshot = await getDocs(personaCollection);

      const personaList = personaSnapshot.docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
      }));
      setPersonas(personaList);
    };

    fetchPersonas();
  }, [user]);

  const handleDelete = async (persona) => {
    if (!user) return;

    const personaRef = doc(db, "users", user.uid, "personas", persona.id);
    await deleteDoc(personaRef);

    setPersonas((prevPersonas) =>
      prevPersonas.filter((p) => p.id !== persona.id)
    );

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
              <strong>{persona.name}</strong> - {persona.role} -{" "}
              {persona.difficulty} - {persona.language}
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
