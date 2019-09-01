import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclDocument} from "../../model/document/acl-document";
import {SetDocuments} from "./set-documents";

export const FETCH_DOCUMENTS_TYPE: string = "fetch_documents_type";

export interface IFetchDocuments extends IGenericActionRequest {
    organizationId: string;
    folderId: string | null;
}

export interface IFetchDocumentsDispatch extends IDispatchMap {
    fetchDocuments: (req?: IFetchDocuments) => Promise<void>;
}

class FetchDocumentsImpl extends AsyncAction<IFetchDocuments> {

    public constructor() {
        super(FETCH_DOCUMENTS_TYPE, "fetchDocuments");
    }

    protected async fetch(api: MiddlewareAPI, req: IFetchDocuments): Promise<void> {
        let path: string = `/document/list/${req.organizationId}`;

        if (req.folderId) {
            path += `?folderId=${req.folderId}`;
        }

        const documents: AclDocument[] | null = await request<AclDocument[]>({
            method: "GET",
            path,
            model: AclDocument,
            api
        });

        if (!documents) {
            throw new HttpError("failed to find documents");
        }

        api.dispatch(SetDocuments.action({
            documents
        }));
    }


}

export const FetchDocuments: FetchDocumentsImpl = new FetchDocumentsImpl();
