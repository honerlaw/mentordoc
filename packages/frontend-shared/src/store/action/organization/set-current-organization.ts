import {ISelectorMap, SyncAction} from "../sync-action";
import {IOrganizationState} from "../../model/organization/organization-state";
import {AclOrganization} from "../../model/organization/acl-organization";
import {IRootState} from "../../model/root-state";
import {IDispatchMap} from "../generic-action";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {plainToClass} from "class-transformer";

const SET_CURRENT_ORGANIZATION_TYPE: string = "set_current_organization_type";

export interface ISetCurrentOrganization {
    currentOrganization: AclOrganization | null;
}

export interface ISetCurrentOrganizationSelector extends ISelectorMap {
    currentOrganization: AclOrganization | null
}

export interface ISetOrganizationsDispatch extends IDispatchMap {
    setCurrentOrganization: (req?: ISetCurrentOrganization) => void;
}

class SetCurrentOrganizationImpl extends SyncAction<IOrganizationState, ISetCurrentOrganization, AclOrganization> {

    public constructor() {
        super(ReducerType.ORGANIZATION, SET_CURRENT_ORGANIZATION_TYPE, "currentOrganization", "setCurrentOrganization");
    }

    public handle(state: IOrganizationState, action: IWrappedAction<ISetCurrentOrganization>): IOrganizationState {
        state = cloneDeep(state);
        if (action.payload) {
            state.currentOrganization = action.payload.currentOrganization;
            this.saveToStorage(state.currentOrganization);
        }
        return state;
    }

    public getSelectorValue(state: IRootState): AclOrganization | null {
        return state.organization.currentOrganization;
    }

    public saveToStorage(org: AclOrganization | null): void {
        if (org === null) {
            window.localStorage.removeItem(SET_CURRENT_ORGANIZATION_TYPE);
        } else {
            window.localStorage.setItem(SET_CURRENT_ORGANIZATION_TYPE, JSON.stringify(org));
        }
    }

    public getFromStorage(): AclOrganization | null {
        const payload: string | null = window.localStorage.getItem(SET_CURRENT_ORGANIZATION_TYPE);
        if (!payload) {
            return null;
        }

        return plainToClass(AclOrganization, JSON.parse(payload));
    }

}

export const SetCurrentOrganization: SetCurrentOrganizationImpl = new SetCurrentOrganizationImpl();
