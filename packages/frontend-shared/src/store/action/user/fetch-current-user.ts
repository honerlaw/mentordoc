import {AsyncAction} from "../async-action";
import {IGenericActionRequest} from "../generic-action-request";
import {MiddlewareAPI} from "redux";
import {request} from "../../../util/request";
import {User} from "../../model/user/user";
import {SetCurrentUser} from "./set-current-user";
import {IDispatchMap} from "../generic-action";

export const FETCH_CURRENT_USER_TYPE: string = "fetch_current_user_type";

export interface IFetchCurrentUser extends IGenericActionRequest {

}

export interface IFetchCurrentUserDispatch extends IDispatchMap {
    fetchCurrentUser: (req?: IFetchCurrentUser) => Promise<void>;
}

export class FetchCurrentUserImpl extends AsyncAction<IFetchCurrentUser> {

    public constructor() {
        super(FETCH_CURRENT_USER_TYPE, "fetchCurrentUser");
    }

    protected async fetch(api: MiddlewareAPI, req?: IFetchCurrentUser): Promise<void> {
        const user: User | null = await request({
            method: "GET",
            path: "/user",
            model: User,
            api
        });

        api.dispatch(SetCurrentUser.action({
            currentUser: user
        }));
    }

}

export const FetchCurrentUser: FetchCurrentUserImpl = new FetchCurrentUserImpl();
