import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclDocument} from "../../model/document/acl-document";
import {UnsetDocuments} from "./unset-documents";
import {FetchDocuments} from "./fetch-documents";
import {FetchFolders} from "../folder/fetch-folders";

export const DELETE_DOCUMENT_TYPE: string = "delete_document_type";

export interface IDeleteDocument extends IGenericActionRequest {
    documentId: string;
}

export interface IDeleteDocumentDispatch extends IDispatchMap {
    deleteDocument: (req: IDeleteDocument) => Promise<void>;
}

class DeleteDocumentImpl extends AsyncAction<IDeleteDocument> {

    public constructor() {
        super(DELETE_DOCUMENT_TYPE, "deleteDocument");
    }

    protected async fetch(api: MiddlewareAPI, req: IDeleteDocument): Promise<void> {
        const document: AclDocument | null = await request<AclDocument>({
            method: "DELETE",
            path: `/document/${req.documentId}`,
            model: AclDocument,
            api,
            body: req
        });

        if (!document) {
            throw new HttpError("failed to create document");
        }

        api.dispatch(UnsetDocuments.action({
            documents: [document]
        }));

        api.dispatch(FetchDocuments.action({
            organizationId: document.model.organizationId,
            folderId: document.model.folderId
        }) as any);
        FetchFolders.findParentAndUpdate(api, document.model.organizationId, document.model.folderId);
    }

}

export const DeleteDocument: DeleteDocumentImpl = new DeleteDocumentImpl();
