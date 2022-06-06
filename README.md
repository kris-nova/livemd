# LiveMD

All your stream meta information, in a single collaborative markdown file called a **state file**.

LiveMD is a live stream utility that enables a streamer to store meta information about their stream in a state file.
Create a new state file for every stream, use the state file to track titles, notes, description, key words, tags, URLs, for every live stream.

Current supported integrations:

 - Twitch
 - YouTube
 - Mastodon
 - HackMD
 - Twitter
 - Discord

# Workflow

Use `livemd` as your main source of live streaming state.

### Create a new stream.

This will create a new local state file called `live.md` in the current directory.

This file will be the centralized state for the duration of your entire live stream.

```bash 
live stream new <title>
live stream new "Hacking live on Kubernetes! Let's Go!"
```

### Add fields to your stream.

The `notify` string is a 280 character string that is used to notify various services such as Twitter, Discord, etc

``` 
live stream update --notify "Had a rough night. Come watch me live. Exclusively on @Twitch. https://twitch.tv/krisnova"
```

The `description` string



