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
  -debug
        enable debug logging
  -version
        print version and exit
```

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
