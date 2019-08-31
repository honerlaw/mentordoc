import {ISelectorMap, SyncAction} from "../sync-action";
import {Organization} from "../../model/organization/organization";
import {IOrganizationState} from "../../model/organization/organization-state";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";

export const SET_ORGANIZATIONS_TYPE: string = "set_organizations_type";

export interface ISetOrganizations {
    organizations: Organization[];
}

export interface ISetOrganizationsSelector extends ISelectorMap {
    organizations: Organization[]
}

export interface ISetOrganizationsDispatch extends IDispatchMap {
    setOrganizations: (req?: ISetOrganizations) => void;
}

export class SetOrganizationsImpl extends SyncAction<IOrganizationState, ISetOrganizations, Organization[]> {

    public constructor() {
        super(ReducerType.ORGANIZATION, SET_ORGANIZATIONS_TYPE, "organizations", "setOrganizations")
    }

    public getSelectorValue(state: IRootState): Organization[] {
        return state.organization.organizations;
    }

}

export const SetOrganizations: SetOrganizationsImpl = new SetOrganizationsImpl();
