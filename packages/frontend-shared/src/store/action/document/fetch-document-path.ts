import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {AclDocument} from "../../model/document/acl-document";
import {DocumentPath} from "../../model/document/document-state";
import {SetDocumentPath} from "./set-document-path";

export const FETCH_DOCUMENT_PATH_TYPE: string = "fetch_document_path_type";

export interface IFetchDocumentPath extends IGenericActionRequest {
    documentId: string;
}

export interface IFetchDocumentPathDispatch extends IDispatchMap {
    fetchDocumentPath: (req?: IFetchDocumentPath) => Promise<void>;
}

class FetchDocumentPathImpl extends AsyncAction<IFetchDocumentPath> {

    public constructor() {
        super(FETCH_DOCUMENT_PATH_TYPE, "fetchDocumentPath");
    }

    protected async fetch(api: MiddlewareAPI, req: IFetchDocumentPath): Promise<void> {
        const resp: Response = await request<Response>({
            method: "GET",
            path: `/document/path/${req.documentId}`,
            api
        });

        const documentPath: DocumentPath = await resp.json();
        api.dispatch(SetDocumentPath.action({
            documentPath
        }));
    }


}

export const FetchDocumentPath: FetchDocumentPathImpl = new FetchDocumentPathImpl();
