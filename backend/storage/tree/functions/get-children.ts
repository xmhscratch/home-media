import {
    extend,
    omit,
    map,
    join,
    compact,
} from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
    TreePlainNode,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string, opts: { includeMeta: boolean, onlyBranch: boolean }): ITreeFuncResult => {
        const { db } = context

        let memo: {
            nodeId: string,
            rootId?: string,
        } = { nodeId }

        let stmt: IStatement = (<IDatabase>db).prepare(`
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
            'node.*',
            opts.includeMeta ? 'meta.title, meta.description, meta.created_date, meta.modified_date, meta.data_source, meta.data_source_type' : '',
            'nodeWithDepth.depth AS depth',
        ]), ',')}
            FROM
                nodes AS node
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
                node.parent = $parentId
                AND node.root = $rootId
                AND node.right <> (
                    CASE WHEN $onlyBranch IS TRUE
                        THEN node.left + 1
                        ELSE -1
                    END
                )
            )
            ORDER BY node.left;
        `)
        stmt.bind(<BindParams>{
            rootId: <SqlValue>memo.rootId,
            parentId: <SqlValue>memo.nodeId,
            onlyBranch: <SqlValue>opts.onlyBranch,
        })

        const results = stmt.all()
        // while (stmt.step()) {
        //     let row = stmt.get()
        //     results.push(omit(row, 'nodeLevel'))
        // }
        // stmt.free()

        return map(results, (o: TreePlainNode) => omit(o, 'nodeLevel'))
    }
}

// stmt = (<IDatabase>db).prepare(`
//     SELECT
//         node.*,
//         (
//           COUNT(parent.id) - (subTree.nLevel + 1)
//         ) AS nodeLevel
//     FROM
//         nodes AS node,
//         nodes AS parent,
//         nodes AS subParent
//     INNER JOIN (
//         SELECT
//             node.id,
//             (COUNT(parent.id) - 1) AS nLevel
//         FROM
//             nodes AS node,
//             nodes AS parent
//         WHERE (node.left BETWEEN parent.left AND parent.right)
//             AND node.id = $nodeId
//             AND node.root = $rootId
//             AND parent.root = $rootId
//         GROUP BY node.id
//         ORDER BY node.left
//     ) AS subTree
//     WHERE (node.left BETWEEN parent.left AND parent.right)
//         AND (node.left BETWEEN subParent.left AND subParent.right)
//         AND subParent.id = subTree.id
//         AND subParent.root = $rootId
//         AND node.root = $rootId
//         AND parent.root = $rootId
//     GROUP BY node.id
//     HAVING nodeLevel = 1
//     ORDER BY node.left;
// `)
