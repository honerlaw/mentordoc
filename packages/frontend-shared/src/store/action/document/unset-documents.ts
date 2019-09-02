import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {AclDocument} from "../../model/document/acl-document";
import {IDocumentState} from "../../model/document/document-state";

export const UNSET_DOCUMENTS_TYPE: string = "unset_documents_type";

export interface IUnsetDocuments {
    documents: AclDocument[];
}

export interface IUnsetDocumentsDispatch extends IDispatchMap {
    unsetDocuments: (req?: IUnsetDocuments) => void;
}

export class UnsetDocumentsImpl extends SyncAction<IDocumentState, IUnsetDocuments, void> {

    public constructor() {
        super(ReducerType.DOCUMENT, UNSET_DOCUMENTS_TYPE, "unsetDocuments", "unsetDocuments")
    }

    public handle(state: IDocumentState, action: IWrappedAction<IUnsetDocuments>): IDocumentState {
        state = cloneDeep(state);
        if (action.payload) {
            action.payload.documents.forEach((doc: AclDocument): void => {
                const key: string = this.getKeyForDocument(doc);

                const documents: AclDocument[] | undefined = state.documentMap[key];

                // no documents exist for the key, so nothing to remove
                if (!documents) {
                    return;
                }

                // remove the document
                for (let i: number = documents.length - 1; i >= 0; --i) {
                    if (documents[i].model.id === doc.model.id) {
                        documents.splice(i, 1);
                        break;
                    }
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): void {
        // does nothing!
    }

    private getKeyForDocument(doc: AclDocument): string {
        return this.getKey(doc.model.organizationId, doc.model.folderId);
    }

    private getKey(orgId: string, folderId: string | null): string {
        if (folderId) {
            return `${orgId}-${folderId}`;
        }
        return orgId;
    }

}

export const UnsetDocuments: UnsetDocumentsImpl = new UnsetDocumentsImpl();
