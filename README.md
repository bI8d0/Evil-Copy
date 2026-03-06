# Evil-Copy 😈💾

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue?logo=windows&logoColor=white)
![Build](https://img.shields.io/badge/Build-Passing-brightgreen)

```
███████╗██╗   ██╗██╗██╗      ██████╗ ██████╗ ██████╗ ██╗   ██╗
██╔════╝██║   ██║██║██║     ██╔════╝██╔═══██╗██╔══██╗╚██╗ ██╔╝
█████╗  ██║   ██║██║██║     ██║     ██║   ██║██████╔╝ ╚████╔╝ 
██╔══╝  ╚██╗ ██╔╝██║██║     ██║     ██║   ██║██╔═══╝   ╚██╔╝  
███████╗ ╚████╔╝ ██║███████╗╚██████╗╚██████╔╝██║        ██║   
╚══════╝  ╚═══╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝        ╚═╝   
```

> **Stealth Backup — Your files are mine now!**

Evil-Copy is a cross-platform USB stealth backup tool written in Go. It silently monitors for newly connected USB drives and automatically copies all their contents to a local directory. Built for **Windows** and **Linux**.

---

## ⚠️ Disclaimer

**This tool is intended for educational and authorized security testing purposes only.** Unauthorized access to or copying of data from devices you do not own or have explicit permission to access is illegal and unethical. The author assumes no responsibility for any misuse of this software. Use at your own risk and always comply with applicable laws and regulations.

---

## ✨ Features

- 🔌 **Automatic USB detection** — Monitors for newly plugged-in removable drives in real time.
- 🪟 **Windows support** — Uses the Windows API (`GetDriveType`) to accurately identify removable drives.
- 🐧 **Linux support** — Uses `fsnotify` to watch `/media` for newly mounted USB devices.
- 📂 **Full device backup** — Recursively copies all files and directories from the USB drive.
- ⚡ **Concurrent file copying** — Uses goroutines with a semaphore to copy up to 5 files in parallel.
- 📝 **Backup logging** — Generates a `backup_log.txt` with the result of each file copy operation.
- 🕵️ **Stealth operation** — Runs silently in the background with no user interaction required.

---

## 📋 Requirements

- [Go](https://go.dev/dl/) 1.23 or later
- Git (to clone the repository)

### Dependencies

| Package | Description |
|---|---|
| [`github.com/fsnotify/fsnotify`](https://github.com/fsnotify/fsnotify) | Cross-platform filesystem notifications (used on Linux) |
| [`golang.org/x/sys`](https://pkg.go.dev/golang.org/x/sys) | Windows system calls (used for `GetDriveType` on Windows) |

---

## 🚀 Getting Started

### Clone the repository

```bash
git clone https://github.com/bI8d0/Evil-Copy.git
cd Evil-Copy
```

### Install dependencies

```bash
go mod download
```

---

## 🔨 Building

Evil-Copy includes a `build.go` script that cross-compiles binaries for both Windows and Linux.

### Build for all platforms

```bash
go run build.go
```

This will generate the following binaries inside the `build/` directory:

| File | Platform |
|---|---|
| `build/evil-copy.exe` | Windows (amd64) |
| `build/evil-copy` | Linux (amd64) |

### Build manually for a specific platform

**Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o build/evil-copy.exe main.go
```

**Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o build/evil-copy main.go
```

---

## ▶️ Usage

Simply run the compiled binary. Evil-Copy will start monitoring for USB drives automatically.

### On Windows

```powershell
.\evil-copy.exe
```

### On Linux

```bash
./evil-copy
```

Once a new USB drive is detected, all files will be copied to a `leaks/` directory located next to the executable:

```
leaks/
└── leaks_20260305_143022/
    ├── backup_log.txt
    ├── Documents/
    │   └── report.pdf
    ├── Photos/
    │   └── vacation.jpg
    └── ...
```

Each backup is stored in a timestamped subdirectory (`leaks_YYYYMMDD_HHMMSS`).

---

## 📁 Project Structure

```
Evil-Copy/
├── .gitignore            # Git ignore rules
├── main.go               # Main application logic (Linux + shared code)
├── main_windows.go       # Windows-specific implementation (removable drive detection)
├── build.go              # Cross-compilation build script
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── LICENSE               # MIT License
├── README.md             # Project documentation
└── build/                # Compiled binaries (generated, git-ignored)
    ├── evil-copy         # Linux binary
    └── evil-copy.exe     # Windows binary
```

---

## 🔧 How It Works

### Windows
1. On startup, the program scans all drive letters (`A:\` to `Z:\`) and records existing removable drives.
2. It then polls every second for new drive letters.
3. When a new removable drive is detected (via `windows.GetDriveType`), it triggers a full backup of the drive contents.

### Linux
1. The program uses `fsnotify` to watch the `/media` directory and its subdirectories.
2. When a new directory is created (indicating a USB mount), it verifies the path structure.
3. After a short delay (to allow the OS to fully mount the device), it triggers a full backup.

### Backup Process
1. Creates a timestamped directory under `leaks/`.
2. Walks the entire source directory tree.
3. Copies files concurrently using goroutines (max 5 simultaneous copies).
4. Logs each operation (success/failure) to `backup_log.txt`.

---

## 📄 License

This project is provided as-is for educational purposes. See [LICENSE](LICENSE) for details.

---

## 👤 Author

**bI8d0**

---

> *"Your files are mine now!"* 😈
