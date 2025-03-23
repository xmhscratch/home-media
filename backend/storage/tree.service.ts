import { Injectable, Scope, OnModuleInit } from '@nestjs/common';

import { StorageService } from './storage.service'
import { Tree } from './tree/tree';
import { IDatabase } from './tree/tree.d';

// import {
//     map,
// } from 'lodash-es'

@Injectable({ scope: Scope.REQUEST })
export class TreeService {

    constructor(private storage: StorageService) { }

    // async onModuleInit(): Promise<void> {
    //     return Promise.resolve(void 0)
    // }

    async customEx(rootId: string, fn: Function) {
        const db: IDatabase = await this.storage.getDb(rootId)

        const tree = new Tree(db, rootId)
        await this.initTree(db)
        await tree.createRoot(rootId)

        const fnResult = await fn(tree, db)

        await this.storage.closeDb(rootId, db)
        return fnResult
    }

    async execute(rootId: string, fnName: string, fnArgs: any[]) {
        return this.customEx(rootId, async (tree: Tree, db: IDatabase) => {
            return tree[fnName](...fnArgs)
        })
    }

    async initTree(db: IDatabase) {
        db.exec(`CREATE TABLE IF NOT EXISTS nodes (
            id varchar(24) PRIMARY KEY NOT NULL,
            root varchar(24) NOT NULL,
            parent varchar(24),
            left int(11) UNIQUE NOT NULL,
            right int(11) UNIQUE NOT NULL,
            level int(11) NOT NULL
        );`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_id ON nodes (id);`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_root ON nodes (root);`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_parent ON nodes (parent);`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_left ON nodes (left);`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_right ON nodes (right);`)

        db.exec(`CREATE TABLE IF NOT EXISTS nodes_meta (
            id varchar(24) PRIMARY KEY NOT NULL,
            root varchar(24) NOT NULL,
            title varchar(512),
            description text,
            created_date datetime NOT NULL,
            modified_date datetime,
            data_source text,
            data_source_type int(1) NOT NULL
        );`)
        db.exec(`CREATE INDEX IF NOT EXISTS node_meta_root ON nodes_meta (root);`)
    }
}
