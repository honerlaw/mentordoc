import {AclOrganization} from "./acl-organization";

export interface IOrganizationState {
    organizations: AclOrganization[] | null
    currentOrganization: AclOrganization | null;
}

export const INITIAL_ORGANIZATION_STATE: IOrganizationState = {
    organizations: null,
    currentOrganization: null
};