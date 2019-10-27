import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {DocumentPath} from "../../model/document/document-state";
import {SetDocumentPath} from "./set-document-path";
import {FetchFolders} from "../folder/fetch-folders";
import {FetchDocuments} from "./fetch-documents";
import {AclOrganization, isAclOrganization} from "../../model/organization/acl-organization";
import {AclFolder, isAclFolder} from "../../model/folder/acl-folder";
import {AclDocument} from "../../model/document/acl-document";
import {SetFolders} from "../folder/set-folders";
import {plainToClass} from "class-transformer";

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

        this.dispatchFolders(api, documentPath);

        api.dispatch(SetDocumentPath.action({
            documentPath
        }));

        this.dispatchRequests(api, documentPath);
    }

    private dispatchFolders(api: MiddlewareAPI, path: DocumentPath): void {
        const folders: AclFolder[] = [];

        path.forEach((item: AclOrganization | AclFolder | AclDocument): void => {
            if (!item.model.id) {
                return;
            }
            if (isAclFolder(item)) {
                folders.push(plainToClass(AclFolder, item));
            }
        });

        api.dispatch(SetFolders.action({
            folders
        }));
    }

    private dispatchRequests(api: MiddlewareAPI, path: DocumentPath): void {
        // @todo potentially make this a flag of some sort? or a different action altogether, might be crappy to
        // @todo fetch this data every single time we get the document path...
        path.forEach((item: AclOrganization | AclFolder | AclDocument): void => {
            if (!item.model.id) {
                return;
            }

            if (isAclOrganization(item)) {
                api.dispatch(FetchFolders.action({
                    organizationId: item.model.id,
                    parentFolderId: null
                }) as any);
                api.dispatch(FetchDocuments.action({
                    organizationId: item.model.id,
                    folderId: null
                }) as any);
            }

            if (isAclFolder(item)) {
                api.dispatch(FetchFolders.action({
                    organizationId: item.model.organizationId,
                    parentFolderId: item.model.parentFolderId
                }) as any);
                api.dispatch(FetchDocuments.action({
                    organizationId: item.model.organizationId,
                    folderId: item.model.id
                }) as any);
            }
        });
    }

}

export const FetchDocumentPath: FetchDocumentPathImpl = new FetchDocumentPathImpl();
