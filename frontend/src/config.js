export const apiAddress = (process.env.NODE_ENV === "production" ? 'http://edge.somnoynadno.ru/api' : 'http://localhost:8080/api');
export const wsAddress = (process.env.NODE_ENV === "production" ? 'ws://edge.somnoynadno.ru/ws' : 'ws://localhost:8080/ws');
