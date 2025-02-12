import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (): ITreeFuncResult => {
        const { db, rootId } = context

        const stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                node.*,
                nodeWithDepth.depth AS depth
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
            WHERE node.id = $nodeId
            LIMIT 1;
        `)
        const results = stmt.get(<BindParams>{
            nodeId: <SqlValue>rootId,
        })
        // stmt.free()

        return results
    }
}
