const express = require('express');

const router = express.Router();
const request = require('request');
const signIn = require('../modules/signIn');
const common = require('../modules/common');

/* create book to a library
  libraryId, openId, isbn
*/
router.post('/', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.isbn) {
    const book = {
      libraryId: Number(req.body.libraryId),
      userId,
      isbn: req.body.isbn
    };
    // let categoryId = isNaN(req.body.categoryId) ? parseInt(req.body.categoryId) : 0;
    // const categoryId = 0;
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId} and isAdmin = 1;`;
        const sqlGetBook = `select * from book where isbn = '${book.isbn}' and libraryId = ${book.libraryId};`;
        connectionGetData.query(
          sqlGetAdmin + sqlGetBook,
          (errorGetData, resultsGetLibraryAndBook) => {
            connectionGetData.release();
            if (!errorGetData) {
              const resultsGetLibrary = resultsGetLibraryAndBook[0];
              const resultsGetBook = resultsGetLibraryAndBook[1];
              if (resultsGetLibrary.length > 0) {
                if (resultsGetBook.length > 0) {
                  res.send({
                    code: 2000,
                    data: {
                      id: resultsGetBook[0].id,
                      userId: resultsGetBook[0].userId,
                      libraryId: resultsGetBook[0].libraryId,
                      title: resultsGetBook[0].title,
                      isbn: resultsGetBook[0].isbn,
                      coverUrl: resultsGetBook[0].coverUrl,
                      withUserId: resultsGetBook[0].withUserId,
                      isDeleted: resultsGetBook[0].isDeleted,
                      datecreated: resultsGetBook[0].datecreated,
                      dateupdated: resultsGetBook[0].dateupdated
                    }
                  });
                }
                else {
                  const doubanGetBoookUrl = common.doubanDataUrl;
                  const options = {
                    url: doubanGetBoookUrl + book.isbn,
                    method: 'GET',
                    headers: {
                      apikey: common.doubanDataAPIKey
                    }
                  };
                  request(options, (errorRequest, response, body) => {
                    if (!errorRequest && response && response.statusCode < 400) {
                      const data = JSON.parse(body);
                      const { title } = data.items[0].volumeInfo;
                      const coverUrl = ''; // data.cover_url;
                      global.sqlPool.getConnection((errorConnectionAddBook, connectionAddBook) => {
                        if (!errorConnectionAddBook) {
                          const sqlAddBook = `insert into book(libraryId, userId, isbn, title, coverUrl, withUserId) values( ${book.libraryId}, '${userId}','${book.isbn}' ,'${title}' ,'${coverUrl}', '');`;
                          connectionAddBook.query(sqlAddBook, (errorAddBook, resultsAddBook) => {
                            connectionAddBook.release();
                            if (!errorAddBook) {
                              if (resultsAddBook.affectedRows > 0) {
                                res.send({
                                  code: 2000,
                                  data: {
                                    id: resultsAddBook.insertId,
                                    userId: book.userId,
                                    libraryId: book.libraryId,
                                    title,
                                    isbn: book.isbn,
                                    coverurl: coverUrl,
                                    withUserId: '',
                                    dateCreated: (new Date()).toUTCString(),
                                    dateUpdated: (new Date()).toUTCString()
                                  }
                                });
                              }
                              else {
                                res.statusCode = 500;
                                res.send(common.database500());
                              }
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.database500(errorAddBook));
                            }
                          });
                        }
                        else {
                          res.statusCode = 500;
                          res.send(common.database500(errorConnectionAddBook));
                        }
                      });
                    }
                    else {
                      // eslint-disable-next-line no-console
                      console.error(errorRequest);
                      res.statusCode = 500;
                      res.send(common.customMessage500("Can't get book from 3rd party API."));
                    }
                  });
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No permission.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.database500(errorGetData));
            }
          }
        );
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.put('/take', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.bookId) {
    const book = {
      id: req.body.bookId,
      libraryId: req.body.libraryId
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetMember = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId};`;
        const sqlGetBook = `select * from book WHERE libraryId = ${book.libraryId} and id = ${book.id};`;
        connectionGetData.query(sqlGetMember + sqlGetBook, (errorGetData, resultsGetBook) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultsGetBook[0].length > 0) {
              if (resultsGetBook[1].length > 0) {
                const getBook = resultsGetBook[1][0];
                if (getBook.withUserId === null || getBook.withUserId === '') {
                  global.sqlPool.getConnection((errorConnectionAddData, connectionAddData) => {
                    if (!errorConnectionAddData) {
                      const sqlAddLog = `insert into log(libraryId, bookId, userId) values(${book.libraryId}, ${getBook.id}, '${userId}');`;
                      const sqlUpdateBook = `update book set withUserId = '${userId}' where libraryId = ${book.libraryId} and id = ${getBook.id};`;
                      connectionAddData.query(
                        sqlAddLog + sqlUpdateBook,
                        (error, resultsTakeBook) => {
                          connectionAddData.release();
                          if (!error) {
                            if (resultsTakeBook[0].affectedRows > 0) {
                              if (resultsTakeBook[1].affectedRows > 0) {
                                res.send({
                                  code: 2000,
                                  data: {
                                    id: resultsTakeBook[0].insertId,
                                    libraryId: book.libraryId,
                                    bookId: getBook.id,
                                    userId,
                                    isReturned: 0,
                                    dateReturned: null,
                                    dateCreated: (new Date()).toUTCString(),
                                    dateUpdated: (new Date()).toUTCString()
                                  }
                                });
                              }
                              else {
                                res.statusCode = 500;
                                res.send(common.customMessage500('Book update failed.'));
                              }
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.customMessage500('Log add failed.'));
                            }
                          }
                          else {
                            res.statusCode = 500;
                            res.send(common.database500(error));
                          }
                        }
                      );
                    }
                    else {
                      res.statusCode = 500;
                      res.send(common.database500(errorConnectionAddData));
                    }
                  });
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500('With someone now.'));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No book.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('No member.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body problem.'));
  }
});

