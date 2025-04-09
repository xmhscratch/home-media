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
                COUNT(*) AS count
            FROM
                nodes AS child,
                (
                    SELECT
                        node.left AS left,
                        node.right AS right
                    FROM
                        nodes AS node
                    WHERE (
                        node.id = $nodeId
                        AND node.root = $rootId
                    )
                    LIMIT 1
                ) AS node
            WHERE (
                child.left BETWEEN node.left AND node.right
                AND child.right BETWEEN node.left AND node.right
                AND child.root = $rootId
            )
            LIMIT 1;
        `)
        let results: { count?: number } = stmt.get(<BindParams>{
            nodeId: <SqlValue>nodeId,
            rootId: <SqlValue>rootId,
        })
        // stmt.free()

        return !isEmpty(results) ? <SqlValue>results['count'] : 0
    }
}
