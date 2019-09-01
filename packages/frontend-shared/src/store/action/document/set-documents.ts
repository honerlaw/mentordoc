import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {AclDocument} from "../../model/document/acl-document";
import {IDocumentState} from "../../model/document/document-state";

export const SET_DOCUMENTS_TYPE: string = "set_documents_type";

export type SelectorValue = (orgId: string, folderId: string | null) => AclDocument[];

export interface ISetDocuments {
    documents: AclDocument[];
}

export interface ISetDocumentsSelector extends ISelectorMap {
    getDocuments: SelectorValue;
}

export interface ISetDocumentsDispatch extends IDispatchMap {
    setDocuments: (req?: ISetDocuments) => void;
}

class SetDocumentsImpl extends SyncAction<IDocumentState, ISetDocuments, SelectorValue> {

    public constructor() {
        super(ReducerType.DOCUMENT, SET_DOCUMENTS_TYPE, "getDocuments", "setDocuments")
    }

    public handle(state: IDocumentState, action: IWrappedAction<ISetDocuments>): IDocumentState {
        state = cloneDeep(state);
        if (action.payload) {
            action.payload.documents.forEach((doc: AclDocument): void => {
                const key: string = this.getKeyForDocument(doc);

                const documents: AclDocument[] | undefined = state.documentMap[key];

                // no folders exist for the key, so just add it
                if (!documents) {
                    state.documentMap[key] = [doc];
                    return;
                }

                // otherwise lets make sure the folder isn't in the map already, and then add it
                const found: AclDocument | undefined = documents
                    .find((existing: AclDocument): boolean => existing.model.id === doc.model.id);
                if (!found) {
                    documents.push(doc);
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorValue {
        return (orgId: string, folderId: string | null): AclDocument[] => {
            const documents: AclDocument[] | undefined = state.document.documentMap[this.getKey(orgId, folderId)];
            if (documents) {
                return documents;
            }
            return [];
        };
    }

    private getKeyForDocument(document: AclDocument): string {
        return this.getKey(document.model.organizationId, document.model.folderId);
    }

    private getKey(orgId: string, folderId: string | null): string {
        if (folderId) {
            return `${orgId}-${folderId}`;
        }
        return orgId;
    }

}

export const SetDocuments: SetDocumentsImpl = new SetDocumentsImpl();
