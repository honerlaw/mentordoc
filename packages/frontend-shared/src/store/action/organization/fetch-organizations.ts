import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {SetOrganizations} from "./set-organizations";
import {HttpError} from "../../model/request-status/http-error";
import {AclOrganization} from "../../model/organization/acl-organization";
import {IRootState} from "../../model/root-state";
import {SetCurrentOrganization} from "./set-current-organization";

export const FETCH_ORGANIZATIONS_TYPE: string = "fetch_organizations_type";

export interface IFetchOrganizations extends IGenericActionRequest {

}

export interface IFetchOrganizationsDispatch extends IDispatchMap {
    fetchOrganizations: (req?: IFetchOrganizations) => Promise<void>;
}

export class FetchOrganizationsImpl extends AsyncAction<IFetchOrganizations> {

    public constructor() {
        super(FETCH_ORGANIZATIONS_TYPE, "fetchOrganizations");
    }

    protected async fetch(api: MiddlewareAPI, req?: IFetchOrganizations): Promise<void> {
        const orgs: AclOrganization[] | null = await request<AclOrganization[]>({
            method: "GET",
            path: "/organization/list",
            model: AclOrganization,
            api
        });

        if (!orgs) {
            throw new HttpError("failed to find organizations");
        }

        // set the organizations
        api.dispatch(SetOrganizations.action({
            organizations: orgs
        }));

        // default to the first org in the list
        let newCurrentOrganization: AclOrganization = orgs[0];

        // get the current organization if there is one
        const rootState: IRootState = api.getState();
        const currentOrganization: AclOrganization | null = rootState.organization.currentOrganization;

        // the current org is set, exists in the list, and is different from the default value, so use that instead
        if (currentOrganization
            && this.isOrganizationInList(currentOrganization, orgs)
            && currentOrganization.model.id !== newCurrentOrganization.model.id) {
            newCurrentOrganization = currentOrganization;
        } else {

            // fallback to the local org in storage instead
            // if the local org exists, is in the list of orgs, and is not the same as the current local org
            const localCurrentOrganization: AclOrganization | null = SetCurrentOrganization.getFromStorage();
            if (localCurrentOrganization
                && this.isOrganizationInList(localCurrentOrganization, orgs)
                && newCurrentOrganization.model.id !== localCurrentOrganization.model.id) {
                newCurrentOrganization = localCurrentOrganization;
            }
        }

        // otherwise it doesn't, so default to the first org they have in their list
        api.dispatch(SetCurrentOrganization.action({
            currentOrganization: newCurrentOrganization
        }));
    }

    private isOrganizationInList(organization: AclOrganization | null, orgs: AclOrganization[]): boolean {
        if (organization === null) {
            return false;
        }
        return orgs.some((org: AclOrganization): boolean => org.model.id === organization.model.id);
    }


}

export const FetchOrganizations: FetchOrganizationsImpl = new FetchOrganizationsImpl();
