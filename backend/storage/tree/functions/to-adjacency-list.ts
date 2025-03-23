import {
    isEmpty,
    forEach,
    set,
} from 'lodash-es'

import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
    TreePlainNode,
} from '../tree.d'

export default function (context: ITree<ITreeArgs>): ITreeFuncContext {
    return (nodeId: string, options?: object): ITreeFuncResult => {
        const { rootId } = context

        const {
            onlyBranch,
            includeMeta,
        }: {
            onlyBranch?: number | boolean,
            includeMeta?: number | boolean,
        } = options || {}

        const evalNodeChild = (node: TreePlainNode) => {
            const childNodes = context.getChildren(node.id, { onlyBranch, includeMeta })

            if (!isEmpty(childNodes)) {
                set(node, 'children', childNodes)
            }

            forEach(childNodes, (node: TreePlainNode) => evalNodeChild(node))
            return node
        }

        const targetNode = context.getNode(nodeId || rootId)
        return evalNodeChild(targetNode as TreePlainNode)
    }
}
