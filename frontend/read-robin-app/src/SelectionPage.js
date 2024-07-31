import React, { useState } from "react";
import Select from "react-select";
import "./SelectionPage.css";

const options = [
  { value: "urlForm", label: "🌐 Webpage" },
  { value: "pdfForm", label: "📄 PDF" },
  { value: "audioForm", label: "🎧 Audio" },
  { value: "videoForm", label: "🎥 Video" },
  { value: "textForm", label: "📝 Text" },
];

function SelectionPage({ setPage }) {
  const [selectedOption, setSelectedOption] = useState(null);

  const handleChange = (option) => {
    setSelectedOption(option);
    setPage(option.value);
  };

  return (
    <div className="selection-page">
      <div className="selection-options">
        <h1>What do you want a quiz generated for?</h1>
        <Select
          value={selectedOption}
          onChange={handleChange}
          options={options}
          className="selection-dropdown"
          placeholder="Select an option..."
        />
      </div>
    </div>
  );
}

export default SelectionPage;
