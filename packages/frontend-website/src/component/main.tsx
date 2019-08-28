import * as React from "react";
import {Route, Switch} from "react-router-dom";
import {LandingPage} from "./marketing/landing-page";
import "./main.scss";
import {SigninPage} from "./auth/signin-page";
import {SignupPage} from "./auth/signup-page";

export class Main extends React.PureComponent<{}, {}> {

    public render(): JSX.Element | null {
        return <Switch>
            <Route exact path={"/"} component={LandingPage} />
            <Route path={"/signin"} component={SigninPage}/>
            <Route path={"/signup"} component={SignupPage}/>
        </Switch>;
    }

}