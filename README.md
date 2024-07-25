# Openchangelog
Openchangelog is an open source, self hostable Changelog Rendering Website.
Changelogs are written in Markdown and can be integrated from different sources, e.g `local` or `GitHub`.

## Configuration
You can configure your Changelog by adapting the `openchangelog.yml` file.
It is typically in `/etc/openchangelog.yml`.

### Look & Feel
**title**: The title is displayed above all Changelog articles.
```yaml
# openchangelog.yml
page:
  title:
```
**subtitle**: The subtitle is displayed below the title.
```yaml
# openchangelog.yml
page:
  subtitle:
```
**Logo**: Your logo is displayed in the header.
```yaml
# openchangelog.yml
page:
  logo:
    src: # url to image
    width: # width of logo as string e.g. 70px
    height: # height of logo as string e.g 30px
    link: # optional link that the logo points to
```

### Local Data Source
You can specify a local file path to a directory containing your Changelog Markdown files.
```yaml
# openchangelog.yml
local:
  filesPath: .testdata
```
### Github Data Source
You can specify your repository and path to a directory inside the repo containing your Changelog Markdown files.
You can **authenticate** via a `Github App` or `Personal Access Token`.  
```yaml
# openchangelog.yml
github:
  owner: # gh username
  repo:
  path: # path inside repo
  auth:
    accessToken: # access token with a access to the specified repo
    # or
    appPrivateKey:
    appInstallationId:
```

### Cache
You can configure a cache to improve latency and avoid hitting rate limits from e.g. Github.  
Internally [httpcache](https://github.com/gregjones/httpcache) is used to cache the request to Github.
You can choose between a `memory`, `disk` and `s3`.
```yaml
# openchangelog.yml
cache: 
  type: # disk, memory, s3
  disk: # used when type is disk
    location: # the file system location of the disk cache
    maxSize: # in bytes
  s3: # used when type is s3
    bucket: # the bucket url, env AWS_ACCESS_KEY_ID and AWS_SECRET_KEY are used as credentials
```

## Writing Changelogs
Each new Changelog, e.g. from a new release, is written in a new Markdown file, this allows adding custom metadata for each article using the markdown Frontmatter.  
All files are stored in the same directory, either local or in remote sources.

### Filename Format
The ordering of the changelog files is important.  
The displayed Articles in the UI are ordered by their filename in descending order. We recommend prefixing the file with the release version.
```
{version}.{title}.md
v0.0.1.darkmode.md
v0.0.2.i81n.md
```

### Content Format
Changelogs are written in Markdown, we are compliant with CommonMark 0.31.2 (thanks to [goldmark](https://github.com/yuin/goldmark)).  
You can use the Frontmatter to specify additional info:
- `title`: Displayed as a bold title at the top of the Changelog Article.
- `description`: Displayed below the title.
- `publishedAt`: An `ISO 8601` datetime that is diplayed next to the Changelog Article.