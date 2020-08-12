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
    maxWidth: 345,
  },
});

export default function ImgMediaCard() {
  const [file, setFile] = useState("");
  const [currImage, setCurrImage] = useState(
    process.env.PUBLIC_URL + "/boat.jpeg"
  );
  console.log("ImgMediaCard -> currImage", currImage);
  const classes = useStyles();

  const fileSelectedHandler = (event) => {
    console.log("ImgMediaCard -> event", event.target.files[0]);
    const fileObject = event.target.files[0];
    setFile(fileObject.name);
    setCurrImage(`${process.env.PUBLIC_URL}/${fileObject.name}`);
  };

  const uploadImage = async () => {
    const url = `${apiURL}/enhance`;
    const body = { image: file };
    try {
      const response = await postCall(url, body);
      const payload = await response.json();
      console.log("uploadImage -> payload", payload);
    } catch (err) {
      console.log("error in enhancing image", err);
    }

    console.log("uploadImage -> url", url);
    console.log("uploadImage -> image", file);
  };

  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardMedia
          component="img"
          alt="Image to be enhanced"
          height="140"
          image={currImage}
          title="Image to be enhanced"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            Please select files in the images folder of the project root
            directory.
          </Typography>
          <input type="file" onChange={fileSelectedHandler} />
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button size="small" color="primary" onClick={uploadImage}>
          Upload
        </Button>
        <Button size="small" color="primary">
          Download
        </Button>
      </CardActions>
    </Card>
  );
}
