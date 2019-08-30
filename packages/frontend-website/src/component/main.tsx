import * as React from "react";
import {Route, Switch} from "react-router-dom";
import {LandingPage} from "./marketing/landing-page";
import "./main.scss";
import {SigninPage} from "./auth/signin-page";
import {SignupPage} from "./auth/signup-page";
import {
    CombineDispatchers,
    ConnectProps, IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    AUTHENTICATION_DATA_KEY, IAuthenticationDataDispatch,
    SetAuthenticationData
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/set-authentication-data";
import {UnsecureRoute} from "./shared/unsecure-route";
import {SecureRoute} from "./shared/secure-route";
import {Dashboard} from "./app/dashboard";

interface IProps extends Partial<IDispatchPropMap<IAuthenticationDataDispatch>> {

}

interface IState {
    isLoading: boolean;
}

@ConnectProps(null, CombineDispatchers(SetAuthenticationData.dispatch))
export class Main extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isLoading: true
        };
    }

    public componentWillMount(): void {
        const data: string | null = window.localStorage.getItem(AUTHENTICATION_DATA_KEY)
        if (data) {
            this.props.dispatch!.setAuthenticationData({
                data: JSON.parse(data)
            });
        }

        this.setState({
            isLoading: false
        });
    }

    public render(): JSX.Element | null {
        if (this.state.isLoading) {
            return null;
        }

        return <Switch>
            <Route exact path={"/"} component={LandingPage}/>
            <SecureRoute redirect={"/signin"} exact={true} path={"/app"} component={Dashboard}/>
            <UnsecureRoute redirect={"/app"} path={"/signin"} component={SigninPage}/>
            <UnsecureRoute redirect={"/app"} path={"/signup"} component={SignupPage}/>
        </Switch>;
    }

}