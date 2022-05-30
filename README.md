# Live

Centralized live stream notes, metadata, and keyword tagging for Twitch, YouTube, and HackMD.

```
export HACKMD_TOKEN="49BY6KGJX4XW9MRKCTOZY8FVHPEZKIC6G0ILOI0D6N3XGQJ9Q7"
export HACKMD_ID="xxYJeDxMRASxW2jIAxjV_g"
```




# Live Tool

The `live` command line tool is written in Go.

Use the `live` tool to manage streams with Twitch, Hackmd, and Exporting to YouTube archive.

The `live` command line tool by default will assume `live.md` as the file we are editing.

```
# Edit ./live.md
live stream <title>    # Create a new live stream (hackmd)
live stream push       # Sync local changes to hackmd
live stream pull       # Overwrite local changes to hackmd

# Update firebot with new hackmd URL (TODO Automate)
live twitch push       # Sync local file to twitch API
live twitch pull       # Overwrite local changes from twitch API
live twitch export     # Export twitch episode to YouTube
```

