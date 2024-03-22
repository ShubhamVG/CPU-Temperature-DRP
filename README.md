# CPU-Temperature-DRP
Discord Rich Presence which adds a rich presence activity to my Discord account telling my CPU temperature.
The only reason I made this is because Dana wanted to see my CPU suffer every time I write and run bad code on my device :p

## Made & Tested on
Made on Linux Mint for Linux.
### How it works?
It works by parsing the data from data from `lm-sensors` (which should be installed on most if not every linux distro). Run ```sensors``` in the terminal to check.
Everyone might have different sensors so run ```sensors -j``` in the terminal to check your sensors and add them to the main.go's `SENSORS_TO_USE` list.