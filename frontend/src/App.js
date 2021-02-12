import React, {useEffect, useState} from 'react';
import {bounce, fadeInLeft} from 'react-animations';
import Radium from 'radium';

import useWebSocket, {ReadyState} from 'react-use-websocket';
import {api} from "./http/API";

import {
    Badge,
    Box,
    Button,
    Container,
    Divider,
    Flex,
    Heading,
    IconButton,
    Link,
    LinkBox,
    Spacer,
    Stack,
    Text,
    Tooltip
} from "@chakra-ui/react";

import {useColorMode, useColorModeValue} from "@chakra-ui/color-mode";
import {MoonIcon, SunIcon} from "@chakra-ui/icons";
import {useBreakpointValue} from "@chakra-ui/media-query";

import moment from "moment";
import 'moment/locale/ru';
import {wsAddress} from "./config";

moment.locale('ru');


export const App = () => {
    const [socketUrl, setSocketUrl] = useState(wsAddress + "/news");
    const [messages, setMessages] = useState([]);

    const {
        readyState,
    } = useWebSocket(socketUrl, {
        onOpen: () => console.log('connection opened'),
        shouldReconnect: (closeEvent) => true,
        onMessage: (m) => {
            messages.unshift(JSON.parse(m.data));
            setMessages(messages);
        }
    });

    const connectionStatus = {
        [ReadyState.CONNECTING]: 'устанавливается...',
        [ReadyState.OPEN]: 'открыто',
        [ReadyState.CLOSING]: 'закрывается...',
        [ReadyState.CLOSED]: 'закрыто',
        [ReadyState.UNINSTANTIATED]: 'не поддерживается',
    }[readyState];

    let [sources, setSources] = React.useState([]);
    let [errorText, setErrorText] = React.useState("");

    useEffect(() => {
        async function fetchSources() {
            return await api.GetAllSources();
        }

        fetchSources().then((s) => {
            setSources(s);
        }).catch((e) => {
            setErrorText(e.toString());
        })
    }, []);

    const {colorMode, toggleColorMode} = useColorMode();

    const adaptiveDirection = useBreakpointValue({base: "column", sm: "row"});
    const adaptiveAlign = useBreakpointValue({base: "center", sm: "stretch"});

    const logoColor = useColorModeValue("gray.600", "gray.400");

    return (
        <Container maxW="6xl" p={4}>
            <Flex direction={"row"} align={"center"}>
                <Stack direction={"row"} m={4} align={adaptiveAlign} style={{cursor: "pointer"}}>
                    <Heading style={styles.bounce} size="xl">EDGE</Heading>
                    <Heading size="xl"> | </Heading>
                    <Heading color={logoColor} size="xl">News</Heading>
                </Stack>
                <Divider orientation="vertical"/>
                <Stack direction={adaptiveDirection} spacing={4} align={adaptiveAlign}>
                    <Button colorScheme="teal" variant="link" isActive={true}>
                        Новости
                    </Button>
                    <Button colorScheme="teal" variant="link">
                        Онлайны
                    </Button>
                </Stack>
                <Spacer/>
                <Stack direction={"row"} m={3}>
                    <IconButton colorScheme="gray" onClick={toggleColorMode} aria-label="Switch-Theme"
                                icon={colorMode === "light" ? <MoonIcon/> : <SunIcon/>}/>
                    <Button colorScheme="teal">?</Button>
                </Stack>
            </Flex>
            <Divider mb={5}/>
            <Flex m={4} direction={adaptiveDirection}>
                <Heading>Агрегатор новостей</Heading>
                <Spacer/>
                <Text p={3} color={"gray.300"}>
                    Статус соединения: {connectionStatus}
                </Text>
            </Flex>
            {messages.map((m, i) => {
                return <LinkBox style={styles.fadeInLeft} key={i} as="article" p="5" borderWidth="1px" rounded="md">
                    <Box as="time" dateTime={m["CreatedAt"]}>
                        {moment(m["CreatedAt"]).format("LLL")}
                    </Box>
                    {
                        sources.length > 0 ?
                            <Tooltip label={
                                `Агрегатор: ${sources[m["SourceID"] - 1]["Name"]} ` +
                                `(${sources[m["SourceID"] - 1]["ScrapperType"]["Name"]})`}>
                                <Badge variant="subtle" color={sources[m["SourceID"] - 1]["Color"]} ml={2}>
                                    {m["Rights"]}
                                </Badge>
                            </Tooltip>
                            : ''
                    }
                    <Heading size="md" my="2">
                        <Link href={m["URL"]} external={true}>
                            {m["Title"]}
                        </Link>
                    </Heading>
                    <Text mb={2}>
                        {m["Description"]}
                    </Text>
                    <Flex direction={adaptiveDirection}>
                        <Text as="i" fontSize="xs">{m["Date"]}</Text>
                        <Spacer/>
                        <Text as="i" fontSize="xs" color={"gray.300"}>
                            {m["Tag"] ? `Тег: ${m["Tag"]}` : `Без тега`}
                        </Text>
                    </Flex>
                </LinkBox>
            })}
        </Container>
    );
};

const styles = {
    bounce: {
        animation: 'x 1s',
        animationName: Radium.keyframes(bounce, 'bounce')
    },
    fadeInLeft: {
        animation: 'x 1s',
        animationName: Radium.keyframes(fadeInLeft, 'fadeInLeft'),
    },
}
