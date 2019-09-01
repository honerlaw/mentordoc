import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclFolder} from "../../model/folder/acl-folder";
import {SetFolders} from "./set-folders";

export const FETCH_FOLDERS_TYPE: string = "fetch_folders_type";

export interface IFetchFolders extends IGenericActionRequest {
    organizationId: string;
    parentFolderId: string | null;
}

export interface IFetchFoldersDispatch extends IDispatchMap {
    fetchFolders: (req?: IFetchFolders) => Promise<void>;
}

class FetchFoldersImpl extends AsyncAction<IFetchFolders> {

    public constructor() {
        super(FETCH_FOLDERS_TYPE, "fetchFolders");
    }

    protected async fetch(api: MiddlewareAPI, req: IFetchFolders): Promise<void> {
        let path: string = `/folder/list/${req.organizationId}`;

        if (req.parentFolderId) {
            path += `?parentFolderId=${req.parentFolderId}`;
        }

        const folders: AclFolder[] | null = await request<AclFolder[]>({
            method: "GET",
            path,
            model: AclFolder,
            api
        });

        if (!folders) {
            throw new HttpError("failed to find folders");
        }

        api.dispatch(SetFolders.action({
            folders
        }));
    }


}

export const FetchFolders: FetchFoldersImpl = new FetchFoldersImpl();
