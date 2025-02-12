import { extend, join, compact } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string, options?: object): ITreeFuncResult => {
        const { db } = context

        const {
            includeMeta,
        }: {
            includeMeta?: number | boolean,
        } = options || {}

        let stmt: IStatement
        let memo: {
            nodeId: string,
            rootId?: string,
        } = { nodeId }

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.root AS rootId
            FROM
                nodes AS node
            WHERE node.id = $nodeId
            LIMIT 1;
        `)
        memo = extend({}, memo, stmt.get(<BindParams>{
            nodeId: <SqlValue>memo.nodeId,
        }))
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            SELECT
                ${join(compact([
                    'parent.*',
                    includeMeta ? 'meta.title, meta.description, meta.created_date, meta.modified_date, meta.data_source, meta.data_source_type' : '',
                    'nodeWithDepth.depth AS depth',
                ]), ',')}
            FROM
                nodes AS node,
                nodes AS parent
            LEFT JOIN
                nodes_meta AS meta
            ON (
                parent.id = meta.id
            )
            LEFT JOIN (
                SELECT
                    subNode_d.id,
                    MAX(node_d.level - parent_d.level) AS depth
                FROM
                    nodes AS node_d,
                    nodes AS subNode_d,
                    nodes AS parent_d
                WHERE (
                        node_d.left BETWEEN parent_d.left AND parent_d.right
                    )
                    AND parent_d.id = subNode_d.id
                    AND node_d.root = parent_d.root
                GROUP BY subNode_d.id
                ORDER BY node_d.left
            ) AS nodeWithDepth
            ON (
                parent.id = nodeWithDepth.id
            )
            WHERE (
                (node.left BETWEEN parent.left AND parent.right)
                AND node.id = $nodeId
                AND parent.root = $rootId
            )
            ORDER BY parent.left;
        `)
        stmt.bind(<BindParams>{
            nodeId: <SqlValue>memo.nodeId,
            rootId: <SqlValue>memo.rootId,
        })

        const results = stmt.all()
        // let results = []
        // while (stmt.step()) {
        //     results.push(stmt.getAsObject())
        // }
        // stmt.free()

        return results
    }
}
