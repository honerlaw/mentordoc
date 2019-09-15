import {SyncAction} from "../sync-action";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {DocumentPath, IDocumentState} from "../../model/document/document-state";
import {IRootState} from "../../model/root-state";

export const SET_DOCUMENT_PATH_TYPE: string = "set_document_path_type";

export  interface ISetDocumentPath {
    documentPath: DocumentPath;
}

export interface ISetDocumentPathDispatch {
    setDocumentPath: (req?: ISetDocumentPath) => void;
}

export interface ISetDocumentPathSelector {
    documentPath: DocumentPath;
}

class SetDocumentPathImpl extends SyncAction<IDocumentState, ISetDocumentPath, DocumentPath | null> {

    public constructor() {
        super(ReducerType.DOCUMENT, SET_DOCUMENT_PATH_TYPE, "documentPath", "setDocumentPath")
    }

    public handle(state: IDocumentState, action: IWrappedAction<ISetDocumentPath>): IDocumentState {
        state = cloneDeep(state);
        if (action.payload) {
            state.documentPath = action.payload.documentPath;
        }
        return state;
    }

    public getSelectorValue(state: IRootState): DocumentPath | null {
        return state.document.documentPath;
    }

}

export const SetDocumentPath: SetDocumentPathImpl = new SetDocumentPathImpl();
