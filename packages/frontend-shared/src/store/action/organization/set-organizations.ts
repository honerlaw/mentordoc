import {ISelectorMap, SyncAction} from "../sync-action";
import {IOrganizationState} from "../../model/organization/organization-state";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {AclOrganization} from "../../model/organization/acl-organization";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";

export const SET_ORGANIZATIONS_TYPE: string = "set_organizations_type";

export interface ISetOrganizations {
    organizations: AclOrganization[] | null;
}

export interface ISetOrganizationsSelector extends ISelectorMap {
    organizations: AclOrganization[] | null
}

export interface ISetOrganizationsDispatch extends IDispatchMap {
    setOrganizations: (req?: ISetOrganizations) => void;
}

export class SetOrganizationsImpl extends SyncAction<IOrganizationState, ISetOrganizations, AclOrganization[]> {

    public constructor() {
        super(ReducerType.ORGANIZATION, SET_ORGANIZATIONS_TYPE, "organizations", "setOrganizations")
    }

    public handle(state: IOrganizationState, action: IWrappedAction<ISetOrganizations>): IOrganizationState {
        state = cloneDeep(state);
        if (action.payload) {
            state.organizations = action.payload.organizations
        }
        return state;
    }

    public getSelectorValue(state: IRootState): AclOrganization[] | null {
        return state.organization.organizations;
    }

}

export const SetOrganizations: SetOrganizationsImpl = new SetOrganizationsImpl();
