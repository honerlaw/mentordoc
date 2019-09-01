import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclFolder} from "../../model/folder/acl-folder";
import {SetFolders} from "./set-folders";
import {UnsetFolders} from "./unset-folders";

export const DELETE_FOLDER_TYPE: string = "delete_folder_type";

export interface IDeleteFolder extends IGenericActionRequest {
    folderId: string;
}

export interface IDeleteFolderDispatch extends IDispatchMap {
    deleteFolder: (req?: IDeleteFolder) => Promise<void>;
}

class DeleteFolderImpl extends AsyncAction<IDeleteFolder> {

    public constructor() {
        super(DELETE_FOLDER_TYPE, "deleteFolder");
    }

    protected async fetch(api: MiddlewareAPI, req: IDeleteFolder): Promise<void> {
        const folder: AclFolder | null = await request<AclFolder>({
            method: "DELETE",
            path:`/folder/${req.folderId}`,
            model: AclFolder,
            api,
            body: req
        });

        if (!folder) {
            throw new HttpError("failed to delete folder");
        }

        api.dispatch(UnsetFolders.action({
            folders: [folder]
        }));
    }

}

export const DeleteFolder: DeleteFolderImpl = new DeleteFolderImpl();
