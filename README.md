# pi-wifi
IOT like wifi configuration for the Raspberry pi

## Installation
First you will need to install the code into the home directory of the raspberry pi. 
Move all of the files and folders outside of the project dir into /home/pi.

Next run `sudo sh setup.sh`

## Usage

1. Once the Raspberry Pi has rebooted, an access point called `pi-wifi` will be created
2. Connect to `pi-wifi` using your device of choice
3. Navigate to `10.0.0.1:3001/views/index`
4. Once the page is loaded wait for the select field to populate with the available networks
5. Fill out the form and submit
6. Once the Raspberry Pi has been rebooted, it will connect to the specified network

