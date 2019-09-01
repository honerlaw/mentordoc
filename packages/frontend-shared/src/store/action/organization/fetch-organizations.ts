import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {IDispatchMap} from "../generic-action";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {SetOrganizations} from "./set-organizations";
import {HttpError} from "../../model/request-status/http-error";
import {AclOrganization} from "../../model/organization/acl-organization";

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

        console.log(orgs);

        api.dispatch(SetOrganizations.action({
            organizations: orgs
        }));
    }


}

export const FetchOrganizations: FetchOrganizationsImpl = new FetchOrganizationsImpl();
