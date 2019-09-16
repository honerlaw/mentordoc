import * as React from "react";
import {Page} from "../shared/page";
import "./account-page.scss";

export class AccountPage extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <Page>
            <div id={"account-page"}>

                <div className={"settings"}>
                    hello
                </div>


            </div>
        </Page>;
    }

}