export type EEnvName = 'development' | 'production'

export type IEnvConfigEndpoint = {
    backend: string
    api: string
    file: string
    frontend: string
}

export type IEnvConfig = {
    env: EEnvName
    endpoint: IEnvConfigEndpoint
}
