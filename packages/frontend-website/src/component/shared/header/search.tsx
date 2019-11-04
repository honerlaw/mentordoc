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
import {Link} from "react-router-dom";
import {DocumentDraft} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/document-draft";
import {Document} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/document";
import "./search.scss";
import * as moment from "moment";

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
            <input type={"text"} placeholder={"search"} onChange={this.onChange}/>
            {this.renderSearchResults()}
        </div>;
    }

    private renderSearchResults(): JSX.Element | null {
        const docs: AclDocument[] | null = this.props.selector!.searchDocuments;
        if (!docs) {
            return null;
        }

        return <ul className={"search-results"}>
            {docs.map((aclDoc: AclDocument): JSX.Element => {
                const doc: Document = aclDoc.model;
                const draft: DocumentDraft = aclDoc.model.drafts[0];
                const timeago: string = moment((draft.publishedAt || draft.updatedAt) / 1e+6).fromNow();

                return <li>
                    <Link to={`/app/${doc.organizationId}/${doc.id}`}>
                        <div className={"title"}>{draft.name}</div>
                        <div className={"timeago"}>last updated {timeago}</div>
                    </Link>
                </li>
            })}
        </ul>;
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