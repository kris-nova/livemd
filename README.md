# Kris Nóva Live

<!-- YouTube Video ID in 3 places -->

[YouTube Video Archive Link](https://www.youtube.com/watch?v=85wlpC6LQfk)

<a href="https://www.youtube.com/watch?v=85wlpC6LQfk
" target="_blank"><img src="http://img.youtube.com/vi/85wlpC6LQfk/hqdefault.jpg" width="480" height="360" border="10" /></a>

<!-- YouTube Video ID in 3 places -->

# <a href="https://twitch.tv/krisnova"><img src ="https://i.imgur.com/1H8qkDT.png" width="60px"> Live on Twitch!</a> 

Live streams exclusively on Twitch. [Follow Kris Nova](https://www.twitch.tv/krisnova) on Twitch to join in on the fun!

#### twitch.tv/krisnova

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

