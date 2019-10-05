import * as React from "react";
import "./header.scss";
import {Link, RouteComponentProps} from "react-router-dom";
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
import {Avatar} from "./avatar";
import {DropdownButton} from "./dropdown-button";
import {WithRouter} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/with-router";

interface IProps extends Partial<IDispatchPropMap<ILogoutDispatch> & ISelectorPropMap<IAuthenticationDataSelector> & RouteComponentProps> {

}

@WithRouter()
@ConnectProps(CombineSelectors(SetAuthenticationData.selector), CombineDispatchers(Logout.dispatch))
export class Header extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"header"}>
            <div id={"header-container"}>
                <div>
                    <Link className={"logo"} to={"/"}>mentordoc</Link>
                </div>

                <div className={"search"}>
                    <input type={"text"} placeholder={"search"} />
                </div>

                <div className={"options"}>
                    {this.renderOptions()}
                </div>
            </div>
        </div>;
    }

    private renderOptions(): JSX.Element[] {
        const options: JSX.Element[] = [];
        if (this.props.selector!.authenticationData) {
            options.push(this.renderDropdown());
        }
        return options;
    }

    private renderDropdown(): JSX.Element {
        return <DropdownButton key={"avatar"} label={<Avatar label={"test"}/>} position={"bottom"} options={
            [
                {
                    label: "account",
                    onClick: () => this.props.history!.push("/app/account")
                },
                {
                    label: "sign out",
                    onClick: () => this.props.dispatch!.logout()
                }
            ]
        }/>;
    }

}