import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {AclDocument} from "../../model/document/acl-document";
import {SetFullDocument} from "./set-full-document";

export const UPDATE_DOCUMENT_TYPE: string = "update_document_type";

export interface IUpdateDocument extends IGenericActionRequest {
    documentId: string;
    draftId: string;
    name: string;
    content: string;
}

export interface IUpdateDocumentDispatch extends IDispatchMap {
    updateDocument: (req?: IUpdateDocument) => Promise<void>;
}

class UpdateDocumentImpl extends AsyncAction<IUpdateDocument> {

    public constructor() {
        super(UPDATE_DOCUMENT_TYPE, "updateDocument");
    }

    protected async fetch(api: MiddlewareAPI, req: IUpdateDocument): Promise<void> {
        const document: AclDocument = await request<AclDocument>({
            method: "PUT",
            path: "/document",
            model: AclDocument,
            api,
            body: req
        });

        api.dispatch(SetFullDocument.action({
            fullDocument: document
        }));
    }


}

export const UpdateDocument: UpdateDocumentImpl = new UpdateDocumentImpl();
