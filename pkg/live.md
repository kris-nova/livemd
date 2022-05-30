# {{ .Title }}
### Episode {{ .I }}
{{ if .YouTubeID }}
[YouTube Video Archive Link](https://www.youtube.com/watch?v={{ .YouTubeID }})
<a href="https://www.youtube.com/watch?v={{ .YouTubeID }}
" target="_blank"><img src="http://img.youtube.com/vi/{{ .YouTubeID }}/hqdefault.jpg" width="480" height="360" border="10" /></a>
# <a href="https://twitch.tv/krisnova"><img src ="https://i.imgur.com/1H8qkDT.png" width="60px"> Live on Twitch!</a> 
{{ end }}
Live streams exclusively on Twitch. [Follow Kris Nova](https://www.twitch.tv/krisnova) on Twitch to join in on the fun!

#### twitch.tv/krisnova

### Configuration

data:
```json
{
  "title": "{{ .Title }}"
}
```
data:


### Notes

 - ...
 - ...

