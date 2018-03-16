# ADZ

Adz is a small tool for showing some quick stats fro Postfix events from a Postfix log file. 

# Usage

See `adz -h` for details. Quick summary: `-a AGEINMINUTES` to set the number of
minutes old an event is to be counted; default is 1 minute. Use `-f FILENAME`
to specify a file other than `/var/log/maillog`.o

# Stats

More to be added, but for now here is what we have.

## Connections

Connects reports the number of 'connect from' events, disconnects the number of
'disconnect from', and 'lost connections' is the count of times 'lost
connection' occurs on an SMTPd event.

## Delivery Stats

Reporting Bounces, Deferrals, and Sends. These are taken from `status=` events.
A message can be responsible for more than one event. For example if a message
comes in that gets deferred due to remote relay issues each deferral is
counted. If the message eventually times out it will also record a bounce, and
if it is delivered it will record a send.

## Local Pickups

This number indicates the number of times a message was picked up from the
local pickup queue. This can occur when queue files are migrated or when a
local invocation of the `sendmail` command is used.

# Dependencies

THe `chisel` library is used for parsing the log file. Oh, and you need the
Postfix log file you want to analyze.

# TODO

* Add delay metrics ?
* Add relay metrics ?

Much of these can be added once the chisel library supports them.
