# noisemaker

Simple load testing tool that could be run as a lightweight docker container or standalone binary. 
This will probably not work very well on machines with a single core cpu (untested). 

Also testing this while approaching limits may behave unexpectedly as this was not tested much. The functions
that simulate load are fairly primitive (especially the CPU one), but looking at resource monitors it appears
to work.

## Basic usage

### Container

Should work with the image alone:

```sh
podman run ghcr.io/madhuravius/noisemaker:v0.1.3 run
```

### CLI

Highlighted by just running the application directly, which can be downloaded
on the [releases page](https://github.com/madhuravius/noisemaker/releases):

```sh
NAME:
   noisemaker - needlessly consume resources and throw it in the bin

USAGE:
   noisemaker [global options] command [command options] [arguments...]

COMMANDS:
   run,     start the trashcan
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cpu value        as a percentage, specify the percentage of CPU you would like to use. Max: 99 (default: 50)
   --mem value        as a percentage, specify the percentage of RAM you would like to use. Max: 99 (default: 50)
   --bandwidth value  in MBps, specify how much bandwidth you want to use upload/download (default: 5)
   --port value       specify port you wish to run the web server for stressing (default: 3000)
   --help, -h         show help
```

## Why?

I wanted to write a simple docker image I could use to simulate CPU, memory, disk, 
and network load. I needed something very simple/readable and controllable for testing
vms and computers at home.

There are probably better tools for this job, but nothing quite fit what I needed in a catch-all:

* [stress-ng](https://github.com/ColinIanKing/stress-ng) and [s-tui](https://github.com/amanusk/s-tui)
* [speedtest](https://www.speedtest.net/apps/cli)
