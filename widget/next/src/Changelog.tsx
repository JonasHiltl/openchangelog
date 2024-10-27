type BaseChangelogProps = {
    className?: string
    page?: number
    pageSize?: number
    baseUrl?: string
    theme?: "dark" | "light"
}

type LocalChangelogProps = BaseChangelogProps & {
    changelogID?: never
    workspaceID?: never
}

type CloudChangelogProps = BaseChangelogProps & {
    changelogID: string
    workspaceID: string
}

type ChangelogProps = LocalChangelogProps | CloudChangelogProps

async function fetchChangelog(args: ChangelogProps): Promise<string> {
    const baseURL = args.baseUrl || "https://app.openchangelog.com/"
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
        <div color-scheme={props.theme} className={props.className} dangerouslySetInnerHTML={{ __html: html }}></div>
    )
}