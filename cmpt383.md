# Lucid Luminator - Polyglot project

### Overall goal of project

Lucid Luminator is a simple web application that uses several different color models/spaces and
perceptual attributes to perform contrast enhancement. Contrast enhancement is a technique used
to increase the visibility between two adjacent objects or reveal richer detail in an image.
There are many techniques to do this. Lucid Luminator performs contrast enhancement by
using histogram equalization.

Histogram Equalization is a common approach for enhancing contrast and brightness with grayscale
images. Extending this to color images is not a straightforward task.

I chose to apply this technique to 3 different color spaces, RGB, YUV, and HSL.
Each model leverages a different perceptual attribute, and has a different effect on contrast enhancement.
This is often tied to the pixel make up / distribution of the source image. HSL leverages
lightness (visual perception of perceived brightness of an area), YUV is tied closely to brightness.
By increasing or decreasing Luma, chroma saturation levels are affected.

Initially, the goal of the application was to implement a single histogram equalization algorithm on
a single color model. However, the color model that I chose was RGB. I found that this was a
poor color space choice. This is documented in Strickland, Kim and Mcdonnell 's paper on color
enhancement in 1987. Histogram equalization in the RGB color space leaves color artifacts. This
is because of its poor correlation with the human visual system, requiring a color space transformation.
This lead me to implement an additional two algorithms for color spaces, HSL and YUV.

The algorithm at a high level is rather simple. Convert the source image from RGB color space to
YUV | HSL, perform histogram equalization and convert back to the RGB color space. The difficulty
was in the details. This was especially painful for HSL (Go's color package does not have a
built in conversion method like it does for YUV :( ). HSL histogram equalization was performed
on saturation for each hue in the image taking into account the maximum allowable saturation
for a given luminance. This makes the saturation a function of both hue and luminance. This
technique was developed by Hague, Weeks and Myler[Hague 1994]. It helps reduce the number of
out-of-gamut colors (colors that don't have a direct conversion from one gamut to another)
and color artifacts.

The result is 3 different contrast enhancement techniques for your visual pleasure. How well
a technique / color model performs is often tied to the make up of the original source image.
With UI/UX in mind, I chose to implement all 3 and display them all at once so that a user is
able to see strengths and weaknesses of each color space as well as get a feel for which
technique may better suit the original poorly contrasted image.

### Languages / parts of systems were implemented with what

I chose to use a fast system language for image processing (contrast enhancement). Go
was a natural choice because of its native support for concurrency. The second challenge to
these contrast enhancement algorithms was to write concurrent versions of them. There
is often ample opportunity to make speed improvements to image processing algorithms
due to their embarrassingly parallel workload. Histogram enhancement is no exception.
I added concurrency wherever it made sense to add. This included histogram generation,
image reading, color space conversion, image writing, etc. The end result is an incredibly
fast image enhancement algorithm.

For displaying the results of contrast enhancement, I chose to implement a web application
with HTML, CSS, and Javascript on the front end. I used React as the front end web framework.
This is because I thought it would be nice to interact with the application through a browser.

To glue these two components, I implemented a simple flask python web server. The web server
has an extremely simple API meant to handle requests from the applications front end and pass
it along to go image processing service where all the computational heavy lifting occurs.

### Methods of communication

There are two main methods of communication between the 3 different language modules. The first
is with a `RPC REST server`. The browser makes HTTP requests on behalf of the front end code to
a python web server. The python web server is running a `RabbitMQ message queue` client. The web
server packages this request and places it onto the RabbitMQ queue to the Go image processing
server. From here, the Go image web server accepts the incoming message and performs the desired
computation. Once finished, the web server returns whether or not the request was processed successfully.
The images are written into React's public directory. This is so React can render
the images (React is not allowed to access files outside of its src directory expect files in
the public folder). The overall system architecture design is a web-queue worker.

Webqueue-worker definition:

https://docs.microsoft.com/en-us/azure/architecture/guide/architecture-styles/web-queue-worker

### Steps to run application

The application uses ports `5000` & `8080`. Make sure those ports are free to run the application.
If they are not free, you may change the ports on the host machine to spin up the VM. This requires
going into the vagrant file and changing the port forwarding configuration.

1. Start up VM and provision

`vagrant up`

This should install all dependencies, start up web server for static files, start up python
API, and go image processing RPC server.

2. Navigate to http://localhost:8080/ (May take a second to load initially)

3. Click "choose file" and choose an image from the folder \$projectroot/client/public

For me on a mac, the folder location looks something like

`/Users/davidpham/Documents/cmpt383/polyglot-project/client/public`

For convenience, I have uploaded a couple images. Ones that I found
particular interesting and highlighted some of the strengths and weaknesses
of histogram equalization on different color models were:

- backyard.jpeg
- boat.jpeg
- lena.png
- couple.jpeg
- argument.png

4. Click upload image

5. All 3 Contrast Enhancements should be shown below for your viewing pleasure

6. Add some of your own images (jpeg / png) to try contrast enhancement

make sure that they are added to the \$projectroot/client/public folder
so that React is able to display them.

7. If you like, you can also go into \$projectroot/api/go/src/image/histogram.go
   There is a global constant variable called `numThreads` which is the number of
   threads for the routine. Feel free to change this to play with the system. However,
   you must stop and restart the RPC client / server to see the changes. You must also
   start the server and client in the same order. Go RPC server -> Python client web server.

### DEBUGGING STEPS

If you are having trouble running the application, try these steps. You may have to run the application
manually. Feel free to email me if you are having issues `dpa35@sfu.ca`. You may have to start the go
server and python server manually. PM2 is managing my front end server so we do not have to worry about that.

comment out lines 178 - 191 of chef script in `cookbooks/polyglot/recipes/default.rb`. These are lines
associated with running both the go server and python server.

Re-provision VM.

`vagrant provision`

`vagrant ssh`

1. Start up Go RabbitMQ server. If successful, (You should see "Awaiting RPC requests") go to step 4.
   (Optional) If running server fails due to missing ampq and webp libraries, run step 3.

`cd /home/vagrant/project/api/go/src/main`
`go run amqp_server.go`

2. Install libraries with bash script in same folder (Skip if server is able to run)

`bash install`
`go run amqp_server.go`

3. Start up Flask server and client rabbitMQ

open another terminal and ssh in.
`vagrant ssh`
`cd /home/vagrant/project/api`
`python3 app.py`

4. Navigate to http://localhost:8080/
