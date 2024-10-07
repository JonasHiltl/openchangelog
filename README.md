<p align="center">
  <img width="750" src="https://github.com/user-attachments/assets/64259c72-17b7-479e-8079-98d7e2b72593"/>
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
  <a href="https://twitter.com/jonasdevs">Twitter</a>
</p>
<br />
<br />
</p>

Openchangelog takes your Markdown files, hosted on GitHub or locally and renders them as a beautiful Changelog Website.
- Dark, Light and System themes
- Automatic RSS feed
- Easy to self-host, just a single config file
- Written in Go → lightweight
- Various integrations, open an issue to request a new integration

## Quickstart
Create an `openchangelog.yml` config file, for more infos see the [configuration](#configuration) section.
```
docker run -v ./openchangelog.yml:/etc/openchangelog.yml:ro -v ./release-notes:/release-notes -p 6001:6001 ghcr.io/jonashiltl/openchangelog:0.1.9
```
Or
```yaml
services:
  openchangelog:
    image: "ghcr.io/jonashiltl/openchangelog:0.1.9"
    ports:
      - "6001:6001"
    volumes:
      - ./release-notes:/release-notes
      - type: bind
        source: openchangelog.yml
        target: /etc/openchangelog.yml
```
Once deployed, your changelog will be available at http://localhost:6001.