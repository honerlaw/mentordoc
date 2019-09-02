import {AclOrganization} from "./acl-organization";

export interface IOrganizationState {
    organizations: AclOrganization[] | null
}

export const INITIAL_ORGANIZATION_STATE: IOrganizationState = {
    organizations: null
};