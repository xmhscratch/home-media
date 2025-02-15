export type IINode = {
    id: string
    root: string
    created_date?: unknown | null
    data_source?: unknown | null
    data_source_type?: unknown | null
    depth?: number
    description?: string | null
    left?: number
    level?: number
    modified_date?: unknown | null
    parent?: string | null
    right?: number
    title?: unknown | null
}

export type ITreeRootNode = {
    rootId: string
    label: string
}

export type FileMetaInfo = {
    path: string
    size: number
}

export type SessionInfo = {
    id: string
    rootId: string
    nodeId: string
    files: FileMetaInfo[]
}

export type IFileListItem = FileMetaInfo & {
    path: string
    sessionId: string
}
