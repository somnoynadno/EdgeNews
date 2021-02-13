import React, {useEffect, useState} from 'react';
import useWebSocket, {ReadyState} from 'react-use-websocket';

import {Badge, Box, Flex, Heading, Link, LinkBox, Spacer, Text, Tooltip} from "@chakra-ui/react";

import {useBreakpointValue} from "@chakra-ui/media-query";

import moment from "moment";
import {wsAddress} from "../config";
import {api} from "../http/API";


export const TextStreamsPage = () => {
    const [socketUrl, setSocketUrl] = useState(wsAddress + "/streams");
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
    let [textStreams, setTextStreams] = React.useState([]);
    let [errorText, setErrorText] = React.useState("");

    useEffect(() => {
        async function fetchSources() {
            return await api.GetAllSources();
        }

        async function fetchActiveStreams() {
            return await api.GetActiveTextStreams();
        }

        fetchSources().then((s) => {
            console.log(s);
            setSources(s);
        }).catch((e) => {
            console.log(e);
            setErrorText("Произошла ошибка, обновите страницу");
        });

        fetchActiveStreams().then((ts) => {
            console.log(ts);
            setTextStreams(ts);
        }).catch((e) => {
            console.log(e);
            setErrorText("Произошла ошибка, обновите страницу");
        });
    }, []);

    const adaptiveDirection = useBreakpointValue({base: "column", sm: "row"});

    return (
        <Box>
            <Flex m={4} direction={adaptiveDirection}>
                <Heading>Агрегатор онлайнов</Heading>
                <Spacer/>
                <Text p={3} color={"gray.300"}>
                    Статус соединения: {connectionStatus}
                </Text>
            </Flex>
            {messages.map((m, i) => {
                return <LinkBox key={i} as="article" p="5" borderWidth="1px" rounded="md">
                    <Box as="time" dateTime={m["CreatedAt"]}>
                        {moment(m["CreatedAt"]).format("LLL")}
                    </Box>
                    {
                        sources.length > 0 && textStreams.length > 0 ?
                            <Tooltip label={
                                `Источник: ${textStreams.find(ts => ts.ID === m["TextStreamID"])?.Source.Name} ` +
                                `(${textStreams.find(ts => ts.ID === m["TextStreamID"])?.Source.URL})`}>
                                <Link href={textStreams.find(ts => ts.ID === m["TextStreamID"])?.URL}>
                                <Badge variant="subtle" color={textStreams.find(ts => ts.ID === m["TextStreamID"])?.Source.Color} ml={2}>
                                    {textStreams.find(ts => ts.ID === m["TextStreamID"])?.Name}
                                </Badge>
                                </Link>
                            </Tooltip>
                            : ''
                    }
                    <Heading size="md" my="2">
                        {m["Title"]}
                    </Heading>
                    <Text mb={2}>
                        {m["Body"]}
                    </Text>
                    <Flex direction={adaptiveDirection}>
                        <Text as="i" fontSize="xs">{m["Time"]}</Text>
                    </Flex>
                </LinkBox>
            })}
        </Box>
    );
};
