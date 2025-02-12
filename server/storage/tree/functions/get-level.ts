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
        const { db } = context

        let stmt: IStatement = (<IDatabase>db).prepare(`
            SELECT
                node.id,
                (COUNT(parent.id) - 1) AS level
            FROM
                nodes AS node,
                nodes AS parent
            WHERE (
                node.id = $nodeId
                AND parent.root = node.root
                AND node.left BETWEEN parent.left AND parent.right
            )
            GROUP BY node.id;
        `)
        let results = stmt.get(<BindParams>{
            nodeId: <SqlValue>nodeId,
        })
        // stmt.free()

        return !isEmpty(results) ? <SqlValue>results['level'] : 0
    }
}