router.put('/return', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.bookId) {
    const book = {
      id: req.body.bookId,
      libraryId: req.body.libraryId
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetMember = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId};`;
        const sqlGetBook = `select * from book WHERE libraryId = ${book.libraryId} and id = ${book.id};`;
        connectionGetData.query(sqlGetMember + sqlGetBook, (errorGetData, resultsGetBook) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultsGetBook[0].length > 0) {
              if (resultsGetBook[1].length > 0) {
                const getBook = resultsGetBook[1][0];
                if (getBook.withUserId !== null && getBook.withUserId !== '') {
                  global.sqlPool.getConnection(
                    (errorConnectionUpdateData, connectionUpdateData) => {
                      if (!errorConnectionUpdateData) {
                        const sqlUpdateLog = `update log set isReturned = 1, dateReturned = NOW() where libraryId = ${book.libraryId} and bookId = ${getBook.id} and userId = '${userId}';`;
                        const sqlUpdateBook = `update book set withUserId = '' where libraryId = ${book.libraryId} and id = ${getBook.id};`;
                        connectionUpdateData.query(
                          sqlUpdateLog + sqlUpdateBook,
                          (error, resultsReturnBook) => {
                            connectionUpdateData.release();
                            if (!error) {
                              if (resultsReturnBook[0].affectedRows > 0) {
                                if (resultsReturnBook[1].affectedRows > 0) {
                                  res.send({
                                    code: 2000,
                                    data: {
                                      id: 0,
                                      libraryId: book.libraryId,
                                      bookId: getBook.id,
                                      userId,
                                      isReturned: 1,
                                      dateReturned: (new Date()).toUTCString()
                                    }
                                  });
                                }
                                else {
                                  res.statusCode = 500;
                                  res.send(common.customMessage500('Book update failed.'));
                                }
                              }
                              else {
                                res.statusCode = 500;
                                res.send(common.customMessage500('Log add failed.'));
                              }
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.database500(error));
                            }
                          }
                        );
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500(errorConnectionUpdateData));
                      }
                    }
                  );
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500("You don't have it."));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No book.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not a member.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.put('/scan/take', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.isbn) {
    const book = {
      libraryId: Number(req.body.libraryId),
      isbn: req.body.isbn
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetMember = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId};`;
        const sqlGetBook = `select * from book WHERE libraryid = ${book.libraryId} and isbn = '${book.isbn}';`;
        connectionGetData.query(sqlGetMember + sqlGetBook, (errorGetData, resultsGetBook) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultsGetBook[0].length > 0) {
              if (resultsGetBook[1].length > 0) {
                const getBook = resultsGetBook[1][0];
                if (getBook.withUserId === null || getBook.withUserId === '') {
                  global.sqlPool.getConnection((errorConnectionAddData, connectionAddData) => {
                    if (!errorConnectionAddData) {
                      const sqlAddLog = `insert into log(libraryId, bookId, userId) values(${book.libraryId}, ${getBook.id}, '${userId}');`;
                      const sqlUpdateBook = `update book set withUserId = '${userId}' where libraryId = ${book.libraryId} and id = ${getBook.id};`;
                      connectionAddData.query(
                        sqlAddLog + sqlUpdateBook,
                        (error, resultsTakeBook) => {
                          connectionAddData.release();
                          if (!error) {
                            if (resultsTakeBook[0].affectedRows > 0) {
                              if (resultsTakeBook[1].affectedRows > 0) {
                                res.send({
                                  code: 2000,
                                  data: {
                                    id: resultsTakeBook[0].insertId,
                                    libraryId: book.libraryId,
                                    bookId: getBook.id,
                                    userId,
                                    isReturned: 0,
                                    dateReturned: null,
                                    dateCreated: (new Date()).toUTCString(),
                                    dateUpdated: (new Date()).toUTCString()
                                  }
                                });
                              }
                              else {
                                res.statusCode = 500;
                                res.send(common.customMessage500('Book update failed.'));
                              }
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.customMessage500('Log add failed.'));
                            }
                          }
                          else {
                            res.statusCode = 500;
                            res.send(common.database500(error));
                          }
                        }
                      );
                    }
                    else {
                      res.statusCode = 500;
                      res.send(common.database500(errorConnectionAddData));
                    }
                  });
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500('With someone now.'));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No book.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not a member.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.put('/scan/return', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.isbn) {
    const book = {
      libraryId: Number(req.body.libraryId),
      isbn: req.body.isbn
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetMember = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId};`;
        const sqlGetBook = `select * from book WHERE libraryId = ${book.libraryId} and isbn = '${book.isbn}';`;
        connectionGetData.query(sqlGetMember + sqlGetBook, (errorGetData, resultsGetBook) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultsGetBook[0].length > 0) {
              if (resultsGetBook[1].length > 0) {
                const getBook = resultsGetBook[1][0];
                if (getBook.withUserId !== null && getBook.withUserId !== '') {
                  global.sqlPool.getConnection(
                    (errorConnectionUpdateData, connectionUpdateData) => {
                      if (!errorConnectionUpdateData) {
                        const sqlUpdateLog = `update log set isReturned = 1, dateReturned = NOW() where libraryId = ${book.libraryId} and bookId = ${getBook.id} and userId = '${userId}';`;
                        const sqlUpdateBook = `update book set withUserId = '' where libraryId = ${book.libraryId} and id = ${getBook.id};`;
                        connectionUpdateData.query(
                          sqlUpdateLog + sqlUpdateBook,
                          (error, resultsReturnBook) => {
                            connectionUpdateData.release();
                            if (!error) {
                              if (resultsReturnBook[0].affectedRows > 0) {
                                if (resultsReturnBook[1].affectedRows > 0) {
                                  res.send({
                                    code: 2000,
                                    data: {
                                      id: 0,
                                      libraryId: book.libraryId,
                                      bookId: getBook.id,
                                      userId,
                                      isReturned: 1,
                                      dateReturned: (new Date()).toUTCString()
                                    }
                                  });
                                }
                                else {
                                  res.statusCode = 500;
                                  res.send(common.customMessage500('Book update failed.'));
                                }
                              }
                              else {
                                res.statusCode = 500;
                                res.send(common.customMessage500('Log add failed.'));
                              }
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.database500(error));
                            }
                          }
                        );
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500(errorConnectionUpdateData));
                      }
                    }
                  );
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500("You don't have it."));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No book.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not a member.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

/* get all books */
router.get('/:libraryId', signIn.verifyRequest, (req, res) => {
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      // eslint-disable-next-line max-len
      // let sql = "select minibook.member.*, minibook.[user].nickname, minibook.[user].avatarurl from minibook.member, minibook.[user] WHERE minibook.member.libraryid = " + req.params.libraryId + " and minibook.member.useropenid = minibook.[user].openid order by minibook.member.isadmin desc, minibook.member.dateupdated OFFSET " + ((page - 1) * pageSize) + " ROWS FETCH NEXT " + pageSize + " ROWS ONLY;";
      // eslint-disable-next-line max-len
      // let count = "select count(*) from minibook.member WHERE libraryid = " + req.params.libraryId + ";";
      const sqlGetBooks = `select * from book WHERE libraryId = ${req.params.libraryId} order by dateUpdated;`;
      const count = '';
      connection.query(sqlGetBooks + count, (error, resultGetBooks) => {
        connection.release();
        if (!error) {
          res.send({
            code: 2000,
            data: resultGetBooks,
            pageSize: 0, // pageSize,
            currentPage: 1, // page,
            totalPage: 1 // Math.ceil(parseInt(result.recordsets[1][0].count) / pageSize),
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

router.get('/:libraryId/:bookId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const book = {
    id: Number(req.params.bookId),
    libraryId: Number(req.params.libraryId)
  };
  global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
    if (!errorConnectionGetData) {
      const sqlGetMember = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId};`;
      const sqlGetBook = `select * from book WHERE libraryid = ${book.libraryId} and id = ${book.id};`;
      const sqlGetLogs = `select * from log WHERE libraryid = ${book.libraryId} and bookId = ${book.id} order by dateUpdated desc;`;
      connectionGetData.query(sqlGetMember + sqlGetBook + sqlGetLogs, (error, resultsGetBook) => {
        connectionGetData.release();
        if (!error) {
          if (resultsGetBook[0].length > 0) {
            const { isAdmin } = resultsGetBook[0][0];
            const currentBook = resultsGetBook[1][0];
            const logs = resultsGetBook[2];
            const doubanGetBoookUrl = common.doubanDataUrl;
            const options = {
              url: doubanGetBoookUrl + currentBook.isbn,
              method: 'GET',
              headers: {
                apikey: common.doubanDataAPIKey
              }
            };
            request(options, (errorRequest, response, body) => {
              if (!errorRequest && response.statusCode < 400) {
                const data = JSON.parse(body);
                res.send({
                  code: 2000,
                  data: {
                    id: book.id,
                    userId: currentBook.userId,
                    libraryId: book.libraryId,
                    title: currentBook.title,
                    isbn: currentBook.isbn,
                    coverUrl: currentBook.coverUrl,
                    withUserId: currentBook.withUserId,
                    dateCreated: currentBook.dateCreated,
                    dateUpdated: currentBook.dateUpdated,
                    logs: isAdmin ? logs : [],
                    details: data
                  }
                });
              }
              else {
                // eslint-disable-next-line no-console
                console.error(errorRequest);
                res.send({
                  code: 2000,
                  data: {
                    id: book.id,
                    userId: currentBook.userId,
                    libraryId: book.libraryId,
                    title: currentBook.title,
                    isbn: currentBook.isbn,
                    coverUrl: currentBook.coverUrl,
                    withUserId: currentBook.withUserId,
                    dateCreated: currentBook.dateCreated,
                    dateUpdated: currentBook.dateUpdated,
                    logs: isAdmin ? logs : [],
                    details: null
                  }
                });
              }
            });
          }
          else {
            res.statusCode = 500;
            res.send(common.customMessage500('Not a member.'));
          }
        }
        else {
          res.statusCode = 500;
          res.send(common.database500(error));
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(errorConnectionGetData));
    }
  });
});

