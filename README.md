# dd-logstats: monitor your logs

[![Build Status](https://travis-ci.org/j-vizcaino/dd-logstats.svg)](https://travis-ci.org/j-vizcaino/dd-logstats)

`dd-logstats` is a simple program for parsing W3C Common Log lines and displaying statistics and alarms about the traffic.

## Features

`dd-logstats` does the following:

* reads log lines from `stdin`
* aggregates statistics according to the section of the URL
* raises alarms on the frontend when traffic is too high 

## Installation

### Dependencies

In order to build `dd-logstats` you need:

* a working Go (>= 1.6) installation
* `make`
* [Glide](https://glide.sh) (>= 0.12)


### Building

```shell
$ git clone https://github.com/j-vizcaino/dd-logstats $GOPATH/src/dd-logstats
$ cd $GOPATH/src/dd-logstats
$ make
```

## Usage

The easiest way to feed `dd-logstats` is by using the `tail` program:

* start the program

```shell
$ tail -F /var/log/nginx/access.log | ./dd-logstats
```

* open a browser to visit [http://localhost:8080/](http://localhost:8080/)


### Options

Full argument options can be printed starting the program with `-help`.

Most important options:

* `-serve`: specifies the server listen address. Example: `:12345`
* `-alarm-threshold`: sets the alarm threshold (average hit for the given time frame)
* `-alarm-timeframe`: time frame for the average hit alarm. Example: `2m0s`
* `-stats-period`: statistics aggregation period. Example: `10s`

