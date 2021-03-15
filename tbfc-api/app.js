var express = require('express');
var app = express();
const cors = require("cors");
app.options("*",cors());
app.use(cors());

var swaggerUi = require('swagger-ui-express');
var swaggerDocument = require('./swagger.json');


var TFBCController = require('./TFBCController');

app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocument));
app.use('/tfbc', TFBCController);

module.exports = app;
