import * as ReactDOM from "react-dom";
import * as React from "react";
import {Main} from "./component/main";
import {Provider} from "react-redux";
import {BrowserRouter} from "react-router-dom";
import {RootStore} from "@honerlawd/mentordoc-frontend-shared/dist/store/store";

const container: HTMLElement = document.createElement("div");
container.id = "mentordoc-mount";
document.getElementsByTagName("body")[0].append(container);

ReactDOM.render(<Provider store={RootStore}>
    <BrowserRouter>
        <Main/>
    </BrowserRouter>
</Provider>, container);