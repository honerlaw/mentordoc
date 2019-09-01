import * as React from "react";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {OrganizationItem} from "./organization-item";
import "./navigator-item.scss";

interface IProps {
    item: AclOrganization;
}

export class NavigatorItem extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element | null {
        if (this.props.item instanceof AclOrganization) {
            return <OrganizationItem organization={this.props.item} />;
        }
        return null;
    }

}