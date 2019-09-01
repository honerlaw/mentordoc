import * as React from "react";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {NavigatorItemView} from "./navigator-item-view";
import {IDropdownButtonOption} from "../../shared/dropdown-button";

interface IProps {
    document: AclDocument;
}

export class DocumentItemView extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <NavigatorItemView title={this.props.document.model.name}
                                  hasChildren={false}
                                  isExpanded={false}
        options={this.getOptions()}/>;
    }

    private getOptions(): IDropdownButtonOption[] {
        const options: IDropdownButtonOption[] = [];

        if (this.props.document.hasAction("delete")) {
            options.push({
                label: "delete",
                onClick: () => console.log("delete")
            });
        }

        return options;
    }

}
