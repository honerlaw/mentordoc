import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclFolder} from "../../model/folder/acl-folder";
import {SetFolders} from "./set-folders";
import {IRootState} from "../../model/root-state";

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

    public findParentAndUpdate(api: MiddlewareAPI, organizationId: string, parentFolderId: string | null): void {
        // basically we need to refetch not the parents' children, but the parents' parents' children instead
        let parentId: string | null = null;
        const state: IRootState = api.getState();
        Object.values(state.folder.folderMap).forEach((folders: AclFolder[]): void => {
            const parent: AclFolder | undefined = folders
                .find((f: AclFolder): boolean => f.model.id === parentFolderId);
            if (parent) {
                parentId = parent.model.parentFolderId;
            }
        });

        api.dispatch(FetchFolders.action({
            organizationId: organizationId,
            parentFolderId: parentId
        }) as any);
    }


}

export const FetchFolders: FetchFoldersImpl = new FetchFoldersImpl();
