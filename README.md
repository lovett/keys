# Keys

Run arbitrary commands headlessly from a secondary keyboard or remotely.

## Why

Every key on a computer keyboard has the potential to do anything, but can't due to preexisting obligations. If you press the "a" key, not seeing the letter "a" appear on the screen is going to be a problem.

If you attach a second keyboard to your computer nothing amazing happens. The "a" key on keyboard 1 does the same thing as the "a" key on keyboard 2.

Keys solves both problems by letting you remap a key on a designated keyboard to any command you'd otherwise run from the command line.

It also provides a browser interface to allow these remapped keys to be pressed remotely.

## Limitations

Physical keyboard support is currently Linux-only.

Commands are run as if they were issued from the command line. Input isn't sent to other applications, as with a Stream Deck.

# Setup

Run `make setup` to install alsa-lib-devel, a prerequisite for playing sounds.

# Build

Run `make build` to compile the application.

# Usage

The application configuration file can be edited with a text editor or through the browser editor. The browser editor also has documentation and examples.

Running `keys setup` will generate a systemd user service to run the server in the background on port 4004. It can be reverse-proxied however you like.

If using with a physical keyboard, use `keys select keyboard` to pick which one to pay attention to. By default, input from all attached keyboards will be used.

Run `keys test sound` to verify that audio is working correctly.

Run `keys test key` to see the name of a pressed key. For letter and number keys this will probably be what you expect, but function and multimedia and numpad keys can be exotic and hardware-specific.

Run `keys start --help` for details on customizing the server.

# Attribution

This project uses icons from [Majesticons](https://github.com/halfmage/majesticons) and sound files from [Google Material Design v2](https://m2.material.io/design/sound/sound-resources.html).
