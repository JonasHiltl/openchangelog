import { Changelog } from "@openchangelog/next"
import { Suspense } from "react";

export default function ChangelogPage() {
  return (
    <div className="mt-10">
      <Suspense fallback={<p>Loading...</p>}>
        <Changelog
          changelogID="cl_cqvk7ich990s5lmf7a2g"
          workspaceID="ws_cqvk7g4h990s5lmf7a1g" theme="dark" />
      </Suspense>
    </div>
  );
}
