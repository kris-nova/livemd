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

Use `livemd` as your main source of livestreaming state.

### Acquire Twitch ID

Go to the [VODs page](https://dashboard.twitch.tv/u/krisnova/content/video-producer) for your channel.

```
https://dashboard.twitch.tv/u/{{channel}}/content/video-producer
```

Find the most recent highlighter link, and pull the ID from the link.

```
https://dashboard.twitch.tv/u/{{channel}}/content/video-producer/highlighter/{{id}}
```

Now you can build the following URLs for your live stream. 

``` 
# Video URL for the specific video
https://www.twitch.tv/videos/{{id}}

# Edit page for the specific video
https://dashboard.twitch.tv/u/krisnova/content/video-producer/edit/{{id}}
```


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

The `description` string is used to populate the YouTube video description, as well as the archive webpage. 

``` 
live stream update --notify "Had a rough night. Come watch me live. Exclusively on @Twitch. https://twitch.tv/krisnova" --description "You should never run this code, ever!"
```

### Go Live Notifications 

```

```




