type ChangelogProps = {
    // specify if using openchangelog cloud
    changelogID?: string
    // specify if using openchangelog cloud
    workspaceID?: string
    page?: number
    pageSize?: number
    // your own hosted openchangelog instance, defaults to https://openchangelog.com
    baseUrl?: string
    theme?: "dark" | "light"
}

async function fetchChangelog(args: ChangelogProps): Promise<string> {
    const baseURL = args.baseUrl || "https://openchangelog.com"
    const params = new URLSearchParams({
        widget: "true",
        ...(args.changelogID && { cid: args.changelogID }),
        ...(args.workspaceID && { wid: args.workspaceID }),
        ...(args.page && { page: args.page.toString() }),
        ...(args.pageSize && { "page-size": args.pageSize.toString() }),
    })
    const url = new URL(`?${params.toString()}`, baseURL)
    const res = await fetch(url, {
        cache: "default",
    })
    if (!res.ok) {
        throw new Error(`failed to render changelog ${res.status}: ${await res.text()}`)
    }
    return res.text()
}

/**
 * Render your Openchangelog changelog on the server.
 * Specify `changelogID` and `workspaceID` if using openchangelog cloud,
 * otherwise use `baseUrl` to point to your hosted instance.
 */
export async function Changelog(props: ChangelogProps) {
    const html = await fetchChangelog(props)

    return (
        <>
            <div color-scheme={props.theme}>
                <div dangerouslySetInnerHTML={{ __html: html }}></div>
            </div>
        </>
    )
}