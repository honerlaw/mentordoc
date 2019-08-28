import * as React from "react";
import {Header} from "./header";
import {AlertList} from "./alert-list";

export class Page extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <div className={"page"}>
            <Header />
            <AlertList/>
            {this.props.children}
        </div>;
    }

}
