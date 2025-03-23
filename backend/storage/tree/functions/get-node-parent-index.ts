import { get, join, compact } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (parentId: string, parentIndex: number, options?: object): ITreeFuncResult => {
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
                ]), ',')},
                (
                    SELECT
                        COUNT(id)
                    FROM
                        nodes AS nthNode
                    WHERE (
                        nthNode.right < node.left
                        AND nthNode.parent = node.parent
                        AND nthNode.root = node.root
                    )
                ) AS parentIndex
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
                        AND parent_d.id = $parentId
                        AND node_d.root = parent_d.root
                    ORDER BY node_d.left
                ) AS nodeWithDepth
            LEFT JOIN
                nodes_meta AS meta
            ON (
                node.id = meta.id
            )
            WHERE (
                node.parent = $parentId
                AND node.root = $rootId
                AND parentIndex = $parentIndex
            )
            LIMIT 1;
        `)
        let results = stmt.get(<BindParams>{
            rootId:        <SqlValue>rootId,
            parentId:      <SqlValue>parentId,
            parentIndex:   <SqlValue>parentIndex,
        })
        // stmt.free()

        return get(results, 'id') ? results : {}
    }
}
