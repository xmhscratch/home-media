import { isEmpty } from 'lodash-es'
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
                node.id,
                MAX(node.level - parent.level) AS depth
            FROM
                nodes AS node,
                nodes AS parent
            WHERE (
                    node.left BETWEEN parent.left AND parent.right
                )
                AND parent.id = $nodeId
                AND node.root = parent.root
            ORDER BY node.left;
        `)
        let results = stmt.get(<BindParams>{
            nodeId: <SqlValue>nodeId,
            rootId: <SqlValue>rootId,
        })
        // stmt.free()

        return !isEmpty(results) ? <SqlValue>results['depth'] : 0
    }
}
