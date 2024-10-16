<p align="center">
  <a href="https://demo.openchangelog.com" target="_blank">
    <img width="750" src="https://github.com/user-attachments/assets/41d11528-0445-4370-8196-33a3c3923f58"/>
  </a>
  <h1 align="center"><b>Openchangelog</b></h1>
</p>
<p align="center">
  Transform your changelog Markdown files to beautiful product updates
  <br />
  <br />
  <a href="https://openchangelog.com">Website</a>
  ·
  <a href="https://openchangelog.com/docs/">Docs</a>
  ·
  <a href="https://cloud.openchangelog.com">Get Started</a>
  ·
  <a href="https://demo.openchangelog.com">Demo</a>
  ·
  <a href="https://twitter.com/jonasdevs">Twitter</a>
</p>
<br />
<br />
</p>

Openchangelog takes your Changelog, hosted on GitHub or locally and renders it as a beautiful Changelog Website.
- Dark, Light and System themes
- Colorful Tags
- Automatic RSS feed
- Password Protection
- Supports [keep a changelog](https://keepachangelog.com/en/1.1.0/) format or 1 Markdown file per release
- Various integrations, open an issue to request a new integration

## Quickstart
Create an `openchangelog.yml` config file, from the sample `openchangelog.example.yml`. For more configuration settings visit our [Docs](https://openchangelog.com/docs/getting-started/self-hosting/#configuration).
```
docker run -v ./openchangelog.yml:/etc/openchangelog.yml:ro -v ./release-notes:/release-notes -p 6001:6001 ghcr.io/jonashiltl/openchangelog:0.3.0
```
Or
```yaml
services:
  openchangelog:
    image: "ghcr.io/jonashiltl/openchangelog:0.3.0"
    ports:
      - "6001:6001"
    volumes:
      - ./release-notes:/release-notes
      - type: bind
        source: openchangelog.yml
        target: /etc/openchangelog.yml
```
Once deployed, your changelog will be available at http://localhost:6001.
