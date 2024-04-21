# Openchangelog

## Filename Format
The ordering of the changelog files is important.  
The displayed Articles in the UI are ordered by their filename in descending order. We recommend prefixing the file with the release version.
```
{version}.{title}.md
v0.0.1.darkmode.md
v0.0.2.i81n.md
```

## Content Format
Changelogs are written in Markdown, we are compliant with CommonMark 0.31.2 (thanks to [goldmark](https://github.com/yuin/goldmark)).  
You can use the Frontmatter to specify additional info:
- `title`: Displayed as a bold title at the top of the Changelog Article.
- `description`: Displayed below the title.
- `publishedAt`: An `ISO 8601` datetime that is diplayed next to the Changelog Article.