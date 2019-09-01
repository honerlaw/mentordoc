import {AclDocument} from "./acl-document";

export interface IDocumentState {
    documentMap: Record<string, AclDocument[]>;
}

export const INITIAL_DOCUMENT_STATE: IDocumentState = {
    documentMap: {}
};
