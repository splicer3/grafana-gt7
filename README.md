# Gran Turismo 7 Telemetry with Grafana

This **Grafana data source plugin** allows for visualization of telemetry data sent by GT7 over the network in a broadcast fashion, using the UDP Port 33740.

This project is derived from a fork I made some time ago from [this project by Alexander Zobnin](https://github.com/splicer3/grafana-gt7), who did an amazing work to create a universal data source for simracing titles. Unfortunately, my goal is **incompatible with titles such as ACC and iRacing** (which have memory-mapped local telemetry files), so I decided to keep the code I developed for usage with GT7 in [this fork](https://github.com/splicer3/grafana-gt7) and start over to keep it simpler.




https://github.com/user-attachments/assets/97d8c7cb-145c-4206-be59-c4dc09d39473


## Goals
The **main goal** is to create a Docker Compose file with as little needed know-how as possible that setups a Grafana instance with automatically provisioned dashboards and data sources, that can run on a Raspberry Pi 2B (which is what I'm using for testing).

This would allow the project to run practically anywhere with as little setup as possible, bar for the Playstation's local IP.

The current to-do list is as follows:
- EV support (would be helpful to find a way to change units dynamically in Grafana)
- Additional flags decoding and visualisation (like TCS, ABS flags)
- A better lap implementation overall (like better conversion of laptimes)
- A smarter dashboard that can use the CarID and the maximum revs sent by GT7

## Features

- Real-time lightweight telemetry data visualization
- Highly customizable dashboard, but with a default one provisioned at Docker Compose startup.
- make-release shell script to create a zip file with everything needed to run it on Docker.
- Playstation's IP editable through Grafana data source options

## Supported titles

Gran Turismo 7 is the only supported title for now. Another project might be created for other games once this is finished.

## Supported platforms
Literally anything that can run Docker. Alternatively, make-release.sh can generate files needed for the following platforms:
- Windows
- Linux
- LinuxARM
- LinuxARM64

## Getting started
1. Make sure you have **npm** and **Go** installed and available, and run `npm install` in the root project folder.
2. [Get started with Docker](https://www.docker.com/get-started/) and make sure you know [how to run Compose](https://docs.docker.com/compose/).
3. Clone this repo and run make-release.sh
4. Choose the OS you need (if you're running this on Docker, it will be one of the Linux options. Generally speaking, it's Linux for normal environments, LinuxARM for old Raspberry Pis and LinuxARM64 for newer Raspberry Pis and Docker running on M series Mac systems).
5. Wait for the script to finish. The zip file containing the required files will be generated in the project's root.
6. Unzip the file in any directory on your target machine.
7. Run `docker compose up -d`, or `docker compose up -d --build` if you made any changes to the code.
8. Connect to `localhost:3000`, and log `username: admin` and `password: admin` as always with Grafana.
9. Go to data source options, find Gran Turismo 7 Telemetry and change the Playstation IP field to your own Playstation's IP
10. Go to dashboards and either build one from scratch or use the default provisioned one.


## Credits
**Alexander Zobnin** for creating the [original simracing telemetry plugin for Grafana](https://github.com/splicer3/grafana-gt7). 

**Nenkai** for his work on GT7 telemetry raw data and decoding. 

**Matthias Küch** for his work on the [GT7 Python telemetry software](https://github.com/snipem/gt7dashboard), which I largely used to understand how to decrypt the incoming GT7 packets.

**[tarnheld](https://www.gtplanet.net/forum/members/tarnheld.265292/)** for [his post on GTPlanet](https://www.gtplanet.net/forum/threads/gt7-is-compatible-with-motion-rig.410728/post-13829730) detailing how to convert GT7's quaternions to euler angles.
