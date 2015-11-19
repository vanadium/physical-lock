# Physical lock

This project defines the software for building a secure physical lock
using the Vanadium stack, along with a commmand-line tool for interacting
with the lock. The software runs on a Raspberry Pi and interacts with the
locks's switches and sensors using GPIO. It runs a Vanadium RPC service
that allows clients to send lock and unlock requests.

Before describing the design in detail we highlight its key distinguishing
aspects:

* Decentralized: There is no single authority on all the locks, no cloud server
  that controls access to all locks from a particular manufacturer. All secrets
  and credentials pertaining to a particular lock are solely held by the lock and
  its clients. Huge compute farms run by hackers all over the world have no single
  point of attack that can compromise the security of multiple locks. 

* No Internet Connectivity Required: The lock does not require internet connectivity
  to function. When you’re right in front of the device, you can communicate with
  it directly without going through the cloud or a third-party service.

* Audited: The lock can keep track of who opened the door, when and how they got
  access.

# Design

An out-of-box lock device exposes the `UnclaimedLock` interface on startup.

```
type UnclaimedLock interface {
     Claim(name string) (security.WireBlessings | error)
}
``` 

This interface has a single method `Claim` that can be used to claim the
device with a specific name and obtain credentials that grant authorization
for subsequently interacting with the lock. Once claimed, the device exposes
the `Lock` interface.

```
type Lock interface {
     // Lock locks the lock.
     Lock() error
     // Unlock unlocks the lock.
     Unlock() error
     // Status returns the current status (locked or unlocked) of the
     // lock.
     Status() (LockStatus | error)
}
```

Clients that possess the appropriate authorization credentials can interact
with the device and send it requests to lock, unlock or determine the status.

## Security Model

An out-of-box lock device comes with pre-installed credentials (private key and
blessings) from its manufacturer. For instance, it may have blessings of the form
`<lock manufacturer>:<serial no>`. The `UnclaimedLock` interface exposed by it
is accessible to everyone.

When the lock is claimed with a specific name, it blesses itseld under the name
and uses that blessing to subsequently authenticate to clients. It also creates
a blessing for the claimer by extending this (self) blessing with the extension
`key`. For e.g., a lock claimed with the name `front-door` will subsequently 
authenticate with the blessing `front-door`, and grant the blessing `front-door:key`
to the claimer.

The blessing granted to the claimer is the _key_ to the lock. It provides access to
the methods on the `Lock` interface exposed by a claimed lock device. Furthermore, all
extensions of this _key_ blessing are also authorized to access the `Lock` interface.

Access to the lock can be shared with other principals by blessing them with the
_key_ blessing. Appropriate caveats may be added to this blessing to limit the scope
of use. For e.g., the key `front-door:key` can be shared with a house guest for a
limited period by extending it as `front-door:key:houseguest` with an expiration caveat.

## Device Discovery

Both claimed and unclaimed lock devices advertise their endpoints using [MDNS] to allow
clients to discover them without going through the cloud or a third-party Mounttable.
All lock devices advertise themselves with names that are prefixed with `lock-`. This
prefix allows clients to easily filter MDNS advertisements for lock devices.

An unclaimed lock device advertises under a name of the form `lock-unclaimed-lock-xxxxxx`
where `xxxxxx` is a random number generated by the lock device. A claimed lock device
advertises itself under the name `lock-<claimed name>`.

Currently we don't have any protection for preventing a rogue device from impersonating
the MDNS name of a lock device. This is something we plan to add in the future. Note
that impersonating the MDNS name of a lock device may lead to denial of service but
does not compromise the security of the lock.

## Locking and Unlocking

The `Lock` interface has methods to lock, unlock and determine the status of the lock
device. The device logs the blessings of the caller for all method requests thereby
maintaining an audit trail. Internally, these methods set or unset certain GPIO pins on the Raspberry Pi 
based on the circuitory described below.

