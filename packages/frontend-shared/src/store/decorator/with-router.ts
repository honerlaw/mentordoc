import {withRouter} from "react-router";

export function WithRouter(): ClassDecorator {
    return withRouter as any;
}