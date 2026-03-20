# freecam

Free your camera from macOS's `ptpcamera` daemon, which grabs exclusive USB/PTP access to any PTP-capable camera (Canon, Nikon, Sony, Fujifilm, Olympus, and more) and blocks third-party tools like Darktable, digiKam, and Lightroom.

## Install

```sh
brew install joennespreuwers/tap/freecam
```

Or build from source:

```sh
git clone https://github.com/joennespreuwers/freecam
cd freecam
go build -o freecam ./cmd/freecam
```

## Usage

```
freecam               # launch live TUI — watches and kills ptpcamera forever
freecam --once        # kill ptpcamera once and exit (no TUI)
freecam --process foo # target a different process instead of ptpcamera
freecam --version     # print version
freecam --help        # show help
```

### TUI keybindings

| Key | Action |
|-----|--------|
| `Q` | Quit |
| `P` | Pause / Resume watching |
| `C` | Clear event log |

## Why

macOS ships a `ptpcamera` daemon that auto-attaches to any PTP camera over USB (Canon, Nikon, Sony, Fujifilm, Olympus, and others). It locks the device so exclusively that no other software can open it — not even with root privileges. `freecam` solves this by detecting and killing the daemon whenever it respawns.

## Distribution

Releases are built with [GoReleaser](https://goreleaser.com) and published to GitHub Releases on every `v*` tag push. A Homebrew formula is auto-generated and pushed to the [homebrew-tap](https://github.com/joennespreuwers/homebrew-tap) repo.

## License

MIT
