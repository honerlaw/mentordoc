import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclDocument} from "../../model/document/acl-document";
import {SetDocuments} from "./set-documents";
import {FetchFolders} from "../folder/fetch-folders";
import {FetchDocuments} from "./fetch-documents";
import {SetFullDocument} from "./set-full-document";
import {FetchFullDocument} from "./fetch-full-document";

export const CREATE_DOCUMENT_DRAFT_TYPE: string = "create_document_draft_type";

export interface ICreateDocumentDraft extends IGenericActionRequest {
    documentId: string;
    name: string;
    content: string;
}

export interface ICreateDocumentDraftDispatch extends IDispatchMap {
    createDocumentDraft: (req?: ICreateDocumentDraft) => Promise<void>;
}

class CreateDocumentDraftImpl extends AsyncAction<ICreateDocumentDraft> {

    public constructor() {
        super(CREATE_DOCUMENT_DRAFT_TYPE, "createDocumentDraft");
    }

    protected async fetch(api: MiddlewareAPI, req: ICreateDocumentDraft): Promise<void> {
        const document: AclDocument | null = await request<AclDocument>({
            method: "POST",
            path: "/document/draft",
            model: AclDocument,
            api,
            body: req
        });

        if (!document) {
            throw new HttpError("failed to create document draft");
        }

        api.dispatch(SetDocuments.action({
            documents: [document]
        }));

        api.dispatch(FetchFullDocument.action({
            documentId: document.model.id
        }) as any);

        api.dispatch(FetchDocuments.action({
            organizationId: document.model.organizationId,
            folderId: document.model.folderId
        }) as any);
        FetchFolders.findParentAndUpdate(api, document.model.organizationId, document.model.folderId);
    }


}

export const CreateDocumentDraft: CreateDocumentDraftImpl = new CreateDocumentDraftImpl();
