import React, { useState } from "react";
import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardActions from "@material-ui/core/CardActions";
import CardContent from "@material-ui/core/CardContent";
import CardMedia from "@material-ui/core/CardMedia";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import "./index.scss";

const useStyles = makeStyles({
  root: {
    maxWidth: 345,
  },
});

export default function ImgMediaCard() {
  const [file, setFile] = useState("");
  const [currImage, setCurrImage] = useState("");
  const classes = useStyles();

  const fileSelectedHandler = (event) => {
    console.log("ImgMediaCard -> event", event.target.files[0]);
    const file = event.target.files[0];
    setFile(file.name);
  };

  const uploadImage = () => {
    const basePath = "/home/vagrant/project/images";
    const fullPathToImage = `${basePath}/${file}`;
    setCurrImage(fullPathToImage);
    console.log("uploadImage -> fullPathToImage", fullPathToImage);
  };

  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardMedia
          component="img"
          alt="Contemplative Reptile"
          height="140"
          image={currImage}
          title="Contemplative Reptile"
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
