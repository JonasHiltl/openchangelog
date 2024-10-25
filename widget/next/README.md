# Openchangelog Next.js Widget

A Next.js Server Component to embed your Openchangelog changelog into your Next.js app.

## Installation
```
npm i @openchangelog/next
```

## Usage
```ts
import { Changelog } from "@openchangelog/next"


export default function ChangelogPage() {
  return (
    <Changelog 
      workspaceID="ws_xxxx"                         // when using Openchangelog cloud
      changelogID="cl_xxxx"                         // when using Openchangelog cloud
      baseUrl="https://your-changelog-instance.com" // when self-hosting
      theme="dark"
    />
  );
}
```

## Suspense
The `Changelog` component is built as an async component, making it compatible with React Suspense. You can wrap it in a Suspense boundary to show loading states while the changelog data is being fetched:

```ts
import { Changelog } from "@openchangelog/next"

export default function ChangelogPage() {
  return (
    <Suspense fallback={<div>Loading changelog...</div>}>
        <Changelog />
    </Suspense>
  );
}
```