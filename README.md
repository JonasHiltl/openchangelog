# Openchangelog
Openchangelog is an open source, self hostable Changelog Rendering Website.
Changelogs are written in Markdown and can be integrated from different sources, e.g `local` or `GitHub`.

## Configuration
You can configure your Changelog by adapting the `config.yaml` file.

### Look & Feel
**title**: The title is displayed above all Changelog articles.
```yaml
# config.yaml
page:
  title:
```
**subtitle**: The subtitle is displayed blow the title.
```yaml
# config.yaml
page:
  subtitle:
```
**Logo**: Your logo is displayed in the header.
```yaml
# config.yaml
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
# config.yaml
local:
  filesPath: .testdata
```
### Github Data Source
You can specify your repository and path to a directory inside the repo containing your Changelog Markdown files.
You can **authenticate** via a `Personal Access Token`.
```yaml
# config.yaml
github:
  owner: # gh username
  repo:
  path: # path inside repo
  auth:
    accessToken: # access token with a access to the specified repo
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