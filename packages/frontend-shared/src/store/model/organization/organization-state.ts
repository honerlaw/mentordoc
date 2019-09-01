import {AclOrganization} from "./acl-organization";

export interface IOrganizationState {
    organizations: AclOrganization[];
}

export const INITIAL_ORGANIZATION_STATE: IOrganizationState = {
    organizations: []
};