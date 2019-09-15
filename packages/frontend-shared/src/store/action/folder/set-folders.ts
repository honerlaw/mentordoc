import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {IFolderState} from "../../model/folder/folder-state";
import {AclFolder} from "../../model/folder/acl-folder";

export const SET_FOLDERS_TYPE: string = "set_folders_type";

export type SelectorSetFoldersChildValue = (orgId: string, parentFolderId: string | null) => AclFolder[];
export type SelectorSetFoldersMapValue = Record<string, AclFolder[]>;
export type SelectorSetFoldersValue = (type: "child" | "map") => SelectorSetFoldersChildValue | SelectorSetFoldersMapValue;

export interface ISetFolders {
    folders: AclFolder[];
}

export interface ISetFoldersSelector extends ISelectorMap {
    getFolders: SelectorSetFoldersValue;
}

export interface ISetFoldersDispatch extends IDispatchMap {
    setFolders: (req?: ISetFolders) => void;
}

export class SetFoldersImpl extends SyncAction<IFolderState, ISetFolders, SelectorSetFoldersValue> {

    public constructor() {
        super(ReducerType.FOLDER, SET_FOLDERS_TYPE, "getFolders", "setFolders")
    }

    public handle(state: IFolderState, action: IWrappedAction<ISetFolders>): IFolderState {
        state = cloneDeep(state);
        if (action.payload) {
            action.payload.folders.forEach((folder: AclFolder): void => {
                const key: string = this.getKeyForFolder(folder);

                const folders: AclFolder[] | undefined = state.folderMap[key];

                // no folders exist for the key, so just add it
                if (!folders) {
                    state.folderMap[key] = [folder];
                    return;
                }

                let found: boolean = false;

                // see if the folder exists already
                for (let i: number = folders.length - 1; i >= 0; --i) {
                    const f: AclFolder = folders[i];

                    // if it does exist, replace it and mark that we found it
                    if (f.model.id === folder.model.id) {
                        found = true;
                        folders[i] = folder;
                    }
                }

                // if it wasn't found then add it
                if (!found) {
                    folders.push(folder);
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorSetFoldersValue {
        return (type: "child" | "map"): SelectorSetFoldersChildValue | SelectorSetFoldersMapValue => {
            if (type === "map") {
                return state.folder.folderMap;
            }
            return (orgId: string, parentFolderId: string | null): AclFolder[] => {
                const folders: AclFolder[] | undefined = state.folder.folderMap[this.getKey(orgId, parentFolderId)];
                if (folders) {
                    return folders;
                }
                return [];
            };
        };
    }

    public getKeyForFolder(folder: AclFolder): string {
        return this.getKey(folder.model.organizationId, folder.model.parentFolderId);
    }

    public getKey(orgId: string, parentFolderId: string | null): string {
        if (parentFolderId) {
            return `${orgId}:${parentFolderId}`;
        }
        return orgId;
    }

}

export const SetFolders: SetFoldersImpl = new SetFoldersImpl();
