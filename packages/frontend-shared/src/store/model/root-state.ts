import {IRequestStatusState} from "./request-status/request-status-state";
import {IUserState} from "./user/user-state";
import {IAlertState} from "./alert/alert-state";
import {IOrganizationState} from "./organization/organization-state";
import {IFolderState} from "./folder/folder-state";
import {IDocumentState} from "./document/document-state";

export interface IRootState {
    user: IUserState;
    requestStatus: IRequestStatusState;
    alert: IAlertState;
    organization: IOrganizationState;
    folder: IFolderState;
    document: IDocumentState;
}
