export let config;
config = {
  wsApiPath: '/ws/v1/connect'
};

if (process.env.NODE_ENV === 'development') {
  config['backendPort'] = 8080;
}