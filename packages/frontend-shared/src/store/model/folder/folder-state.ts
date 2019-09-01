import {AclFolder} from "./acl-folder";

interface IFolderMap {
    [key: string]: AclFolder[];
}

export interface IFolderState {
    folderMap: IFolderMap;
}

export const INITIAL_FOLDER_STATE: IFolderState = {
    folderMap: {}
};