package main

import (
    "../image"
    "fmt"
)

func main() {
    fmt.Println("Hello World!")
    fileName := "train.jpg"
    // jpeg images
    format, err := image.GuessImageFormat(fileName)
    if err != nil {
        panic(err)
    }
    fmt.Println(format)

    image.YUVHistogramEquilization(fileName)

    // png images

}