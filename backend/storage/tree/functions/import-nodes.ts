import { map, omit, pick } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import insertNodeMeta from './insert-node-meta'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
    TreePlainNode,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (plainNodes?: Array<TreePlainNode>): ITreeFuncResult => {
        const { db } = context;
        const { rootId } = context;

        (<IDatabase>db)
            .prepare('DELETE FROM nodes WHERE nodes.root = $rootId;')
            .bind(<BindParams>{
                rootId: <SqlValue>rootId,
            })
            .run()

        return Promise.all(
            map(plainNodes, (item: TreePlainNode) => new Promise((resolve) => {
                let stmt: IStatement = (<IDatabase>db).prepare(`INSERT OR IGNORE INTO nodes (
                    id, root, parent, left, right, level
                ) VALUES (
                    $nodeId, $rootId, $parentId, $left, $right, $level
                )`)

                const params: BindParams & {
                    nodeId: SqlValue,
                    rootId: SqlValue,
                    parentId: SqlValue,
                    left: SqlValue,
                    right: SqlValue,
                    level: SqlValue,
                } = {
                    nodeId: <SqlValue>(item.id || item._id),
                    rootId: <SqlValue>item.root,
                    parentId: <SqlValue>item.parent,
                    left: <SqlValue>item.left,
                    right: <SqlValue>item.right,
                    level: <SqlValue>(item.level || 0),
                }
                stmt.bind(params)
                stmt.run()

                item.id = <string>params.nodeId
                item.root = <string>params.rootId

                const meta = pick(
                    omit(
                        { ...params, ...item },
                        ['_id', 'parent', 'parentId', 'left', 'right', 'level', 'nodeId', 'rootId'],
                    ),
                    ['id', 'root', 'title', 'description', 'createdDate', 'modifiedDate', 'dataSource', 'dataSourceType'],
                )
                return resolve(meta)
            }))
        )
            .then((results) => Promise.all(
                map(results, (metaItem: object) => insertNodeMeta(context)(metaItem))
            ))
    }
}
