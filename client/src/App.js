import React from "react";
import "./App.css";
import AppBar from "./components/AppBar";
import UploadCard from "./components/UploadCard";
import Typography from "@material-ui/core/Typography";

function App() {
  return (
    <div className="App">
      <AppBar />
      <div className="title">
        <Typography variant="h4" color="inherit">
          {" "}
          Lucid Images
        </Typography>
        <Typography variant="h6" color="inherit">
          {" "}
          Upload an image for histogram equalization
        </Typography>
      </div>
      <div className="App-upload-card">
        <UploadCard />
      </div>
    </div>
  );
}

export default App;
