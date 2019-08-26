import * as ReactDOM from "react-dom";
import * as React from "react";
import {Main} from "./component/main";
import {Provider} from "react-redux";
import {RootStore} from "@honerlawd/mentordoc-frontend-shared/dist/store/store";
import {BrowserRouter} from "react-router-dom";

const container: HTMLElement = document.createElement("div");
container.id = "mentordoc-mount";
document.getElementsByTagName("body")[0].append(container);

ReactDOM.render(<Provider store={RootStore}>
    <BrowserRouter>
        <Main/>
    </BrowserRouter>
</Provider>, container);