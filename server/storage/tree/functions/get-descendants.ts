import { join, compact } from 'lodash-es'
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
            onlyBranch,
            includeMeta,
        }: {
            onlyBranch?: number | boolean,
            includeMeta?: number | boolean,
        } = options || {}

        let stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                ${join(compact([
                    'node.*',
                    includeMeta ? 'meta.title, meta.description, meta.created_date, meta.modified_date, meta.data_source, meta.data_source_type' : '',
                    'nodeWithDepth.depth AS depth',
                ]), ',')}
            FROM
                nodes AS node,
                nodes AS parent
            LEFT JOIN
                nodes_meta AS meta
            ON (
                node.id = meta.id
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
                nodeWithDepth.id = node.id
            )
            WHERE (
                    node.left BETWEEN parent.left
                    AND parent.right
                )
                AND parent.id = $nodeId
                AND node.root = $rootId
                AND node.right <> (
                    CASE WHEN $onlyBranch IS TRUE
                        THEN node.left + 1
                        ELSE -1
                    END
                )
            ORDER BY node.left;
        `)
        stmt.bind(<BindParams>{
            nodeId:     <SqlValue>nodeId,
            rootId:     <SqlValue>rootId,
            onlyBranch: <SqlValue>onlyBranch,
        })

        const results = stmt.all()
        // let results = []
        // for  (const item of stmt.iterate()) {
        //     results.push(item)
        // }
        // stmt.free()

        return results
    }
}
