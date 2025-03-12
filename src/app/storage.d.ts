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
    subtitles?: Array<{ stream_index: number, codec_name: string, lang_code: string, lang_title: string }>
    dubs?: Array<{ stream_index: number, codec_name: string, lang_code: string, lang_title: string }>
    duration?: float
    sourceReady?: boolean
}

export type SessionInfo = {
    id: string
    rootId: string
    nodeId: string
    files: { [fileKey: string]: FileMetaInfo }
}

export type ISocketMessage = {
    event: number
    payload?: string | { stage: number, message?: string | object }
}

export type IFileListItem = FileMetaInfo & {
    path: string
    fileKey: string
    sessionId: string
    nodeId: string
    stage?: number
    message?: string
    isDownloading: boolean
    isCompleted: boolean
    isProgressing: boolean
}
