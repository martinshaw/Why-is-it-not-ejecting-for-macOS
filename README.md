# Why-is-it-not-ejecting-for-macOS

## Overview

**Why-is-it-not-ejecting-for-macOS** is a command-line tool for macOS that helps you diagnose why an external drive (such as a USB stick or external hard disk) cannot be ejected. It provides detailed information about which processes or files are preventing the ejection, making it easier to resolve the issue.

## How It Works

The tool leverages several native macOS utilities and system calls to trace the cause of the ejection problem:

1. **diskutil**: Used to list and manage disks and volumes. The tool first identifies the target disk and gathers its details using `diskutil`.
2. **lsof**: Lists open files on the system. The tool uses `lsof` to find files on the target disk that are currently open, which can block ejection.
3. **ps**: Shows running processes. The tool correlates open files with the processes that are using them, providing a clear view of which applications or background processes are responsible.

By combining these steps, the tool gives you actionable information to close or terminate the right processes, allowing you to safely eject your drive.

## Why This Is Useful

- **No more guessing**: Instead of manually hunting for which app is blocking your drive, get a clear report.
- **Saves time**: Quickly resolve ejection issues without rebooting or force-unmounting.
- **Safer for your data**: Avoid data loss or corruption by ensuring all files are properly closed before ejection.

## How to Compile and Run

This project is written in Go. To compile it, you need to have Go installed (version 1.18 or later recommended).

1. Clone the repository:
   ```sh
   git clone https://github.com/martinshaw/Why-is-it-not-ejecting-for-macOS.git
   cd Why-is-it-not-ejecting-for-macOS
   ```
2. Build the executable:
   ```sh
   go build
   sudo chmod +X ./ejecting
   ```
3. Run it
   ```sh
   ./ejecting -ui menubar # for native macOS menubar GUI
   ./ejecting -format json # for JSON output
   ```
4. If you choose the menubar option, you can interact with the GUI to kill processes directly from there.

## Example Output

For:

```sh
./ejecting -format indent
```

Output:

```json
[
  {
    "Disk": {
      "DeviceIdentifier": "disk4s1",
      "MountPoint": "/Volumes/Samsung USB",
      "Size": 513283060224,
      "VolumeName": "Samsung USB"
    },
    "OpenFiles": [
      {
        "CommandName": "IINA",
        "CommandPath": "/Applications/IINA.app/Contents/MacOS/IINA",
        "PID": 28246,
        "User": "martin",
        "Name": "/Volumes/Samsung USB/Audiobooks/Tim S. Grover - Relentless.mp3"
      }
    ]
  }
]
```

## Requirements

- macOS (Tested on Monterey and later)
- Go 1.18+ (Don't have go? Just download the precompiled binary from the releases page)
- Access to `diskutil`, `lsof`, and `ps` (All standard on macOS - Don't worry about this!)
