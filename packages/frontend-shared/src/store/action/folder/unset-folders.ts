import {ISelectorMap, SyncAction} from "../sync-action";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {IFolderState} from "../../model/folder/folder-state";
import {AclFolder} from "../../model/folder/acl-folder";

export const UNSET_FOLDERS_TYPE: string = "unset_folders_type";

export type SelectorValue = (orgId: string, parentFolderId: string | null) => AclFolder[];

export interface IUnsetFolders {
    folders: AclFolder[];
}

export interface IUnsetFoldersDispatch extends IDispatchMap {
    unsetFolders: (req?: IUnsetFolders) => void;
}

export class UnsetFoldersImpl extends SyncAction<IFolderState, IUnsetFolders, void> {

    public constructor() {
        super(ReducerType.FOLDER, UNSET_FOLDERS_TYPE, "unsetFolders", "unsetFolders")
    }

    public handle(state: IFolderState, action: IWrappedAction<IUnsetFolders>): IFolderState {
        state = cloneDeep(state);
        if (action.payload) {
            action.payload.folders.forEach((folder: AclFolder): void => {
                const key: string = this.getKeyForFolder(folder);

                const folders: AclFolder[] | undefined = state.folderMap[key];

                // no folders exist for the key, so nothing to remove
                if (!folders) {
                    return;
                }

                // remove the folder
                for (let i: number = folders.length - 1; i >= 0; --i) {
                    if (folders[i].model.id === folder.model.id) {
                        folders.splice(i, 1);
                        break;
                    }
                }
            });
        }
        return state;
    }

    public getSelectorValue(state: IRootState): void {
        // does nothing!
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

export const UnsetFolders: UnsetFoldersImpl = new UnsetFoldersImpl();
