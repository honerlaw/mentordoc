import {SyncAction} from "../sync-action";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {IDocumentState} from "../../model/document/document-state";
import {AclDocument} from "../../model/document/acl-document";
import {IRootState} from "../../model/root-state";

export const SET_SEARCH_DOCUMENTS_TYPE: string = "set_search_documents_type";

export  interface ISetSearchDocuments {
    searchDocuments: AclDocument[] | null;
}

export interface ISetSearchDocumentsDispatch {
    setSearchDocuments: (req?: ISetSearchDocuments) => void;
}

export interface ISetSearchDocumentsSelector {
    searchDocuments: AclDocument[] | null;
}

class SetSearchDocumentsImpl extends SyncAction<IDocumentState, ISetSearchDocuments, AclDocument[] | null> {

    public constructor() {
        super(ReducerType.DOCUMENT, SET_SEARCH_DOCUMENTS_TYPE, "searchDocuments", "setSearchDocuments")
    }

    public handle(state: IDocumentState, action: IWrappedAction<ISetSearchDocuments>): IDocumentState {
        state = cloneDeep(state);
        if (action.payload) {
            state.searchDocuments = action.payload.searchDocuments;
        }
        return state;
    }

    public getSelectorValue(state: IRootState): AclDocument[] | null {
        return state.document.searchDocuments;
    }

}

export const SetSearchDocuments: SetSearchDocumentsImpl = new SetSearchDocumentsImpl();
