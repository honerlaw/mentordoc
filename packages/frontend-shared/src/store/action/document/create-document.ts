import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclDocument} from "../../model/document/acl-document";
import {SetDocuments} from "./set-documents";

export const CREATE_DOCUMENT_TYPE: string = "create_document_type";

export interface ICreateDocument extends IGenericActionRequest {
    organizationId: string;
    folderId: string | null;
    name: string;
    content: string;
}

export interface ICreateDocumentDispatch extends IDispatchMap {
    createDocument: (req?: ICreateDocument) => Promise<void>;
}

class CreateDocumentImpl extends AsyncAction<ICreateDocument> {

    public constructor() {
        super(CREATE_DOCUMENT_TYPE, "createDocument");
    }

    protected async fetch(api: MiddlewareAPI, req: ICreateDocument): Promise<void> {
        const document: AclDocument | null = await request<AclDocument>({
            method: "POST",
            path: "/document",
            model: AclDocument,
            api,
            body: req
        });

        if (!document) {
            throw new HttpError("failed to create document");
        }

        api.dispatch(SetDocuments.action({
            documents: [document]
        }));
    }


}

export const CreateDocument: CreateDocumentImpl = new CreateDocumentImpl();
