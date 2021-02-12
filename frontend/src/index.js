import React from 'react';
import ReactDOM from 'react-dom';
import {App} from './App';

import {ChakraProvider} from "@chakra-ui/react";
import {ColorModeScript} from "@chakra-ui/color-mode";

import theme from "./theme";

ReactDOM.render(
    <React.StrictMode>
        <ChakraProvider>
            <ColorModeScript initialColorMode={theme.config.initialColorMode} />
            <App/>
        </ChakraProvider>
    </React.StrictMode>,
    document.getElementById('root')
);
