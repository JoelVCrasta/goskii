# goskii

goskii is a command-line tool for converting both local or web images into ASCII art. Designed for simplicity and efficiency, it allows users to create artistic text-based art of their images directly from the terminal.
Future updates aim to include video support, making goskii a comprehensive tool for making ASCII art.

![Version](https://img.shields.io/badge/Version-1.0-blue.svg)

## Features

- Fast image-to-ASCII conversion
- Covert web images to ASCII
- Save and render ASCII art
- 13 ASCII character sets
- Customize image size

## Installation

Note: This installation script only works for x86-64 architecture. For ARM or other architectures, you can build it from source using Go on your system.

### Linux/Mac

Execute the below command and you are set.

```
curl -sSL https://raw.githubusercontent.com/JoelVCrasta/goskii/refs/heads/main/install.sh | bash
```

### Windows

Open Powershell as administrator and execute the below command.

```
Set-ExecutionPolicy Bypass -Scope Process -Force; Invoke-WebRequest -Uri "https://raw.githubusercontent.com/JoelVCrasta/goskii/refs/heads/main/install.ps1" -OutFile install.ps1; ./install.ps1
```

## Usage

| Options         | Type     | Description                                                        |
| :-------------- | :------- | :----------------------------------------------------------------- |
| `--path, -p`    | `string` | Path to the image file (Required)                                  |
| `--charset, -c` | `int`    | Character set to use (1 - 13). Default is 1                        |
| `--help, -h`    | `flag`   | Show help information for goskii                                   |
| `--output, -o`  | `string` | Output folder path. Default is current directory                   |
| `--render, -r`  | `string` | Render the contents of the ASCII art file                          |
| `--showset, -s` | `flag`   | Display all available character sets                               |
| `--version, -v` | `flag`   | Show the version of goskii                                         |
| `--width, -w`   | `int`    | Width of the ASCII art (1 - 500). Default adjusts to terminal size |

## Examples

Covert local image to ASCII art.

```
goskii -p ./example.png
```

Convert web-hosted image to ASCII art.

```
goskii -p "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQE67yVsNZ_zCtVCT7_bdIbzSib4BSuOwBFhg&s"
```

Specify custom ASCII art width.

```
goskii -p ./example.png -w 250
```

Show available character sets

```
goskii -s
```

Use a specific character set

```
goskii -p ./example.png -c 10
```
