# youtube-livestream-merge

Lets say that you have a youtube livestream running on your desktop and you want to leave the house, create a 2nd stream, and then load the 2nd stream to play in a browser on the 1st stream running on the desktop. Why do this automatically? Well, I live in a tall building with an elivator that kills phone service. I can't. It has to be automated.

Please note, a lot of this code was not written by me. It is example code. FYI.


## HOW TO RUN

you need to generate a client_secret.json

to do this, go to the google youtube API credentials and create a credential for a Desktop app. Then download the client_secret_garblygoookasdfasdfweedasdf.json and rename it client_secret.json in the root of this project

then

`go run src/main.go`
