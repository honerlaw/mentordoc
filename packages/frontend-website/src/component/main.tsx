import * as React from "react";
import {Route, Switch} from "react-router-dom";
import {LandingPage} from "./marketing/landing-page";
import "./main.scss";
import {SigninPage} from "./auth/signin-page";
import {SignupPage} from "./auth/signup-page";
import {
    CombineDispatchers, CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    AUTHENTICATION_DATA_KEY, IAuthenticationDataDispatch, IAuthenticationDataSelector,
    SetAuthenticationData
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/set-authentication-data";
import {UnsecureRoute} from "./shared/unsecure-route";
import {SecureRoute} from "./shared/secure-route";
import {Dashboard} from "./app/dashboard";
import {
    FetchCurrentUser,
    IFetchCurrentUserDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/fetch-current-user";
import {
    ISetCurrentUserSelector,
    SetCurrentUser
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/set-current-user";
import {
    FetchOrganizations,
    IFetchOrganizationsDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/fetch-organizations";
import {
    ISetOrganizationsSelector,
    SetOrganizations
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-organizations";

interface IProps extends Partial<ISelectorPropMap<IAuthenticationDataSelector & ISetCurrentUserSelector & ISetOrganizationsSelector>
    & IDispatchPropMap<IAuthenticationDataDispatch & IFetchCurrentUserDispatch & IFetchOrganizationsDispatch>> {

}

interface IState {
    isLoading: boolean;
}

@ConnectProps(
    CombineSelectors(SetAuthenticationData.selector, SetCurrentUser.selector, SetOrganizations.selector),
    CombineDispatchers(SetAuthenticationData.dispatch, FetchCurrentUser.dispatch, FetchOrganizations.dispatch)
)
export class Main extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isLoading: true
        };
    }

    public async componentWillMount(): Promise<void> {
        const data: string | null = window.localStorage.getItem(AUTHENTICATION_DATA_KEY);
        if (data) {
            this.props.dispatch!.setAuthenticationData({
                data: JSON.parse(data)
            });
        }

        if (data) {
            await this.props.dispatch!.fetchCurrentUser();
            await this.props.dispatch!.fetchOrganizations();
        }

        this.setState({isLoading: false});
    }

    public render(): JSX.Element | null {
        if (this.state.isLoading || (this.props.selector!.authenticationData && (!this.props.selector!.currentUser || !this.props.selector!.organizations))) {
            return null;
        }

        return <Switch>
            <Route exact path={"/"} component={LandingPage}/>
            <SecureRoute redirect={"/signin"} exact={true} path={"/app/:orgId?/:docId?"} component={Dashboard}/>
            <UnsecureRoute redirect={"/app"} path={"/signin"} component={SigninPage}/>
            <UnsecureRoute redirect={"/app"} path={"/signup"} component={SignupPage}/>
        </Switch>;
    }

}