# Touch for Windows

**Touch for Windows** is an advanced version of the standard `touch` command, designed to work on Windows. This tool provides enhanced functionality like customizable timestamps, directory creation, and flexible file permission management.

## Features

- **Custom timestamps**: Set a specific creation timestamp (supports date and/or time).
- **Directory creation**: Create directories with ease.
- **Overwrite protection**: Prevent overwriting existing files.
- **Force overwrite**: Replace existing files without prompting.
- **File creation**: Create files.
- **Permissions management**: Set file permissions for specific users (e.g., USER, ADMIN).
- **Cross-platform support**: Works seamlessly on Windows.

## Prerequisites

- **Go**: You need Go 1.16 or later installed on your machine. You can download it from [here](https://golang.org/dl/).
- **Windows**: Designed to work on Windows, but can be used on other platforms if adjusted accordingly.

## Installation

To install **Touch for Windows**, follow these steps:

### 1. Clone the repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/PauWol/Windows-Touch.git
```
### 2. Build the project

Navigate to the project directory and build the project using Go:

```bash
cd touch
go build
```

### 3. (Optional) Move the executable to a directory in your PATH

You can move the touch.exe file to a directory in your system’s PATH for easier access. For example:

```	bash
mv touch.exe C:\Windows\System32\
```

This will allow you to use touch from any command prompt window without navigating to the project directory.

## Usage
Once the touch.exe binary is ready, you can use it via the command line.

### Command Format
```bash
touch [flags] [file(s)]
```
### Flags
- `-d`, `--directory`: Create directories instead of files.
- `-f`, `--force`: Overwrite existing files.
- `-t`, `--timestamp`: Set the creation timestamp for the file (format: YYYY-MM-DD HH:MM:SS). You can also provide just a date (YYYY-MM-DD) or time (HH:MM) to modify only the respective part.
- `-p`, `--permissions`: Set the file permissions. Accepts values like USER or ADMIN.

### Examples
Create a file with a custom timestamp:

```bash
touch -t "2025-02-09 15:00:00" myfile.txt
```

Create a directory:

```bash
touch -d mydirectory
```

Force overwrite an existing file:

```bash
touch -f myfile.txt
```

Set specific permissions (user or admin):

```bash
touch -p "USER" myfile.txt
```

### Timestamp Formatting
You can use the following formats for the `--timestamp` flag:

**Full Date and Time**: ``2025-02-09 15:00:00``
**Only Date**: ``2025-02-09``
**Only Time**: ``15:00``

When providing just a time, the date will remain unchanged, and vice versa.

## Troubleshooting

**Invalid Timestamp Format**: Ensure you provide the timestamp in the correct format. If only time or only date is provided, the other part will stay unchanged.
**Permission Errors**: If you encounter permission issues, try running the command as an Administrator or adjust file permissions as needed.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Note**: This tool is inspired by the classic `touch` command but specifically built to work well on Windows with extended functionality.
