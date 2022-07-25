const express = require('express');

const router = express.Router();
const signIn = require('../modules/signIn');
const common = require('../modules/common');

/* self book borow logs */
router.get('/:libraryId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const log = {
    libraryId: Number(req.params.libraryId),
    userId
  };
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      const sqlGetLogs = `select book.title, log.* from log, book WHERE log.libraryId = ${log.libraryId} and log.userId = '${log.userId}' and log.bookId = book.id order by log.dateUpdated desc;`;
      connection.query(sqlGetLogs, (error, results) => {
        connection.release();
        if (!error) {
          res.send({
            code: 2000,
            data: results
          });
        }
        else {
          res.statusCode = 500;
          res.send(common.database500(error));
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(err));
    }
  });
});

/* self book borow logs */
router.get('/:libraryId/:bookId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const log = {
    libraryId: Number(req.params.libraryId),
    bookId: Number(req.params.bookId),
    userId
  };
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      const sqlGetLogs = `select book.title, log.* from log, book WHERE log.libraryId = ${log.libraryId} and log.userId = '${log.userId}' and bookId = ${log.bookId} and log.bookId = book.id order by log.dateUpdated desc;`;
      connection.query(sqlGetLogs, (error, results) => {
        connection.release();
        if (!error) {
          res.send({
            code: 2000,
            data: results
          });
        }
        else {
          res.statusCode = 500;
          res.send(common.database500(error));
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(err));
    }
  });
});

module.exports = router;
