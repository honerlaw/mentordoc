import * as React from "react";
import {debounce} from "lodash";
import {
    CombineDispatchers, CombineSelectors,
    ConnectProps,
    IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    ISearchDocumentsDispatch, SearchDocuments
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/search-documents";
import {
    ISetSearchDocumentsSelector, SetSearchDocuments
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-search-documents";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";

interface IProps extends Partial<IDispatchPropMap<ISearchDocumentsDispatch> & ISelectorPropMap<ISetSearchDocumentsSelector>> {

}

interface IState {
    searchQuery: string;
}

@ConnectProps(CombineSelectors(SetSearchDocuments.selector), CombineDispatchers(SearchDocuments.dispatch))
export class Search extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            searchQuery: ""
        };

        this.onChange = this.onChange.bind(this);
        // debounce this instead of onChange, otherwise, we cant get the correct input value from the event
        this.handleOnChange = debounce(this.handleOnChange.bind(this), 250);
    }

    public render(): JSX.Element {
        return <div className={"search"}>
            <input type={"text"} placeholder={"search"} onChange={this.onChange} />
            <div className={"search-results"}>
                {this.renderSearchResults()}
            </div>
        </div>;
    }

    private renderSearchResults(): JSX.Element[] {
        const docs: AclDocument[] | null = this.props.selector!.searchDocuments;
        if (!docs) {
            return [];
        }

        return docs.map((doc: AclDocument): JSX.Element => {
            return <span>{doc.model.drafts[0].name}</span>
        });
    }

    private async onChange(event: React.ChangeEvent<HTMLInputElement>): Promise<void> {
        await this.handleOnChange(event.target.value);
    }

    private async handleOnChange(searchQuery: string): Promise<void> {
        this.setState({
            searchQuery
        });

        await this.props.dispatch!.searchDocuments({
            searchQuery
        });
    }

}