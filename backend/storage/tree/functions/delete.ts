import {
    extend,
    pickBy,
    identity,
    map,
    first,
    pullAt,
    forEach,
} from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string): ITreeFuncResult => {
        const { db } = context

        let stmt: IStatement
        let memo: {
            nodeId: string,
            nodeLeft?:  number,
            nodeRight?:  number,
            rootId?: string,
            nodeWidth?:  number,
        } = { nodeId }

        stmt = (<IDatabase>db).prepare(`
            SELECT
                node.id AS nodeId,
                node.root AS rootId,
                node.left AS nodeLeft,
                node.right AS nodeRight,
                (node.right - node.left + 1) AS nodeWidth
            FROM
                nodes as node
            WHERE (
                node.id = $nodeId
                AND node.parent IS NOT NULL
            );
        `)
        memo = extend({}, memo,
            pickBy(stmt.get(<BindParams>{
                nodeId: <SqlValue>nodeId,
            }), identity)
        )
        // stmt.free()

        if (!memo.nodeId) {
            return `Node ${memo.nodeId} does not exist`
        }

        const targetNodes = context.getDescendants(memo.nodeId)
        memo = extend({}, memo, { targetNodes })

        const nodeIds = map(targetNodes as Array<object>, 'id')

        // let $nodeId: string
        // do {
        // } while (nodeIds.length)

        forEach(nodeIds, ($nodeId: string) => {
            $nodeId = first(pullAt(nodeIds, 0))
            stmt = (<IDatabase>db).prepare(`
                DELETE FROM nodes
                WHERE id = $nodeId;
            `)
            stmt.run(<BindParams>{ nodeId: <SqlValue>$nodeId })
            // stmt.free()
        })

        stmt = (<IDatabase>db).prepare(`
            DELETE FROM nodes
            WHERE (
                left BETWEEN $nodeLeft AND $nodeRight
            )
            AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeLeft:  <SqlValue>memo.nodeLeft,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE
                nodes
            SET
                left = left - $nodeWidth
            WHERE left > $nodeRight
                AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeWidth: <SqlValue>memo.nodeWidth,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        stmt = (<IDatabase>db).prepare(`
            UPDATE
                nodes
            SET
                right = right - $nodeWidth
            WHERE right > $nodeRight
                AND root = $rootId;
        `)
        stmt.run(<BindParams>{
            nodeWidth: <SqlValue>memo.nodeWidth,
            nodeRight: <SqlValue>memo.nodeRight,
            rootId:    <SqlValue>memo.rootId,
        })
        // stmt.free()

        return null
    }
}
