# Lucid Luminator - Polyglot project

### Overall goal of project

What is the overall goal of the project (i.e. what does it do, or what problem is it solving)?

Lucid Luminator is a simple web application that uses several different color models and
perceptual attributes to perform contrast enhancement. This is accomplished specifically
by using histogram equalization.

Histogram Equalization is a common approach for enhancing contrast and brightness with grayscale
images. Extending this to color images is not a straight forward task.

I chose to apply this technique to 3 different color spaces, RGB, YUV, and HSL.
Each model leverages a different perceptual attribute, and has a different effect on contrast enhancement.

This is often tied to the pixel make up / distribution of the source image. HSL leverages
lightness (visual perception of perceived brightness of an area), YUV is tied closely to brightness.
By increasing or decreasing Luma, chroma saturation levels are affected.

Initially, the goal of the application was to implement a single histogram equalization algorithm on
a single color model. However, the color model that I chose was RGB. I found quickly that this was a
poor color space choice. This is documented in Strickland, Kim and Mcdonnell 's paper on color
enhancement in 1987. Histogram equalization in the RGB color space leaves color artifacts. This
is because of its poor correlation with the human visual system, requiring a color space transformation. This lead me to implement an additional two algorithms for color spaces, HSL and YUV.

The algorithm at a high level is rather simple. Convert the image from RGB color space to
YUV | HSL, perform histogram equalization and convert back to the RGB color space. The difficulty
was in the details. This was especially painful for HSL (Go's color package does not have a
built in conversion method like it does for YUV :( ). HSL histogram equalization was performed
on saturation for each hue in the image taking into account the maximum allowable saturation
for a given luminence. This makes the saturation a function of both hue and luminance. This
technique was developed by Hague, Weeks and Myler[Hague 1994]. It helps reduce the number of
out-of-gamut colors (colors that don't have a direct conversion from one gamut to another)
and color artifacts.

The result is 3 different contrast enhancement techniques for your visual pleasure.

### Languages / parts of systems were implemented with what

I chose to use a fast system language for image processing (contrast enhancement). Go
was a natural choice because of it's native support for concurrency. The second part to
these contrast enhancement algorithms was to write concurrent versions of them. There
is often ample opportunity to make speed improvements to image processing algorithms
due to their embarrassingly parallel workload. Histogram enhancement is no exception.
I added concurrency wherever it made sense to add. This included histogram generation,
image reading, color space conversion, image writing, etc.

For displaying the results of contrast enhancement, I chose to implement a web application
with HTML, CSS, and Javascript on the front end. I used React as the front end web framework.
This is because I thought it would be nice to interact with the application through a browser.

To glue these two components, I implemented a simple flask python web server. The web server
has an extremely simple API meant to handle requests from the applications front end and pass
it along to go image processing service where all the computational heavy lifting occurs.

### Methods of communication

There are two main methods of communication between the 3 different language modules. The first
is with a RPC REST server. The browser makes HTTP requests on behalf of my front end code to
a python web server. The python web server is running a RabbitMQ message queue client. The web
server packages this request and places it onto the RabbitMQ queue to the Go image processing
server. From here, the Go image web server accepts the incoming message and performs the desired
computation. Once finished, the web server returns whether or not the request was processed successfully. The images are written into React's public directory. This is so React can render
the images (React is not allowed to access files outside of its src directory expect files in
the public folder).

### Steps to run application

1. run `vagrant up`
