import { get, join, compact } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string, options?: object): ITreeFuncResult => {
        const { db, rootId } = context

        const {
            includeMeta,
        }: {
            includeMeta?: number | boolean,
        } = options || {}

        const nodeInfo = context.getNode(nodeId)
        const prevNodeParent = get(nodeInfo, 'parent', null)
        const prevNodeRight = get(nodeInfo, 'left', 2) - 1

        const stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                ${join(compact([
                    'node.*',
                    includeMeta ? 'meta.title, meta.description, meta.created_date, meta.modified_date, meta.data_source, meta.data_source_type' : '',
                    'nodeWithDepth.depth AS depth',
                ]), ',')}
            FROM
                nodes as node,
                (
                    SELECT
                        MAX(node_d.level - parent_d.level) AS depth
                    FROM
                        nodes AS node_d,
                        nodes AS parent_d
                    WHERE (
                            node_d.left BETWEEN parent_d.left AND parent_d.right
                        )
                        AND parent_d.id = $nodeId
                        AND node_d.root = parent_d.root
                    ORDER BY node_d.left
                ) AS nodeWithDepth
            LEFT JOIN
                nodes_meta AS meta
            ON (
                node.id = meta.id
            )
            WHERE node.root = $rootId
                AND node.right = $prevNodeRight
                AND node.parent = $prevNodeParent
            LIMIT 1;
        `)
        const results = stmt.get(<BindParams>{
            nodeId:         <SqlValue>nodeId,
            rootId:         <SqlValue>rootId,
            prevNodeParent: <SqlValue>prevNodeParent,
            prevNodeRight:  <SqlValue>prevNodeRight,
        })
        // stmt.free()

        return results
    }
}
