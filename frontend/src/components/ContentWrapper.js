import React from 'react';
import {Route} from "react-router-dom";

import {
    Box,
    Button,
    Container,
    Divider,
    Flex,
    Heading,
    IconButton,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Spacer,
    Link,
    Stack,
    Text,
    useDisclosure
} from "@chakra-ui/react";

import {useColorMode, useColorModeValue} from "@chakra-ui/color-mode";
import {useBreakpointValue} from "@chakra-ui/media-query";
import {MoonIcon, QuestionIcon, QuestionOutlineIcon, SunIcon} from "@chakra-ui/icons";

import {TextStreamsPage} from "../pages/TextStreamsPage";
import {NewsPage} from "../pages/NewsPage";

import history from "../history";


export const ContentWrapper = () => {
    const {colorMode, toggleColorMode} = useColorMode();

    const adaptiveDirection = useBreakpointValue({base: "column", sm: "row"});
    const adaptiveAlign = useBreakpointValue({base: "center", sm: "stretch"});

    const logoColor = useColorModeValue("gray.600", "gray.400");

    const {isOpen, onOpen, onClose} = useDisclosure();

    if (window.location.pathname === '/') history.push('/news');
    else return (<Container maxW="6xl" p={4}>
            <Flex direction={"row"} align={"center"}>
                <Stack onClick={() => window.location.reload()} direction={"row"}
                       m={4} align={adaptiveAlign} style={{cursor: "pointer"}}>
                    <Heading size="xl">EDGE</Heading>
                    <Heading size="xl"> | </Heading>
                    <Heading color={logoColor} size="xl">News</Heading>
                </Stack>
                <Divider orientation="vertical"/>
                <Stack direction={adaptiveDirection} spacing={4} align={adaptiveAlign}>
                    <Button colorScheme="teal" onClick={() => history.push('/news')}
                            variant="link" isActive={window.location.pathname === '/news'}>
                        Новости
                    </Button>
                    <Button colorScheme="teal" onClick={() => history.push('/streams')}
                            variant="link" isActive={window.location.pathname === '/streams'}>
                        Онлайны
                    </Button>
                </Stack>
                <Spacer/>
                <Stack direction={"row"} m={3}>
                    <IconButton colorScheme="gray" onClick={toggleColorMode} aria-label="Switch-Theme"
                                icon={colorMode === "light" ? <MoonIcon/> : <SunIcon/>}/>
                    <IconButton onClick={onOpen} colorScheme="gray"
                                icon={colorMode === "light" ? <QuestionOutlineIcon/> : <QuestionIcon/>}/>
                </Stack>
            </Flex>
            <Divider mb={5}/>
            <Box>
                <Route exact path="/news" component={NewsPage}/>
                <Route exact path="/streams" component={TextStreamsPage}/>
            </Box>

            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>О сервисе</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Text>
                            <strong>EDGE News</strong> представляет собой агрегатор последних новостных сводок
                            и текстовых онлайн трансляций с самых популярных новостных ресурсов
                            в режиме реального времени.
                        </Text>
                        <br/>
                        <Text>
                            Для мониторинга интересных вам событий воспользуйтесь кнопной фильтрации в
                            разделе "Агрегатор Онлайнов".
                        </Text>
                        <Divider m={4}/>
                        <Text fontSize="sm">Разработчик: <Link href={"https://github.com/somnoynadno"}>@somnoynadno</Link></Text>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={onClose}>
                            Понятно
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Container>
    );
};
