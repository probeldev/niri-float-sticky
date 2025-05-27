# Niri Sticky Float  
*A utility to make floating windows visible across all workspaces in [niri](https://github.com/YaLTeR/niri) — similar to "sticky windows" in other compositors.*  

## Why?  
Niri doesn’t natively support global floating windows. This tool forces float windows to persist on every workspace, mimicking the `sticky` behavior from X11/Wayland compositors like Sway or KWin.  

## Installation

### Via Go:
```bash
go install github.com/probeldev/niri-float-sticky@latest
```


## Usage

To automatically launch the utility on niri startup, add this line to your niri configuration:

```kdl
spawn-at-startup "niri-floating-fixer"
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

