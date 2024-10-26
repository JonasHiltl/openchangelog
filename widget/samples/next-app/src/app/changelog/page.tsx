import { Changelog } from "@openchangelog/next"
import { Suspense } from "react";

export default function ChangelogPage() {
  return (
    <div>
      <h1>My Next App</h1>
      <Suspense fallback={<p>Loading...</p>}>
        <Changelog changelogID="cl_cqvk7ich990s5lmf7a2g" workspaceID="ws_cqvk7g4h990s5lmf7a1g" theme="light" />
      </Suspense>
    </div>
  );
}
