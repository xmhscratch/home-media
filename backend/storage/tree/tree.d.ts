export interface ITreeFuncContext extends Function { }
export type ITreeFuncResult = unknown
export type BindParams = unknown
export type SqlValue = unknown

export const TreeFuncArgs: Array<T>
export type TreePlainNode = {
    id: string,
    _id?: string,
    root: string,
    parent?: string,
    left: number,
    right: number,
    level: number,
    children?: Array<TreePlainNode>,
}

export type TreeMetaNode = {
    id: string,
    root: string,
    title: string,
    description?: string,
    createdDate?: string,
    modifiedDate?: string,
    dataSource: string,
    dataSourceType: int,
}

export type ITreeArgs = {
    db: IDatabase | undefined
    rootId: string
}

export interface ITree<ITreeArgs> {
    [x: unknown]: unknown

    db: IDatabase | undefined
    rootId: string

    _memoizer: object

    destroy(): Promise<unknown>

    getNewID(): string

    create(...TreeFuncArgs): ITreeFuncResult
    import(...TreeFuncArgs): ITreeFuncResult
    moveTo(...TreeFuncArgs): ITreeFuncResult
    delete(...TreeFuncArgs): ITreeFuncResult

    getNode(...TreeFuncArgs): ITreeFuncResult
    getRootNode(...TreeFuncArgs): ITreeFuncResult
    getPrevNode(...TreeFuncArgs): ITreeFuncResult
    getNodeByParentIndex(...TreeFuncArgs): ITreeFuncResult
    getPaths(...TreeFuncArgs): ITreeFuncResult
    getLevel(...TreeFuncArgs): ITreeFuncResult
    getDepth(...TreeFuncArgs): ITreeFuncResult
    getIndexOf(...TreeFuncArgs): ITreeFuncResult
    getChildren(...TreeFuncArgs): ITreeFuncResult
    getDescendants(...TreeFuncArgs): ITreeFuncResult
    countChilds(...TreeFuncArgs): ITreeFuncResult
    toAdjacencyList(...TreeFuncArgs): ITreeFuncResult
    toLinearList(...TreeFuncArgs): ITreeFuncResult

    memoize(fnName: string, origFn: Function, fnArgs: unknown): ITreeFuncResult
}

export interface IMemoizer extends object {
    [x: string]: Map
}

export interface IDatabase {
    prepare(stmt: string): IStatement
    transaction(fn: Function): Function
    pragma(string: string, options?: object): object | number | string
    backup(stmt: string, options?: object): Promise<unknown>
    serialize(options?: object): Buffer
    function(name: string, fn: Function): IDatabase
    function(name: string, options: object, fn: Function): IDatabase
    aggregate(name: string, options: object): IDatabase
    table(name: string, definition: object): IDatabase
    loadExtension(path: string, entryPoint?: unknown): IDatabase
    exec(stmt: string): IStatement
    close(): IStatement
}

export interface IStatement {
    run(...bindParams: Array<unknown>): object
    get(...bindParams: Array<unknown>): object
    all(...bindParams: Array<unknown>): Array<unknown> | object
    iterate(...bindParams: Array<unknown>): unknown
    pluck(toggleState?: boolean): IStatement
    expand(toggleState?: boolean): IStatement
    raw(toggleState?: boolean): IStatement
    columns(): Array<unknown> | object
    bind(...bindParams: Array<unknown>): IStatement
}
