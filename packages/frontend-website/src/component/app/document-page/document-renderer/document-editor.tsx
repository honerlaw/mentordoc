import * as React from "react";
import * as tuiEditor from "tui-editor";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import "./document-editor.scss";
import {debounce} from "lodash";
import {
    CombineDispatchers,
    ConnectProps,
    IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    IUpdateDocumentDispatch,
    UpdateDocument
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/update-document";

interface IProps extends Partial<IDispatchPropMap<IUpdateDocumentDispatch>> {
    document: AclDocument;
}

interface IState {
    name: string;
    content: string;
}

@ConnectProps(null, CombineDispatchers(UpdateDocument.dispatch))
export class DocumentEditor extends React.PureComponent<IProps, IState> {

    private editorRef: React.RefObject<any>;
    private editor: tuiEditor.default;

    public constructor(props: IProps) {
        super(props);

        this.state = {
            name: props.document.model.drafts[0].name,
            content: props.document.model.drafts[0].content!.content
        };

        this.editorRef = React.createRef();
        this.onNameChange = this.onNameChange.bind(this);

        // debounce both of these for perf reasons, fetching markdown is expensive, saving is expensive
        this.onContentChange = debounce(this.onContentChange.bind(this), 500);
        this.save = debounce(this.save.bind(this), 500);
    }

    public componentDidMount(): void {
        // @todo sort out the typings
        const Editor: any = tuiEditor;

        this.editor = new Editor({
            el: this.editorRef.current,
            initialValue: this.state.content,
            usageStatistics: false,
            events: {
                change: this.onContentChange
            },
            initialEditType: 'wysiwyg',
            previewStyle: 'vertical',
            height: 'auto'
        });
    }

    public render(): JSX.Element | null {
        return <div className={"document-editor"}>
            <input type={"text"} value={this.state.name} onChange={this.onNameChange}/>
            <div ref={this.editorRef}/>
        </div>;
    }

    private onNameChange(event: React.ChangeEvent<HTMLInputElement>): void {
        const name: string = event.target.value;
        this.setState({name});

        this.save(name, this.editor.getMarkdown());
    }

    private onContentChange(): void {
        const content: string = this.editor.getMarkdown();
        this.setState({content});

        this.save(this.state.name, content);
    }

    private save(name: string, content: string): void {
        this.props.dispatch!.updateDocument({
            documentId: this.props.document.model.id,
            draftId: this.props.document.model.drafts[0].id,
            name,
            content
        });
    }

}