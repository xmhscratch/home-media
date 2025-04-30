import { map, omit, isEmpty, pick } from 'lodash-es'
import { IDatabase, IStatement, BindParams, SqlValue } from '../tree.d'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
    TreeMetaNode,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (metaItem?: TreeMetaNode): ITreeFuncResult => {
        const { db } = context;

        if (isEmpty(metaItem)) { return }

        if (isEmpty(
            omit(pick(metaItem, [
                'id', 'root',
                'nodeId',
                'rootId',
                'title',
                'description',
                'createdDate',
                'modifiedDate',
                'dataSource',
                'dataSourceType',
            ]), ['id', 'root']))
        ) { return }

        let stmt: IStatement = (<IDatabase>db).prepare(`INSERT OR IGNORE INTO nodes_meta (
            id, root, title, description, created_date, modified_date, data_source, data_source_type
        ) VALUES (
            $nodeId, $rootId, $title, $description, $createdDate, $modifiedDate, $dataSource, $dataSourceType
        )`)
        const params: BindParams & {
            nodeId: SqlValue,
            rootId: SqlValue,
            title: SqlValue,
            description: SqlValue,
            createdDate: SqlValue,
            modifiedDate: SqlValue,
            dataSource: SqlValue,
            dataSourceType: SqlValue,
        } = {
            nodeId: <SqlValue>metaItem!.id,
            rootId: <SqlValue>metaItem!.root,
            title: <SqlValue>metaItem!.title,
            description: <SqlValue>metaItem!.description,
            createdDate: <SqlValue>metaItem!.createdDate || (new Date()).toISOString(),
            modifiedDate: <SqlValue>metaItem!.modifiedDate,
            dataSource: <SqlValue>metaItem!.dataSource,
            dataSourceType: <SqlValue>metaItem!.dataSourceType || 0,
        }
        stmt.bind(params)
        stmt.run()

        return metaItem?.id
    }
}
