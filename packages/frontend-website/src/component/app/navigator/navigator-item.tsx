import * as React from "react";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";

interface IProps {
    item: AclOrganization;
}

export class NavigatorItem extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div>{this.props.item.model.name}</div>
    }

}