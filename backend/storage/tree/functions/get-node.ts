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
            includeMeta,
        }: {
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
            WHERE (
                node.id = $nodeId
                AND node.root = $rootId
            )
            LIMIT 1;
        `)
        stmt.bind(<BindParams>{
            nodeId: <SqlValue>nodeId,
            rootId: <SqlValue>rootId,
        })
        let results = stmt.get()
        // stmt.free()

        return results
    }
}
