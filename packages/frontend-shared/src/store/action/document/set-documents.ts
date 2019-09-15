import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {AclDocument} from "../../model/document/acl-document";
import {IDocumentState} from "../../model/document/document-state";

export const SET_DOCUMENTS_TYPE: string = "set_documents_type";

export type SelectorSetDocumentsChildValue = (orgId: string, folderId: string | null) => AclDocument[];
export type SelectorSetDocumentsMapValue = Record<string, AclDocument[]>;
export type SelectorSetDocumentsValue = (type: "map" | "child") => SelectorSetDocumentsChildValue | SelectorSetDocumentsMapValue;

export interface ISetDocuments {
    documents: AclDocument[];
}

export interface ISetDocumentsSelector extends ISelectorMap {
    getDocuments: SelectorSetDocumentsValue;
}

export interface ISetDocumentsDispatch extends IDispatchMap {
    setDocuments: (req?: ISetDocuments) => void;
}

class SetDocumentsImpl extends SyncAction<IDocumentState, ISetDocuments, SelectorSetDocumentsValue> {

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

                let found: boolean = false;

                // see if the folder exists already
                for (let i: number = documents.length - 1; i >= 0; --i) {
                    const d: AclDocument = documents[i];

                    // if it does exist, replace it and mark that we found it
                    if (d.model.id === doc.model.id) {
                        found = true;
                        documents[i] = doc;
                    }
                }

                // if it wasn't found then add it
                if (!found) {
                    documents.push(doc);
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorSetDocumentsValue {
        return (type: "child" | "map"): SelectorSetDocumentsChildValue | SelectorSetDocumentsMapValue => {
            if (type === "map") {
                return state.document.documentMap;
            }
            return (orgId: string, folderId: string | null): AclDocument[] => {
                const documents: AclDocument[] | undefined = state.document.documentMap[this.getKey(orgId, folderId)];
                if (documents) {
                    return documents;
                }
                return [];
            };
        }
    }

    private getKeyForDocument(document: AclDocument): string {
        return this.getKey(document.model.organizationId, document.model.folderId);
    }

    private getKey(orgId: string, folderId: string | null): string {
        if (folderId) {
            return `${orgId}:${folderId}`;
        }
        return orgId;
    }

}

export const SetDocuments: SetDocumentsImpl = new SetDocumentsImpl();
