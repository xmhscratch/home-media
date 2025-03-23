import {
    extend,
    isEqual,
    toNumber,
    toString,
} from 'lodash-es'
import ObjectID from 'bson-objectid'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string, parentId: string, adjacentId?: string |  number): ITreeFuncResult => {
        const { db } = context

        let stmt: IStatement
        let memo: {
            nodeId: string,
            parentId: string,
            rootId?: string,
            parentRight?:  number,
            nodeLevel?:  number,
            adjacentId?: string,
            parentLevel?:  number,
            newPosition?:  number,
            parentLeft?:  number,
            nodeLeft?:  number,
            nodeRight?:  number,
            nodeWidth?:  number,
        } = { nodeId, parentId }

        if (isEqual(nodeId, parentId)) {
            return true
        }

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.root AS rootId,
                node.left AS nodeLeft,
                node.right AS nodeRight,
                node.level AS nodeLevel,
                node.right - node.left + 1 AS nodeWidth
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
                node.left AS parentLeft,
                node.right AS parentRight,
                node.level AS parentLevel
            FROM
                nodes AS node
            WHERE (
                node.id = $parentId
                AND node.root = $rootId
            )
            LIMIT 1;
        `)
        memo = extend({}, memo, stmt.get(<BindParams>{
            parentId:  <SqlValue>memo.parentId,
            rootId:    <SqlValue>memo.rootId
        }))
        // stmt.free()

        if (ObjectID.isValid(toString(adjacentId))) {
            stmt = (<IDatabase>db).prepare(`
                SELECT
                    parent.id AS adjacentId,
                    parent.left AS adjacentLeft,
                    parent.right AS adjacentRight,
                    parent.right AS newPosition
                FROM
                    nodes AS node,
                    nodes AS parent
                WHERE (
                    node.id = $adjacentId
                    AND node.left BETWEEN parent.left AND parent.right
                    AND node.root = $rootId
                    AND parent.level = $parentLevel + 1
                )
                LIMIT 1;
            `)
            memo = extend({}, memo, stmt.get(<BindParams>{
                parentLevel:   <SqlValue>memo.parentLevel,
                adjacentId:    <SqlValue>memo.adjacentId,
                rootId:        <SqlValue>memo.rootId,
            }))
            // stmt.free()
        }

        if (!memo.newPosition) {
            if (toNumber(adjacentId) < 0) {
                memo.newPosition = memo.parentRight
            }
            else {
                memo.newPosition = memo.parentLeft
            }
        }

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET
                left = 0 - left,
                right = 0 - right
            WHERE (
                left >= $nodeLeft
                AND right <= $nodeRight
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            nodeLeft:  <SqlValue>memo.nodeLeft,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET left = left - $nodeWidth
            WHERE (
                left > $nodeRight
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            nodeWidth: <SqlValue>memo.nodeWidth,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET right = right - $nodeWidth
            WHERE (
                right > $nodeRight
                AND root = $rootId
            );
        `)
        stmt.run(<BindParams>{
            nodeWidth: <SqlValue>memo.nodeWidth,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET left = left + $nodeWidth
            WHERE (
                CASE WHEN $newPosition > $nodeRight
                    THEN left > ($newPosition - $nodeWidth)
                    ELSE left > $newPosition
                END
            )
            AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeWidth:     <SqlValue>memo.nodeWidth,
            newPosition:   <SqlValue>memo.newPosition,
            nodeRight:     <SqlValue>memo.nodeRight,
            rootId:        <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET right = right + $nodeWidth
            WHERE (
                CASE WHEN $newPosition > $nodeRight
                    THEN right > ($newPosition - $nodeWidth)
                    ELSE right > $newPosition
                END
            )
            AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeWidth:     <SqlValue>memo.nodeWidth,
            newPosition:   <SqlValue>memo.newPosition,
            nodeRight:     <SqlValue>memo.nodeRight,
            rootId:        <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET
                left = CASE WHEN $newPosition > $nodeRight
                    THEN 0 - left + ($newPosition - $nodeRight)
                    ELSE 0 - left + ($newPosition - $nodeRight + $nodeWidth)
                END,
                right = CASE WHEN $newPosition > $nodeRight
                    THEN 0 - right + ($newPosition - $nodeRight)
                    ELSE 0 - right + ($newPosition - $nodeRight + $nodeWidth)
                END,
                level = level + ($parentLevel - ($nodeLevel - 1))
            WHERE left <= 0 - $nodeLeft
                AND right >= 0 - $nodeRight
                AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeWidth:     <SqlValue>memo.nodeWidth,
            newPosition:   <SqlValue>memo.newPosition,
            nodeLeft:      <SqlValue>memo.nodeLeft,
            nodeRight:     <SqlValue>memo.nodeRight,
            parentLevel:   <SqlValue>memo.parentLevel,
            nodeLevel:     <SqlValue>memo.nodeLevel,
            rootId:        <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE nodes
            SET parent = $parentId
            WHERE id = $nodeId
                AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            parentId:  <SqlValue>memo.parentId,
            nodeId:    <SqlValue>memo.nodeId,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        return true
    }
}
