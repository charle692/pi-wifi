# pi-wifi
IOT like wifi configuration for the Raspberry pi

## Installation
First off, your home directory will need to be `/home/pi`.
Next you will need to clone the project from github using `git clone https://github.com/charle692/pi-wifi.git`.
Next move into the project dir: `cd pi_wifi`.
And finally run `sudo sh setup.sh`.
You will be asked to reboot for changes to take affect. Once rebooted, you should see a network called `pi-wifi`.

## Usage
1. Once the Raspberry Pi has rebooted, an access point called `pi-wifi` will be created
2. Connect to `pi-wifi` using your device of choice
3. Navigate to `10.0.0.1:3001/views/index`
4. Once the page is loaded wait for the select field to populate with the available networks
5. Fill out the form and submit
6. Once the Raspberry Pi has been rebooted, it will connect to the specified network

