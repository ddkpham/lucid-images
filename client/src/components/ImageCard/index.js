import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardActions from "@material-ui/core/CardActions";
import CardContent from "@material-ui/core/CardContent";
import CardMedia from "@material-ui/core/CardMedia";
import Typography from "@material-ui/core/Typography";

const useStyles = makeStyles({
  root: {
    width: "40vw",
    margin: "10px",
  },
  media: {
    height: 300,
    backgroundSize: "contain",
  },
});

export default function ImageCard(props) {
  const classes = useStyles();
  const { image, format, isOriginal } = props;

  const img = isOriginal ? image : `enhanced-${format}-${image}`;
  return (
    <Card className={classes.root}>
      <CardContent>
        <Typography gutterBottom variant="h5" component="h2">
          {image.split(".")[0]}
        </Typography>
        <Typography variant="body2" color="textSecondary" component="p">
          {isOriginal ? "Original" : `${format} histogram equalization`}
        </Typography>
      </CardContent>
      <CardActionArea>
        <CardMedia
          className={classes.media}
          image={process.env.PUBLIC_URL + `/${img}`}
          title="Image"
        />
      </CardActionArea>
      <CardActions></CardActions>
    </Card>
  );
}
