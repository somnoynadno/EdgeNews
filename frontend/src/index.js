import React from 'react';
import ReactDOM from 'react-dom';
import {Router} from "react-router-dom";

import App from './App';

import {ChakraProvider} from "@chakra-ui/react";
import {ColorModeScript} from "@chakra-ui/color-mode";

import theme from "./theme";
import history from './history';

import moment from "moment";
import 'moment/locale/ru';

moment.locale('ru');


ReactDOM.render(
    <React.StrictMode>
        <ChakraProvider>
            <ColorModeScript initialColorMode={theme.config.initialColorMode}/>
            <Router history={history}>
                <App/>
            </Router>
        </ChakraProvider>
    </React.StrictMode>,
    document.getElementById('root')
);
