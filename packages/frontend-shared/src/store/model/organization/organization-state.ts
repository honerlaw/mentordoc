import {Organization} from "./organization";

export interface IOrganizationState {
    organizations: Organization[];
}

export const INITIAL_ORGANIZATION_STATE: IOrganizationState = {
    organizations: []
};