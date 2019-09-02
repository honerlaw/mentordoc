import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {MiddlewareAPI} from "redux";
import {AclDocument} from "../../model/document/acl-document";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {SetDocuments} from "./set-documents";
import {SetFullDocument} from "./set-full-document";

export const FETCH_FULL_DOCUMENT_TYPE: string = "fetch_full_document_type";

export interface IFetchFullDocument extends IGenericActionRequest {
    documentId: string;
}

export interface IFetchFullDocumentDispatch {
    fetchFullDocument: (req: IFetchFullDocument) => Promise<void>;
}

class FetchFullDocumentImpl extends AsyncAction<IFetchFullDocument> {

    public constructor() {
        super(FETCH_FULL_DOCUMENT_TYPE, "fetchFullDocument");
    }

    protected async fetch(api: MiddlewareAPI, req: IFetchFullDocument): Promise<void> {
        const document: AclDocument | null = await request<AclDocument>({
            method: "GET",
            path: `/document/${req.documentId}`,
            model: AclDocument,
            api
        });

        if (!document) {
            throw new HttpError("failed to find documents");
        }

        api.dispatch(SetFullDocument.action({
            fullDocument: document
        }));
    }

}

export const FetchFullDocument: FetchFullDocumentImpl = new FetchFullDocumentImpl();
