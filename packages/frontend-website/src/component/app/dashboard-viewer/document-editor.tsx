import * as React from "react";
import * as tuiEditor from "tui-editor";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import "./document-editor.scss";
import {onChangeSetState} from "../../../util";
import {debounce} from "lodash";

interface IProps {
    document: AclDocument;
}

interface IState {
    name: string;
    content: string;
}

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
        this.onContentChange = debounce(this.onContentChange.bind(this), 400);
    }

    public componentDidMount(): void {
        // @todo sort out the typings
        const Editor: any = tuiEditor;

        this.editor = new Editor({
            el: this.editorRef.current,
            initialValue: this.state.content,
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
            <input type={"text"} value={this.state.name} onChange={onChangeSetState<IState>("name", this)}/>
            <div ref={this.editorRef}/>
        </div>;
    }

    private onContentChange(): void {
        this.setState({
            content: this.editor.getMarkdown()
        });
    }

}