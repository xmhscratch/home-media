import { get } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string): ITreeFuncResult => {
        const { db, rootId } = context

        let stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                indexes.nodeId,
                indexes.nodeIndex
            FROM
                (
                    SELECT
                        node.id AS nodeId,
                        ROW_NUMBER () OVER ( 
                            ORDER BY node.left
                        ) AS nodeIndex
                    FROM
                        nodes AS node
                    WHERE (
                        node.root = $rootId
                    )
                ) AS indexes
            WHERE (
                indexes.nodeId = $nodeId
            )
            LIMIT 1;
        `)
        let results = stmt.get(<BindParams>{
            nodeId: <SqlValue>nodeId,
            rootId: <SqlValue>rootId,
        })
        // stmt.free()

        return get(results, 'nodeIndex', -1)
    }
}
