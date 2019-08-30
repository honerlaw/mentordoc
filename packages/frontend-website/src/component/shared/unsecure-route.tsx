import * as React from "react";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    IAuthenticationDataSelector,
    SetAuthenticationData
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/set-authentication-data";
import {Redirect, Route, RouteProps} from "react-router";

interface IProps extends Partial<ISelectorPropMap<IAuthenticationDataSelector>>, RouteProps {
    redirect: string;
}

@ConnectProps(CombineSelectors(SetAuthenticationData.selector))
export class UnsecureRoute extends React.PureComponent<IProps, {}> {

    public constructor(props: IProps) {
        super(props);
    }

    public render(): JSX.Element {
        // they are logged in, so redirect
        if (this.props.selector!.authenticationData) {
            return <Redirect to={this.props.redirect} />;
        }
        return <Route {...this.props} />;
    }

}
