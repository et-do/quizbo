import React from "react";
import "./SelectionPage.css";

function SelectionPage({ setPage }) {
  const handleSelectionChange = (event) => {
    const page = event.target.value;
    if (page) {
      setPage(page);
    }
  };

  return (
    <div className="selection-page-custom">
      <div className="selection-options-custom">
        <h1>What do you want a quiz generated for?</h1>
        <select
          className="selection-dropdown-custom"
          onChange={handleSelectionChange}
          defaultValue=""
        >
          <option value="" disabled>
            Select an option...
          </option>
          <option value="urlForm">🌐 Webpage</option>
          <option value="pdfForm">📄 PDF</option>
          <option value="audioForm">🎧 Audio</option>
          <option value="videoForm">🎥 Video</option>
          <option value="textForm">📝 Text</option>
        </select>
      </div>
    </div>
  );
}

export default SelectionPage;
