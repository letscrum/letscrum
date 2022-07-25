var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');

const swaggerUi = require('swagger-ui-express');
const swaggerDocument = require('./swagger-output.json');
const DebugControl = require('./utilities/debug');
var sqlPool = require('./modules/sqlPool');

var app = express();

app.use((_req, res, next) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
  res.setHeader('Access-Control-Allow-Headers', '*');
  res.setHeader('Access-Control-Allow-Credentials', true);
  next();
});

const indexRouter = require('./routes/index');
const usersRouter = require('./routes/users');
const librariesRouter = require('./routes/libraries');
const membersRouter = require('./routes/members');
const booksRouter = require('./routes/books');
const logsRouter = require('./routes/logs');

app.disable('etag');

app.use(
  '/swagger',
  swaggerUi.serve,
  swaggerUi.setup(swaggerDocument)
);

// view engine setup
// app.set('views', path.join(__dirname, 'views'));
// app.set('view engine', 'pug');

global.sqlPool = sqlPool;

// JWT configuration
global.jwtSecret = process.env.JWT_SECRET;
global.jwtExpiresIn = 30 * 24 * 3600;
// global.jwtExpiresIn = 60;

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));
app.use(DebugControl.log.request());

app.use('/v1/', indexRouter);
app.use('/v1/users', usersRouter);
app.use('/v1/libraries', librariesRouter);
app.use('/v1/members', membersRouter);
app.use('/v1/books', booksRouter);
app.use('/v1/logs', logsRouter);

// catch 404 and forward to error handler
app.use((_req, _res, next) => {
  next(createError(404));
});

// error handler
app.use((err, req, res) => {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

module.exports = app;