### Equipment
- 10KΩ resistor
- 1KΩ resistor
- Magnetic switch (normally open circuit - closed when the sensor moves away)
  (e.g. [P/N 8601 Magnetic Switch](http://www.amazon.com/gp/product/B0009SUF08/ref=oh_aui_detailpage_o03_s00?ie=UTF8&psc=1))
- For now, an active buzzer to simulate a relay.
  Will fill in the relay details here once we have that setup.
- Breadboard, jumper cables, ribbon cable - or whatever is needed to connect to the RPi's GPIO pins

### Circuitry

Apologies for this unconventional, possibly unreadable circuit representation. Putting it down
so that the author can remember! TODO(ashankar): Clean it up!

The pin number assignments here work both for RaspberryPi Model B/B+ and RaspberryPi2-ModelB.

```
---(Pin 1)-------/\/\(10KΩ)/\/\---------(COM port of magnetic switch)
                                  \
                                   \----/\/\(1KΩ)/\/\---------(Pin 15 = GPIO22)
                                                          \
                                                           \----(LED)-----|
                                                                          |
                                                                          |
                                          (N.O. port of magnetic switch)--|
                                                                          |
                                                                          |
                                         (-ve terminal of active buzzer)--|
                                                                          |
                                                                          |
                                                                          |
                                                           (Pin 6 = GND)--|

---(Pin 11 = GPIO17)-----------(+ terminal of active buzzer)
```

The above ciruit is meant to be a simulation of an actual lock device wherein
locking and unlocking simply makes a buzzer ring. In particular the `Lock` and
`Unlock` calls update the status of the pin `GPIO17`, and the `Status` call
checks the status of the ping `GPIO22`.

# Deployment

Building the lock service for the RaspberryPi

```
jiri go get -u github.com/davecheney/gpio
JIRI_PROFILE=arm jiri go install v.io/x/lock/lockd
scp $JIRI_ROOT/release/projects/physical-lock/go/bin/lockd <rpi_scp_location>
```

If building without the `arm` profile, there are no physical switches/relays
and instead a simulated hardware is used that uses the interrupt signal (SIGINT)
to simulate locking/unlocking externally.

The lock service can be started by running the following command in the
directory where the `lockd` binary was copied.

```
lockd --v23.credentials=<creds dir> --config-dir=<config dir>
```

The `--config-dir` flag specifies a path to a directory where configuration files 
can be saved, and the `--v23.credentials` flag specifies a path to the Vanadium
credentials directory for `lockd`. The configuration directory must be persisted
accross restarts; emptying it would amount to a "factory reset".

# Sample Usage

We describe a few commands for discovering and interacting with
a lock device.

The first step is to build the command-line tool.

```
jiri go install v.io/x/lock/lock
```

In the rest of this section, we specify all commands
 relative to the directory containing the `lock` tool.

All commands must be run under an [agent] that has a blessing from
`dev.v.io`. This can be achieved by starting `bash` under an agent and
using the `principal` tool to obtain a blessing from `dev.v.io`.

```
jiri go install v.io/x/ref/cmd/principal
jiri go install v.io/x/ref/services/agent/...

$JIRI_ROOT/release/go/bin/agent --v23.credentials=<creds> bash
$JIRI_ROOT/release/go/bin/principal seekblessings
``` 

Above, `<creds>` points to a directory where the tool's Vanadium
credentials would be persisted.

## Scanning for lock devices

Find the names of nearby lock devices (on the same network as this tool)

```
lock scan
```

This command prints out the names of nearby lock devices. The name of an
unclaimed lock device always begins with `unclaimed-`, and the name of a claimed
lock device is the name under which it was claimed.

## Claiming an unclaimed device
The `claim` command can be used to claim an unclaimed device. For instance,
the following command claims the device `unclaimed-lock-xxxxxx` with the
name `front-door`. (Here `unclaimed-lock-xxxxxx` is the name of the unclaimed
lock device obtained from the `scan` command.)

```
lock claim unclaimed-lock-xxxxxx front-door
```

The lock would now authenticate with the name `front-door` and a subsequent invocation
of `scan` should print the name `front-door`.

## Listing available keys
The `listkeys` command lists the set of available physical-lock keys and the names
of the locks to which they apply.

```
lock listkeys
```

For example, as a result of the previous claim, the key `front-door:key` would show up
as being available to this tool for the lock `front-door`.

## Locking and Unlocking
The lock `front-door` can be locked using

```
lock lock front-door
```

and unlocked using 

```
lock unlock front-door
```

The `status` command can be used to determine the current status of the lock.

```
lock status front-door
```

## Sharing keys
The `lock` tool can also be used to share keys with nearby users
(on the same network as this tool). This involves the following steps.

First, the receiver must run the `recvkeys` command on their `lock` tool.

```
// At the receiver
lock recvkeys
```

As a result of this command, the receiver's `lock` tool would start a server for
receiving keys and advertise the endpoint of the server over MDNS.

Next, the sender must then search its neighborhood for users waiting to receiving keys

```
// At the sender
lock users
```

This command would print the email addresses of all users waiting to receive keys.
Once the email address of the receiver is visible then this command can be stopped
and `sendkey` command can be used to send a key.

For instance, the following command sends the key for lock `front-door` to the
 user `john.smith@gmail.com` that is only valid for 10 minutes.

```
// At the sender
lock sendkey --for=10m front-door john.smith@gmail.com friend
```

Above, `friend` is the category used to classify the receiver.

The `sendkey` command would cause the receiver to see a prompt specifying 
the name of the key received and the blessings of the sender. Once the 
receiver confirms saving the key, he/she would be able to use it to interact
with the lock `front-door` for next 10 minutes. Executing the `listkeys`
command would reveal the key `front-door:key:friend` along with its expiration time.

# Future Work

1) Auditing: Make lock service persist all requests received by it in an audit log.
This may be a struct of the form:

```
AuditLog(startTime, endTime time.Time) []AuditLog | error
type AuditLog struct {
  Blessing string
  Action LockStatus
  Timestamp time.Time
}
```

We'd also have to work out how to control access to this `AuditLog`. One option
is to use caveats - so when "making" a new key one can choose to insert the
"noadmin" caveat?

2) Caveats: Currently the `sendkey` command only supports expiration commands.
In the future, we would like to support at least the following caveats:

* Time-Range Caveat: This caveat would restrict the key so that it can only be used
  within a specific time range.
* Ask-Permission Caveat: This is a third-party caveat that requires the holder of the key
  to ask the granter for permission before using a key.

[MDNS]: http://en.wikipedia.org/wiki/Multicast_DNS
[agent]: https://github.com/vanadium/docs/blob/master/glossary.md#agent
