import * as React from "react";
import "./header.scss";
import {Link} from "react-router-dom";

export class Header extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <div id={"header"}>
            <Link className={"logo"} to={"/"}>mentordoc</Link>
        </div>;
    }

}