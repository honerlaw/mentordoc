import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {IFolderState} from "../../model/folder/folder-state";
import {AclFolder} from "../../model/folder/acl-folder";

export const SET_FOLDERS_TYPE: string = "set_folders_type";

export type SelectorValue = (orgId: string, parentFolderId: string | null) => AclFolder[];

export interface ISetFolders {
    folders: AclFolder[];
}

export interface ISetFoldersSelector extends ISelectorMap {
    getFolders: SelectorValue;
}

export interface ISetFoldersDispatch extends IDispatchMap {
    setFolders: (req?: ISetFolders) => void;
}

export class SetFoldersImpl extends SyncAction<IFolderState, ISetFolders, SelectorValue> {

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

                // otherwise lets make sure the folder isn't in the map already, and then add it
                const found: AclFolder | undefined = folders
                    .find((existing: AclFolder): boolean => existing.model.id === folder.model.id);
                if (!found) {
                    folders.push(folder);
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorValue {
        return (orgId: string, parentFolderId: string | null): AclFolder[] => {
            const folders: AclFolder[] | undefined = state.folder.folderMap[this.getKey(orgId, parentFolderId)];
            if (folders) {
                return folders;
            }
            return [];
        };
    }

    private getKeyForFolder(folder: AclFolder): string {
        return this.getKey(folder.model.organizationId, folder.model.parentFolderId);
    }

    private getKey(orgId: string, parentFolderId: string | null): string {
        if (parentFolderId) {
            return `${orgId}-${parentFolderId}`;
        }
        return orgId;
    }

}

export const SetFolders: SetFoldersImpl = new SetFoldersImpl();
