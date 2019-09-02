import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {HttpError} from "../../model/request-status/http-error";
import {AclFolder} from "../../model/folder/acl-folder";
import {SetFolders} from "./set-folders";
import {FetchFolders} from "./fetch-folders";
import {IRootState} from "../../model/root-state";

export const CREATE_FOLDER_TYPE: string = "create_folder_type";

export interface ICreateFolder extends IGenericActionRequest {
    organizationId: string;
    parentFolderId: string | null;
    name: string;
}

export interface ICreateFolderDispatch extends IDispatchMap {
    createFolder: (req?: ICreateFolder) => Promise<void>;
}

class CreateFolderImpl extends AsyncAction<ICreateFolder> {

    public constructor() {
        super(CREATE_FOLDER_TYPE, "createFolder");
    }

    protected async fetch(api: MiddlewareAPI, req: ICreateFolder): Promise<void> {
        const folder: AclFolder | null = await request<AclFolder>({
            method: "POST",
            path: "/folder",
            model: AclFolder,
            api,
            body: req
        });

        if (!folder) {
            throw new HttpError("failed to create folder");
        }

        api.dispatch(SetFolders.action({
            folders: [folder]
        }));

        FetchFolders.findParentAndUpdate(api, folder.model.organizationId, folder.model.parentFolderId);
    }


}

export const CreateFolder: CreateFolderImpl = new CreateFolderImpl();
