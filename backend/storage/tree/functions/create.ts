import { extend, has, unset } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (parentId: string, nodeId?: string): ITreeFuncResult => {
        const { db } = context

        let stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                node.id AS nodeId
            FROM
                nodes AS node
            WHERE node.id = $nodeId
            LIMIT 1;
        `)
        const existResult = stmt.get(<BindParams>{
            nodeId: <SqlValue>nodeId,
        })
        if (existResult && has(existResult, 'nodeId')) {
            return false
        }

        nodeId = nodeId || context.getNewID()
        let memo: {
            nodeId: string,
            parentId: string,
            rootId?: string,
            parentRight?: number,
            nodeLevel?: number,
        } = { nodeId, parentId }

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.root AS rootId
            FROM
                nodes AS node
            WHERE node.id = $parentId
            LIMIT 1;
        `)
        memo = extend({}, memo, stmt.get(<BindParams>{
            parentId: <SqlValue>memo.parentId,
        }))
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            INSERT OR IGNORE INTO nodes (
                id, root, parent, left, right, level
            )
            VALUES (
                $nodeId, $rootId, $parentId, 0, 0, 0
            );
        `)
        stmt.run(<BindParams>{
            nodeId: <SqlValue>memo.nodeId,
            rootId: <SqlValue>memo.rootId,
            parentId: <SqlValue>memo.parentId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.left AS nodeLeft,
                node.right AS nodeRight
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
                node.right AS parentRight
            FROM
                nodes AS node
            WHERE (
                node.id = $parentId
                AND node.root = $rootId
            )
            LIMIT 1;
        `)
        memo = extend({}, memo, stmt.get(<BindParams>{
            parentId: <SqlValue>memo.parentId,
            rootId: <SqlValue>memo.rootId,
        }))
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET left = left + 2
            WHERE (
                left >= $parentRight
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            parentRight: <SqlValue>memo.parentRight,
            rootId: <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET right = right + 2
            WHERE (
                right >= $parentRight
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            parentRight: <SqlValue>memo.parentRight,
            rootId: <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET
                left = $parentRight,
                right = $parentRight + 1
            WHERE (
                id = $nodeId
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            parentRight: <SqlValue>memo.parentRight,
            nodeId: <SqlValue>memo.nodeId,
            rootId: <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.id,
                (COUNT(parent.id) - 1) AS nodeLevel
            FROM
                nodes AS node,
                nodes AS parent
            WHERE (
                (node.left BETWEEN parent.left AND parent.right)
                AND node.id = $nodeId
                AND parent.root = node.root
            )
            GROUP BY node.id;
        `)
        memo = extend({}, memo, stmt.get(<BindParams>{
            nodeId: <SqlValue>memo.nodeId,
        }))
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET level = $nodeLevel
            WHERE (
                id = $nodeId
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            nodeLevel: <SqlValue>memo.nodeLevel,
            nodeId: <SqlValue>memo.nodeId,
            rootId: <SqlValue>memo.rootId,
        })
        // stmt.free()

        return { nodeId }
    }
}
