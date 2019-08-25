import {AuthenticationData} from "./authentication-data";
import {User} from "./user";

export interface IUserState {
    currentUser: User | null;
    authenticationData: AuthenticationData | null;
}

export const USER_INITIAL_STATE: IUserState = {
    currentUser: null,
    authenticationData: null,
};
