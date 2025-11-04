# Keys

Use a regular keyboard as a macro pad to run arbitrary commands headlessly.

## Why

Every key on a computer keyboard has the potential to do anything, but necessity usually makes that impossible. Most keys already serve some purpose, so reassignment isn't practical.

Modifiers like Alt and Control can help, but not all combinations are ergonomic. Macro pads are another option: mini supplemental keyboards with a handful of keys that can be programmed to do what can't be achieved from the main keyboard alone.

But if you have a spare keyboard lying around, why not use that.

Or a tablet or phone that can run a virtual keyboard through a web browser. That can work too.

## Limitations

Physical keyboard support is currently Linux-only.

Commands are run as if they were issued from the command line. There is no application-specific integration as with a Stream Deck, or relay/redirection of key events.

# Setup

Run `make setup` to install alsa-lib-devel, a prerequisite for playing sounds.

# Attribution

This project uses icons from [Majesticons](https://github.com/halfmage/majesticons) and sound files from [Google Material Design v2](https://m2.material.io/design/sound/sound-resources.html).
