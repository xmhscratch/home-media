import {
    ITree,
    ITreeArgs,
    ITreeFuncContext,
    ITreeFuncResult,
} from '../tree.d'

export default (context: ITree<ITreeArgs>): ITreeFuncContext => {
    return (nodeId: string, options?: object): ITreeFuncResult => {
        const { rootId } = context

        const {
            onlyBranch,
            includeMeta,
        }: {
            onlyBranch?: number | boolean,
            includeMeta?: number | boolean,
        } = options || {}

        return context.getDescendants(nodeId || rootId, { onlyBranch, includeMeta })
    }
}
