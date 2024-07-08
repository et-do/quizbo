import React from "react";

function SelectionPage({ setPage }) {
  return (
    <div className="selection-options">
      <button className="selection-button" onClick={() => setPage("quizForm")}>
        Webpage
      </button>
      <button className="selection-button" disabled>
        PDF
      </button>
      <button className="selection-button" disabled>
        Audio
      </button>
    </div>
  );
}

export default SelectionPage;