router.get('/isbn/:libraryId/:isbn', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const book = {
    libraryId: req.params.libraryId,
    isbn: req.params.isbn
  };
  global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
    if (!errorConnectionGetData) {
      const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId} and isAdmin = 1`;
      connectionGetData.query(sqlGetAdmin, (error, resultsGetAdmin) => {
        connectionGetData.release();
        if (!error) {
          if (resultsGetAdmin.length > 0) {
            const doubanGetBoookUrl = common.doubanDataUrl;
            const options = {
              url: doubanGetBoookUrl + book.isbn,
              method: 'GET',
              headers: {
                apikey: common.doubanDataAPIKey
              }
            };
            request(options, (errorRequest, response, body) => {
              if (!errorRequest && response.statusCode < 400) {
                const data = JSON.parse(body);
                const { title } = data.items[0].volumeInfo;
                const coverUrl = ''; // data.cover_url;
                res.send({
                  code: 2000,
                  data: {
                    id: 0,
                    userId,
                    title,
                    isbn: book.isbn,
                    coverUrl,
                    details: data
                  }
                });
              }
              else {
                // eslint-disable-next-line no-console
                console.error(errorRequest);
                res.statusCode = 500;
                res.send(common.customMessage500('Load book failed from 3rd party API.'));
              }
            });
          }
          else {
            res.statusCode = 500;
            res.send(common.customMessage500('Not an admin.'));
          }
        }
        else {
          res.statusCode = 500;
          res.send(common.database500(error));
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(errorConnectionGetData));
    }
  });
});

router.put('/manage/return', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.id) {
    const book = {
      libraryId: Number(req.body.libraryId),
      id: Number(req.body.id)
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId} and isAdmin = 1;`;
        const sqlGetLog = `select * from log WHERE libraryId = ${book.libraryId} and bookId = ${book.id} and isReturned = 0;`;
        connectionGetData.query(sqlGetAdmin + sqlGetLog, (errorGetData, resultsGetAdminAndLog) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultsGetAdminAndLog[0].length > 0) {
              if (resultsGetAdminAndLog[1].length > 0) {
                const log = resultsGetAdminAndLog[1][0];
                global.sqlPool.getConnection((errorConnectionUpdateData, connectionUpdateData) => {
                  if (!errorConnectionUpdateData) {
                    const sqlUpdateLog = `update log set isReturned = 1, dateReturned = NOW() where id = ${log.id};`;
                    const sqlUpdateBook = `update book set withUserId = '' where libraryId = ${book.libraryId} and id = ${book.id} and withUserId = '${log.userId}';`;
                    connectionUpdateData.query(
                      sqlUpdateLog + sqlUpdateBook,
                      (error, resultsReturnBook) => {
                        connectionUpdateData.release();
                        if (!error) {
                          if (resultsReturnBook[0].affectedRows > 0) {
                            if (resultsReturnBook[1].affectedRows > 0) {
                              res.send({
                                code: 2000,
                                data: {
                                  id: log.id,
                                  libraryId: log.libraryId,
                                  bookId: log.bookId,
                                  userId: log.userId,
                                  isReturned: 1,
                                  dateReturned: (new Date()).toUTCString()
                                }
                              });
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.customMessage500('Book update failed.'));
                            }
                          }
                          else {
                            res.statusCode = 500;
                            res.send(common.customMessage500('Log add failed.'));
                          }
                        }
                        else {
                          res.statusCode = 500;
                          res.send(common.database500(error));
                        }
                      }
                    );
                  }
                  else {
                    res.statusCode = 500;
                    res.send(common.database500(errorConnectionUpdateData));
                  }
                });
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No log.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not an admin.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.put('/manage/remove', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.id) {
    const book = {
      libraryId: Number(req.body.libraryId),
      id: Number(req.body.id),
      isDeleted: Number(req.body.delete)
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId} and isAdmin = 1;`;
        connectionGetData.query(sqlGetAdmin, (errorGetData, resultGetAdmin) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultGetAdmin.length > 0) {
              global.sqlPool.getConnection((errorDeleteData, connectionDeleteData) => {
                if (!errorDeleteData) {
                  const sqlSetDeleteBook = `update book set isDeleted = ${book.isDeleted} WHERE libraryId = ${book.libraryId} and id = ${book.id};`;
                  connectionDeleteData.query(sqlSetDeleteBook, (error, resultSetDeleteBook) => {
                    connectionDeleteData.release();
                    if (!error) {
                      if (resultSetDeleteBook.affectedRows > 0) {
                        res.send({
                          code: 2000,
                          data: {
                            id: book.id,
                            libraryid: book.libraryId,
                            isDeleted: book.isDeleted
                          }
                        });
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500());
                      }
                    }
                    else {
                      res.statusCode = 500;
                      res.send(common.database500(error));
                    }
                  });
                }
                else {
                  res.statusCode = 500;
                  res.send(common.database500(errorDeleteData));
                }
              });
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not admin.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.delete('/', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.id) {
    const book = {
      id: Number(req.body.id),
      libraryId: Number(req.body.libraryId)
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${book.libraryId} and isAdmin = 1;`;
        connectionGetData.query(sqlGetAdmin, (errorGetData, resultGetAdmin) => {
          connectionGetData.release();
          if (!errorGetData) {
            if (resultGetAdmin.length > 0) {
              global.sqlPool.getConnection((errorDeleteData, connectionDeleteData) => {
                if (!errorDeleteData) {
                  const sqlDeleteBook = `delete from book WHERE libraryId = ${book.libraryId} and id = ${book.id};`;
                  const sqlDeleteLogs = `delete from log WHERE libraryId = ${book.libraryId} and bookId = ${book.id};`;
                  connectionDeleteData.query(
                    sqlDeleteBook + sqlDeleteLogs,
                    (error, resultDeleteBookAndLogs) => {
                      connectionDeleteData.release();
                      if (!error) {
                        if (resultDeleteBookAndLogs[0].affectedRows > 0) {
                          res.send({
                            code: 2000,
                            data: {
                              id: book.id,
                              libraryid: book.libraryId
                            }
                          });
                        }
                        else {
                          res.statusCode = 500;
                          res.send(common.customMessage500('No book.'));
                        }
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500(error));
                      }
                    }
                  );
                }
                else {
                  res.statusCode = 500;
                  res.send(common.database500(errorDeleteData));
                }
              });
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not admin.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetData));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionGetData));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

module.exports = router;
