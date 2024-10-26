import { Changelog } from "@openchangelog/next"
import { Suspense } from "react";

export default function ChangelogPage() {
  return (
    <div>
      <h1>My Next App</h1>
      <Suspense fallback={<p>Loading...</p>}>
        <Changelog baseUrl="http://localhost:6001" changelogID="cl_crcrelibeadmd11vnkl0" workspaceID="ws_crco2lqbeadg93uuhjk0" theme="light" />
      </Suspense>
    </div>
  );
}
