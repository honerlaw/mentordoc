import {SyncAction} from "../sync-action";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {IDocumentState} from "../../model/document/document-state";
import {AclDocument} from "../../model/document/acl-document";
import {IRootState} from "../../model/root-state";

export const SET_FULL_DOCUMENT_TYPE: string = "set_full_document_type";

export  interface ISetFullDocument {
    fullDocument: AclDocument | null;
}

export interface ISetFullDocumentDispatch {
    setFullDocument: (req?: ISetFullDocument) => void;
}

export interface ISetFullDocumentSelector {
    fullDocument: AclDocument | null;
}

class SetFullDocumentImpl extends SyncAction<IDocumentState, ISetFullDocument, AclDocument | null> {

    public constructor() {
        super(ReducerType.DOCUMENT, SET_FULL_DOCUMENT_TYPE, "fullDocument", "setFullDocument")
    }

    public handle(state: IDocumentState, action: IWrappedAction<ISetFullDocument>): IDocumentState {
        state = cloneDeep(state);
        if (action.payload) {
            state.fullDocument = action.payload.fullDocument;
        }
        return state;
    }

    public getSelectorValue(state: IRootState): AclDocument | null {
        return state.document.fullDocument;
    }

}

export const SetFullDocument: SetFullDocumentImpl = new SetFullDocumentImpl();
