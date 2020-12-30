export let config;

if (process.env.NODE_ENV === 'development') {
  config = {
    backendUrl: 'ws://localhost:8080/ws/v1/connect'
  };
} else {
  config = {
    backendUrl: 'ws://happiness-door-bot-dev.herokuapp.com/ws/v1/connect'
  };
}