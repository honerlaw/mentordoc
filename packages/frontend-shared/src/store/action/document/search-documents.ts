import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclDocument} from "../../model/document/acl-document";
import {SetSearchDocuments} from "./set-search-documents";

export const SEARCH_DOCUMENTS_TYPE: string = "search_documents_type";

export interface ISearchDocuments extends IGenericActionRequest {
    searchQuery: string;
}

export interface ISearchDocumentsDispatch extends IDispatchMap {
    searchDocuments: (req?: ISearchDocuments) => Promise<void>;
}

class SearchDocumentsImpl extends AsyncAction<ISearchDocuments> {

    public constructor() {
        super(SEARCH_DOCUMENTS_TYPE, "searchDocuments");
    }

    protected async fetch(api: MiddlewareAPI, req: ISearchDocuments): Promise<void> {
        let path: string = `/document/search?query=${req.searchQuery}`;

        const documents: AclDocument[] | null = await request<AclDocument[]>({
            method: "GET",
            path,
            model: AclDocument,
            api
        });

        if (!documents) {
            throw new HttpError("failed to find documents");
        }

        api.dispatch(SetSearchDocuments.action({
            searchDocuments: documents
        }));
    }


}

export const SearchDocuments: SearchDocumentsImpl = new SearchDocumentsImpl();
