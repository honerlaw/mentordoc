import * as React from "react";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {OrganizationItemView} from "./organization-item-view";
import {AclFolder} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/folder/acl-folder";
import {FolderItemView} from "./folder-item-view";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {DocumentItemView} from "./document-item-view";

interface IProps {
    item: AclOrganization | AclFolder | AclDocument;
}

export class NavigatorItem extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element | null {
        if (this.props.item instanceof AclOrganization) {
            return <OrganizationItemView organization={this.props.item} />;
        }
        if (this.props.item instanceof AclFolder) {
            return <FolderItemView folder={this.props.item} />;
        }
        return <DocumentItemView document={this.props.item} />;
    }

}