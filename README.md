# uptime-girl

An Uptime Robot robot that automatically creates monitors based on container labels.

## Proposed

This app will make use of Docker container labels to create [Uptime Robot](https://uptimerobot.com) monitors using the [API](https://uptimerobot.com/api).

It will be designed to be spun up with zero configuration (apart from setting the API key, obviously) and left to periodically check containers for changes to labels.

## Why?

I got sick of constantly adding/removing monitors from the Uptime Robot web UI every time I deploy or end-of-life a website or API service. Since I already run everything in Docker and already make use of a couple of apps that use labels to automate things, I thought why not automate this too!
