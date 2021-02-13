import React, {useEffect, useState} from 'react';
import useWebSocket, {ReadyState} from 'react-use-websocket';

import {Alert, AlertIcon, Badge, Box, Flex, Heading, Link, LinkBox, Spacer, Text, Tooltip} from "@chakra-ui/react";
import {useBreakpointValue} from "@chakra-ui/media-query";

import moment from "moment";
import {wsAddress} from "../config";
import {api} from "../http/API";


export const NewsPage = () => {
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

        async function fetchLastNews() {
            return await api.GetLastNews(25);
        }

        fetchSources().then((s) => {
            setSources(s);
            fetchLastNews().then((ln) => {
                setMessages(ln);
            }).catch((e) => {
                console.log(e);
                setErrorText("Произошла ошибка, попробуйте обновить страницу или зайти позже");
            });
        }).catch((e) => {
            console.log(e);
            setErrorText("Произошла ошибка, попробуйте обновить страницу или зайти позже");
        });
    }, []);

    const adaptiveDirection = useBreakpointValue({base: "column", sm: "row"});

    return (
        <Box>
            <Flex m={4} direction={adaptiveDirection}>
                <Heading>Агрегатор новостей</Heading>
                <Spacer/>
                <Text p={3} color={"gray.300"}>
                    Статус соединения: {connectionStatus}
                </Text>
            </Flex>

            {errorText ? <Alert status="error"><AlertIcon />{errorText}</Alert> : ''}

            {
                messages.map((m, i) => {
                    return <LinkBox key={i} as="article" p="5" borderWidth="1px" rounded="md">
                        <Box as="time" dateTime={m["CreatedAt"]}>
                            {moment(m["CreatedAt"]).format("LLL")}
                        </Box>
                        {
                            sources.length > 0 ?
                                <Tooltip label={
                                    `Агрегатор: ${sources[m["SourceID"] - 1]["Name"]} ` +
                                    `(${sources[m["SourceID"] - 1]["ScrapperType"]["Name"]})`}>
                                    <Badge variant="outline" color={sources[m["SourceID"] - 1]["Color"]} ml={2}>
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
                                {m["Tag"] ? `Tag: ${m["Tag"]}` : `Без тега`}
                            </Text>
                        </Flex>
                    </LinkBox>
                })
            }
        </Box>
    );
};
