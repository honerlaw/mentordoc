import * as React from "react";
import "./header.scss";
import {Link} from "react-router-dom";
import {
    CombineDispatchers,
    CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {ILogoutDispatch, Logout} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/logout";
import {
    IAuthenticationDataSelector,
    SetAuthenticationData
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/set-authentication-data";

interface IProps extends Partial<IDispatchPropMap<ILogoutDispatch> & ISelectorPropMap<IAuthenticationDataSelector>> {

}

@ConnectProps(CombineSelectors(SetAuthenticationData.selector), CombineDispatchers(Logout.dispatch))
export class Header extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"header"}>
            <div>
                <Link className={"logo"} to={"/"}>mentordoc</Link>
            </div>

            <div className={"options"}>
                {this.renderOptions()}
            </div>
        </div>;
    }

    private renderOptions(): JSX.Element[] {
        const options: JSX.Element[] = [];
        if (this.props.selector!.authenticationData) {
            options.push(<div key={"logout"} className={"option"}
                              onClick={() => this.props.dispatch!.logout()}>logout</div>);
        }
        return options;
    }

}