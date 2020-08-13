import React, { useState } from "react";
import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardActions from "@material-ui/core/CardActions";
import CardContent from "@material-ui/core/CardContent";
import CardMedia from "@material-ui/core/CardMedia";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import { apiURL } from "../../config";
import { postCall } from "../../apiCalls";
import "./index.scss";

const useStyles = makeStyles({
  root: {
    maxWidth: 500,
  },
  btnContainer: {
    display: "flex",
    justifyContent: "center",
  },
});

export default function ImgMediaCard(props) {
  const { image } = props;
  const [file, setFile] = useState("");
  console.log("ImgMediaCard -> file", file);
  const [currImage, setCurrImage] = useState(
    process.env.PUBLIC_URL + `/${image}`
  );
  console.log("ImgMediaCard -> currImage", currImage);
  const classes = useStyles();

  const fileSelectedHandler = (event) => {
    console.log("ImgMediaCard -> event", event.target.files[0]);
    const fileObject = event.target.files[0];
    setFile(fileObject.name);
    setCurrImage(`${process.env.PUBLIC_URL}/${fileObject.name}`);

    // clear enhanced image
    localStorage.removeItem("enhancedImage");
  };

  const uploadImage = async () => {
    const url = `${apiURL}/enhance`;
    const body = { image: file };
    if (
      file &&
      (file.includes(".png") || file.includes(".jpg") || file.includes(".jpeg"))
    ) {
      try {
        localStorage.setItem("enhancedImage", file);
        const response = await postCall(url, body);
        const payload = await response.json();
        console.log("uploadImage -> payload", payload);

        // set enhancedImage in localstorage
      } catch (err) {
        // localStorage.removeItem("enhancedImage");
        console.log("error in enhancing image", err);
      }
    } else {
      alert("Please select a png / jpeg file");
    }
  };

  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardMedia
          component="img"
          alt="Image to be enhanced"
          height="300"
          image={currImage}
          title="Image to be enhanced"
          className="uploadcard-original-img"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            Please select images from the $root/client/public folder. Feel free
            to add your own jpeg/png images.
          </Typography>
          <input type="file" onChange={fileSelectedHandler} />
        </CardContent>
      </CardActionArea>
      <CardActions className={classes.btnContainer}>
        <Button size="small" color="primary" onClick={uploadImage}>
          Upload
        </Button>
      </CardActions>
    </Card>
  );
}

ImgMediaCard.defaultProps = {
  image: "welcome.jpg",
};
