import React, {useEffect, useState} from 'react';
import useWebSocket, {ReadyState} from 'react-use-websocket';

import {
    Alert,
    AlertIcon,
    Badge,
    Box,
    Button,
    Checkbox,
    Flex,
    Heading,
    Link,
    LinkBox,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Spacer,
    Stack,
    Text,
    Tooltip,
    useDisclosure,
} from "@chakra-ui/react";

import {useBreakpointValue} from "@chakra-ui/media-query";

import moment from "moment";
import {wsAddress} from "../config";
import {api} from "../http/API";
import {sortArrayByKey} from "../utils";


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
        },
        filter: (m) => selectedStreams.includes((JSON.parse(m.data)).TextStreamID),
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
    let [selectedStreams, setSelectedStreams] = useState([]);
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
            setErrorText("Произошла ошибка, обновите страницу или попробуйте снова позже");
        });

        fetchActiveStreams().then((ts) => {
            console.log(ts);
            setTextStreams(ts);
        }).catch((e) => {
            console.log(e);
            setErrorText("Произошла ошибка, обновите страницу или попробуйте снова позже");
        });
    }, []);

    const preloadSelectedStreams = async () => {
        setMessages([]);
        setErrorText('');

        try {
            let m = [];
            for (let s of selectedStreams) {
                let data = await api.GetMessagesByTextStreamID(s);
                m = [...m, ...data];
            }
            m = sortArrayByKey(m, "CreatedAt", false);
            setMessages(m);
        } catch (error) {
            setErrorText("Произошла ошибка, попробуйте снова позже");
        }
    }

    let dividedStreams = [];

    const adaptiveDirection = useBreakpointValue({base: "column", sm: "row"});
    const adaptiveMargin = useBreakpointValue({base: 0, md: 2, lg: 4, xl: 6});

    const {isOpen, onOpen, onClose} = useDisclosure();

    return (
        <Box>
            <Flex m={4} direction={adaptiveDirection}>
                <Heading>Агрегатор онлайнов</Heading>
                <Button ml={adaptiveMargin} onClick={onOpen}>Фильтрация</Button>
                <Spacer/>
                <Text p={3} color={"gray.300"}>
                    Статус соединения: {connectionStatus}
                </Text>
            </Flex>

            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>Активные трансляции</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Stack pl={2} mt={1} spacing={1}>
                            {
                                sources ? sortArrayByKey(textStreams, "SourceID").map((ts, i) => {
                                    const checkbox = <Checkbox
                                        key={i}
                                        isChecked={selectedStreams.indexOf(ts.ID) !== -1}
                                        onChange={(e) => {
                                            let ss = [...selectedStreams];

                                            if (e.target.checked) {
                                                ss.push(ts.ID);
                                                setSelectedStreams(ss);
                                            } else {
                                                ss = ss.filter(item => item !== ts.ID);
                                                setSelectedStreams(ss);
                                            }
                                        }}
                                    >
                                        {ts.Name} &nbsp;
                                        <Badge>{
                                            ts.LastStreamUpdate ? moment(ts.LastStreamUpdate).format("hh:mm") : ''
                                        }</Badge>
                                    </Checkbox>
                                    if (dividedStreams.find((elem) => elem === ts.SourceID)) {
                                        return checkbox;
                                    } else {
                                        dividedStreams.push(ts.SourceID);
                                        const heading = <Heading p={2} as="h5" size="sm">
                                            {sources.find((elem) => elem.ID === ts.SourceID).Name}
                                        </Heading>
                                        return <div>
                                            {heading}
                                            {checkbox}
                                        </div>
                                    }
                                }) : ''
                            }
                        </Stack>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={() => {
                            onClose();
                            preloadSelectedStreams()
                        }}>
                            Применить
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>

            {errorText ? <Alert status="error"><AlertIcon/>{errorText}</Alert> : ''}

            {sources ? messages.map((m, i) => {
                let body = m?.Body.trim().split('\n');
                return <LinkBox key={i} as="article" p="5" borderWidth="1px" rounded="md">
                    <Box as="time" dateTime={m["CreatedAt"]}>
                        {moment(m["CreatedAt"]).format("LLL")}
                    </Box>
                    {
                        sources.length > 0 && textStreams.length > 0 ?
                            <Tooltip label={
                                `Источник: ${textStreams?.find(ts => ts.ID === m["TextStreamID"])?.Source.Name} ` +
                                `(${textStreams.find(ts => ts.ID === m["TextStreamID"])?.Source.URL})`}>
                                <Link href={textStreams?.find(ts => ts.ID === m["TextStreamID"])?.URL}>
                                    <Badge variant="subtle"
                                           color={textStreams?.find(ts => ts.ID === m["TextStreamID"])?.Source.Color}
                                           ml={2}>
                                        {textStreams?.find(ts => ts.ID === m["TextStreamID"])?.Name}
                                    </Badge>
                                </Link>
                            </Tooltip>
                            : ''
                    }
                    <Heading size="md" my="2">
                        {m["Title"]}
                    </Heading>
                    <Text mb={2}>
                        {body.map((text, index) => {
                            if (index === body.length - 1) return <span key={index}>{text}</span>
                            return <span key={index}>{text}<br/><br/></span>
                        })}
                    </Text>
                    <Flex direction={adaptiveDirection}>
                        <Text as="i" fontSize="xs">{m["Time"]}</Text>
                    </Flex>
                </LinkBox>
            }) : ''}
        </Box>
    );
};
