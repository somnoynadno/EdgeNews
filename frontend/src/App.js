import React from 'react';
import {
    Route,
    Switch,
    withRouter
} from "react-router-dom";

import {ContentWrapper} from "./components/ContentWrapper";


const App = () => {
    return (
        <Switch>
            <Route path='/' component={ContentWrapper} />
        </Switch>
    );
};

export default withRouter(App);
