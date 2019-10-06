import * as React from "react";
import {Page} from "../shared/page";
import "./organization-page.scss";

export class OrganizationPage extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <Page>
            <div id={"organization-page"}>
                <div className={"settings"}>
                    org
                </div>
            </div>
        </Page>;
    }

}