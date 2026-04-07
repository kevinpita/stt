# stt (Send To Telegram)

A simple Go CLI tool to send files to a Telegram chat using a bot.

## Installation

### Using `go install`

If you have Go installed on your system, you can install `stt` directly from the source:

```bash
go install github.com/kevinpita/stt@latest
```

Ensure your `GOPATH/bin` is in your `PATH`.

### From Releases

You can download the pre-built binary for Linux from the [GitHub Releases](https://github.com/kevinpita/stt/releases) page.

1. Download the latest binary.
2. Make it executable: `chmod +x stt`
3. Move it to your path: `sudo mv stt /usr/local/bin/`

## Configuration

The easiest way to configure `stt` is to use the built-in setup command:

```bash
stt --setup
```

This will prompt you for your Telegram bot token and chat ID and save them to `~/.config/stt.conf`.

Alternatively, you can manually create the file:

```toml
token = "YOUR_BOT_TOKEN"
chat_id = "YOUR_CHAT_ID"
```

## Usage

Send one or more files to your Telegram chat:

```bash
stt file1.jpg file2.pdf video.mp4
```

`stt` automatically detects the file type (photo, video, or document) and uses the appropriate Telegram API endpoint.
