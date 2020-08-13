import React, { Fragment } from "react";
import "./App.css";
import AppBar from "./components/AppBar";
import UploadCard from "./components/UploadCard";
import Typography from "@material-ui/core/Typography";
import ImageCard from "./components/ImageCard";

const enhancedImage = localStorage.getItem("enhancedImage");
console.log("enhancedImage", enhancedImage);

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
          Upload an image that needs contrast improvement
        </Typography>
      </div>
      <div className="App-upload-card">
        <UploadCard />
      </div>
      {enhancedImage ? (
        <Fragment>
          <div className="App-enhanced-image-title">
            {" "}
            <Typography variant="h4" color="inherit">
              {" "}
              Contrast Enhancement Results
            </Typography>{" "}
          </div>
          <div className="App-enhanced-image-container">
            <ImageCard image={enhancedImage} format={"RGB"} isOriginal={true} />
            <ImageCard image={enhancedImage} format={"RGB"} />
            <ImageCard image={enhancedImage} format={"YUV"} />
            <ImageCard image={enhancedImage} format={"HSL"} />
          </div>
        </Fragment>
      ) : null}
    </div>
  );
}

export default App;
