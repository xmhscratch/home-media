import ObjectID from 'bson-objectid'
import * as hashObject from 'object-hash'
import {
    isEmpty,
    has,
    memoize,
    join,
    forEach,
} from 'lodash-es'

import getRootNode from './functions/get-root-node'
import getPaths from './functions/get-paths'
import createNode from './functions/create'
import getNode from './functions/get-node'
import getLevel from './functions/get-level'
import getDepth from './functions/get-depth'
import getIndexOf from './functions/get-index-of'
import getChildren from './functions/get-children'
import getDescendants from './functions/get-descendants'
import countChilds from './functions/count-childs'
import moveTo from './functions/move-to'
import deleteNode from './functions/delete'
import toAdjacencyList from './functions/to-adjacency-list'
import importNodes from './functions/import-nodes'
import insertNodeMeta from './functions/insert-node-meta'
import toLinearList from './functions/to-linear-list'
import getPrevNode from './functions/get-prev-node'
import getNodeByParentIndex from './functions/get-node-parent-index'

import {
    ITree,
    ITreeArgs,
    TreePlainNode,
    TreeMetaNode,
    TreeFuncArgs,
    IMemoizer,
} from './tree.d'
import { IDatabase, IStatement, BindParams, SqlValue } from './tree.d'

export const getNewID = (): string => ObjectID().toHexString()

export class Tree implements ITree<ITreeArgs> {

    static MODIFIER_FUNCTIONS = ['create', 'import', 'insert', 'moveTo', 'delete']

    db: IDatabase | undefined
    rootId: string

    _memoizer: IMemoizer = {}

    constructor(db: IDatabase, rootId?: string) {
        this.db = db
        this.rootId = rootId || getNewID()

        return this
    }

    async destroy() {
        return this.db?.close()
    }

    getNewID(): string { return getNewID() }
    static getNewID(): string { return getNewID() }

    async createRoot(rootId: string) {
        const stmt: IStatement = this.db.prepare(`INSERT OR IGNORE INTO nodes (
            id, root, parent, left, right, level
        ) VALUES (
            $rootId, $rootId, NULL, 1, 2, 0
        );`)
        stmt.bind(<BindParams>{ rootId: <SqlValue>rootId })
        stmt.run()

        return this
    }

    create = (...args: any[]) => this.memoize('create', createNode(this), args)
    import = (...args: any[]) => this.memoize('import', importNodes(this), args)
    insert = (...args: any[]) => this.memoize('insert', insertNodeMeta(this), args)
    moveTo = (...args: any[]) => this.memoize('moveTo', moveTo(this), args)
    delete = (...args: any[]) => this.memoize('delete', deleteNode(this), args)

    getNode = (...args: any[]) => this.memoize('getNode', getNode(this), args)
    getRootNode = (...args: any[]) => this.memoize('getRootNode', getRootNode(this), args)
    getPrevNode = (...args: any[]) => this.memoize('getPrevNode', getPrevNode(this), args)
    getNodeByParentIndex = (...args: any[]) => this.memoize('getNodeByParentIndex', getNodeByParentIndex(this), args)
    getPaths = (...args: any[]) => this.memoize('getPaths', getPaths(this), args)
    getLevel = (...args: any[]) => this.memoize('getLevel', getLevel(this), args)
    getDepth = (...args: any[]) => this.memoize('getDepth', getDepth(this), args)
    getIndexOf = (...args: any[]) => this.memoize('getIndexOf', getIndexOf(this), args)
    getChildren = (...args: any[]) => this.memoize('getChildren', getChildren(this), args)
    getDescendants = (...args: any[]) => this.memoize('getDescendants', getDescendants(this), args)
    countChilds = (...args: any[]) => this.memoize('countChilds', countChilds(this), args)
    toAdjacencyList = (...args: any[]) => this.memoize('toAdjacencyList', toAdjacencyList(this), args)
    toLinearList = (...args: any[]) => this.memoize('toLinearList', toLinearList(this), args)

    memoize(fnName: string, invokeFn: Function, fnArgs: any[]) {
        let fnResult: any

        // fnResult = invokeFn(...fnArgs)
        const fnKeyGetter = (fnArgs: any[]): string => !isEmpty(fnArgs) ? `${fnName}_${hashObject(fnArgs)}` : fnName
        if (!has(this._memoizer, `${fnKeyGetter(fnArgs)}`)) {
            this._memoizer[`${fnKeyGetter(fnArgs)}`] = memoize<any>(invokeFn, fnKeyGetter) 
        }

        if ((new RegExp(`^(${join(Tree.MODIFIER_FUNCTIONS, '|')})$`, 'g')).test(fnName)) {
            forEach(this._memoizer, (v, k) => {
                if (has(this._memoizer, k)) {
                    this._memoizer[k].cache.clear()
                }
            })
            fnResult = invokeFn.apply(this, fnArgs)
        }
        else {
            const memoFn = this._memoizer[`${fnKeyGetter(fnArgs)}`]
            fnResult = memoFn.apply(this, fnArgs)
        }

        return fnResult
    }
}

export default Tree
export type { ITree, TreePlainNode, TreeMetaNode }
