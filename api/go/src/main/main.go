package main

import (
    "../image"
    "fmt"
)

func main() {
    fmt.Println("Running enhancement!")
    fileName := "boat.jpeg"
    // jpeg images
    format, err := image.GuessImageFormat(fileName, true)
    if err != nil {
        panic(err)
    }
    fmt.Println(format)

    image.YUVHistogramEquilization(fileName, true)
    image.RGBHistogramEquilization(fileName, true)
    image.HSLHistogramEquilization(fileName, true)

    // png images

}