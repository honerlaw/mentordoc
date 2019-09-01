import * as React from "react";
import {Page} from "../shared/page";
import {Navigator} from "./navigator";
import {DocumentViewer} from "./document-viewer";
import "./dashboard.scss";

export class Dashboard extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <Page>
            <div id={"dashboard-container"}>
                <Navigator/>
                <DocumentViewer/>
            </div>
        </Page>;
    }

}
