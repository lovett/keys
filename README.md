Trigger shell commands locally from a secondary keyboard or remotely from the browser.

## Why

If you plug more than one keyboard into your computer, they do the same thing. That's not useful.

A keyboard attached to a headless computer is also not very useful due to the lack of feedback.

In both scenarios the keyboard has limited utility because each key can only perform one action, even if that action makes no sense at the time.

This application lets you remap a key on a designated keyboard to a command you'd otherwise run from the command line. It also provides a browser interface so that these remapped keys can be pressed remotely.

## Limitations

Keyboard support is Linux-only.

Commands are run as if they were issued from the command line. Input isn't sent to other applications.

## Usage

Running `keys setup` will generate a systemd user service to run the server in the background on port 4004.

Run `keys start` to start the server directly. See `keys start --help` for further options.

If using a physical keyboard, use `keys select keyboard` to pick which one to pay attention to. By default, input from all attached keyboards will be used.

Run `keys test sound` to verify that audio is working correctly.

Run `keys test key` to see the name of a pressed key. For letter and number keys this will probably be what you expect, but others can be exotic.

## API

There is an OpenAPI spec at `localhost:4004/openapi.yaml`

## Shell Client

A POSIX shell script can be downloaded from `localhost:4004/util/keys.sh` to interact with the server remotely via curl.

## Development

Run `scripts/setup.sh` to install system packages.

Run `scripts/build.sh` to compile the application.

## Attribution

This project uses icons from [Majesticons](https://github.com/halfmage/majesticons) and sound files from [Google Material Design v2](https://m2.material.io/design/sound/sound-resources.html).
