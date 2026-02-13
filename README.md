# Niri Float Sticky  
*A utility to make floating windows visible across all workspaces in [niri](https://github.com/YaLTeR/niri) — similar to "sticky windows" in other compositors.*  

## Why?  
Niri doesn’t natively support global floating windows. This tool forces float windows to persist on every workspace, mimicking the `sticky` behavior from X11/Wayland compositors like Sway or KWin.  

## Installation

### Via Go:
```bash
go install github.com/probeldev/niri-float-sticky@latest
```

### Via [AUR](https://aur.archlinux.org/packages/niri-float-sticky) (maintained by [jamlotrasoiaf](https://github.com/jamlotrasoiaf)/[brainworms2002](https://aur.archlinux.org/account/brainworms2002)):
```bash
paru -S niri-float-sticky
```

### Via Nix:
```bash
nix profile install github:probeldev/niri-float-sticky 
```

## Usage

To automatically launch the utility on niri startup, add this line to your niri configuration:

```kdl
spawn-at-startup "niri-float-sticky"
```

### Command Line Options

```bash
Usage of niri-float-sticky:
  -allow-moving-to-foreign-monitors
        allow moving to foreign monitors
  -app-id value
        only move floating windows with app-id matching given patterns
  -debug
        enable debug logging
  -disable-auto-stick
        disable auto sticking for all windows
  -ipc string
        send IPC command to daemon: set_sticky, unset_sticky, toggle_sticky
  -title value
        only move floating windows with title matching this pattern
  -version
        print version and exit
```

Notes:

- `-app-id` and `-title` can be provided multiple times to specify different patterns.
- Each flag accepts regex, so you can also provide multiple patterns in one flag (e.g. `foo|bar`).
- Internally all patterns are combined into a single regex with alternatives

Example with debug log:
```bash
niri-float-sticky -debug >> /tmp/niri-float-sticky.log

# Configuring logrotate
cat <<EOF | sudo tee /etc/logrotate.d/niri-float-sticky >/dev/null
/tmp/niri-float-sticky.log {
    daily
    rotate 5
    compress
    missingok
    notifempty
    copytruncate
    maxsize 10M
    su root root
}
EOF
```


### IPC

`niri-float-sticky` exposes a simple UNIX socket IPC interface to control window stickiness at runtime.

The daemon creates a socket at:

```
$XDG_RUNTIME_DIR/niri-float-sticky.sock
```

This allows external commands (for example, keybindings) to toggle stickiness of the currently focused window.

### How It Works

The binary has two modes:

1. **Daemon mode** (default)
   Runs the event loop and listens for IPC commands.

2. **Client mode** (`-ipc`)
   Sends a command to the running daemon and exits immediately.


### Available IPC Commands

```bash
niri-float-sticky -ipc set_sticky
niri-float-sticky -ipc unset_sticky
niri-float-sticky -ipc toggle_sticky
```

* `set_sticky` — force window to be sticky
* `unset_sticky` — remove manual override
* `toggle_sticky` — toggle sticky state



### Example: Keybinding in niri

You can bind stickiness toggling to a key:

```
binds {
    Mod+G { 
        spawn "niri-float-sticky" "-ipc" "toggle_sticky"; 
    }
}
```

### Implementation Details

The IPC protocol uses a JSON message over a UNIX domain socket:

```json
{
  "action": "toggle_sticky",
  "window_id": 123456
}
```


## Contributing

We welcome all contributions! To get started:

1. **Open an Issue** to:
   - Report bugs
   - Suggest new features
   - Ask questions

2. **Create a Pull Request** for:
   - Bug fixes
   - New functionality
   - Documentation improvements


## License

This project is licensed under the **MIT License**.
