# {{ .Title }}

 - [Twitch Edit Page](https://dashboard.twitch.tv/u/{{ .TwitchChannel }}/content/video-producer/edit/{{ .TwitchID }})
 - [Twitch Video Page](https://www.twitch.tv/videos/{{ .TwitchID }})

{{ .Description }}{{ if .Notify }}

### Notification 🔔

```
{{ .Notify }}
```
{{ end }}{{ if .Twitter }}{{ .Twitter }}{{ end }}
{{if .Links }}### References{{ range .Links }}
 - {{ .Markdown }}{{ end }}{{ end }}
---

⚠ Do Not Edit Below This Line ⚠ Do Not Edit Below This Line ⚠

---