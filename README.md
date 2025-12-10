# Keys

Trigger shell commands locally from a secondary keyboard or remotely from the browser.

## Why

Every key on a computer keyboard has the potential to do anything.

A single key usually has a single job. If you press the "a" key, either that letter appears on the screen or some other predetermined application-specific action occurs. Connect a second keyboard and its "a" key will perform the exact same job.

This application lets you remap a key on a designated keyboard to any command you'd otherwise run from the command line. Now that second "a" key can do something entirely unrelated.

The application also provides a browser interface so that remapped keys can be pressed remotely.

# Usage Scenario

You have an always-on Linux server and a spare keyboard for troubleshooting. Your TV is connected to an second adjacent Linux machine that suspends when idle. It supports Wake-on-LAN, but you have no way to trigger it while standing in front of the TV and yelling doesn't help. Like a savage, you resort to using the power button which is inconveniently located on the back of the machine which itself is inconveniently tucked away on its knee-level shelf.  With keys, you instead configure the "t" key to run the wake command from your server.

Sometimes you've already plopped yourself on the couch before remembering to turn on the TV computer. Walking 10 paces to the keyboard with the magic "t" key is obviously impossible. With keys, you instead take out your phone and navigate to the browser interface and get the job done from there.

# Limitations

Physical keyboard support is currently Linux-only.

Commands are run as if they were issued from the command line. Input isn't sent to other applications. This is not a Stream Deck.

# Setup

Run `scripts/setup.sh` to install alsa-lib-devel, a prerequisite for sound effects.

# Build

Run `scripts/build.sh` to compile the application.

# Usage

The application configuration file can be edited with a text editor or through the browser interface. The browser interface has documentation and examples.

Running `keys setup` will generate a systemd user service to run the server in the background on port 4004. It can be reverse-proxied however you like.

If using a physical keyboard, use `keys select keyboard` to pick which one to pay attention to. By default, input from all attached keyboards will be used.

Run `keys test sound` to verify that audio is working correctly.

Run `keys test key` to see the name of a pressed key. For letter and number keys this will probably be what you expect, but function and multimedia and numpad keys can be more exotic.

Run `keys start --help` for details on customizing the server. Changing the server port and turning off web or keyboard input are supported.

# API

There is an OpenAPI spec at `localhost:4004/openapi.yaml` in case you're into that sort of thing.

# Attribution

This project uses icons from [Majesticons](https://github.com/halfmage/majesticons) and sound files from [Google Material Design v2](https://m2.material.io/design/sound/sound-resources.html).
