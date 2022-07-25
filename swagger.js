const swaggerAutogen = require('swagger-autogen')();

const doc = {
  info: {
    title: 'LetScrum API',
    description: 'Description',
    version: '1.0.0'
  },
  host: 'localhost:3999',
  schemes: ['http']
};

const outputFile = './swagger-output.json';
const endpointsFiles = ['app.js'];

/* NOTE: if you use the express Router, you must pass in the
   'endpointsFiles' only the root file where the route starts,
   such as index.js, app.js, routes.js, ... */

swaggerAutogen(outputFile, endpointsFiles, doc);
