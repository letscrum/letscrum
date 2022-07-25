const express = require('express');

const router = express.Router();
const signIn = require('../modules/signIn');
const common = require('../modules/common');

/* create library */
router.post('/', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (userId && req.body.name && !Number.isNaN(req.body.setDefault)) {
    const library = {
      userId,
      name: req.body.name,
      coverUrl: '',
      description: req.body.description ? req.body.description : ''
    };
    const setDefault = !!req.body.setDefault;
    global.sqlPool.getConnection((errorConnectionAddLibrary, connectionAddLibrary) => {
      if (!errorConnectionAddLibrary) {
        const sqlAddLibrary = `insert into library(userId, name, coverUrl, description) values('${library.userId}', '${library.name}', '${library.coverUrl}', '${library.description}');`;
        connectionAddLibrary.query(sqlAddLibrary, (errorAddLibrary, resultsAddLibrary) => {
          connectionAddLibrary.release();
          if (!errorAddLibrary) {
            if (resultsAddLibrary.affectedRows > 0) {
              global.sqlPool.getConnection((errorConnectionUpdateUser, connectionUpdateUser) => {
                if (!errorConnectionAddLibrary) {
                  let sqlUpdateUser = '';
                  if (setDefault === true) {
                    sqlUpdateUser = `update user set libraryId = ${resultsAddLibrary.insertId} where id = '${library.userId}';`;
                  }
                  const sqlAddMember = `insert into member(libraryId, userId, nickname, isAdmin) values(${resultsAddLibrary.insertId}, '${library.userId}', '', 1);`;
                  connectionUpdateUser.query(sqlUpdateUser + sqlAddMember, (error) => {
                    connectionUpdateUser.release();
                    if (!error) {
                      res.send({
                        code: 2000,
                        data: {
                          id: resultsAddLibrary.insertId,
                          userId: library.userId,
                          name: library.name,
                          code: common.getLibraryJoinCode(
                            library.userId,
                            resultsAddLibrary.insertId
                          ),
                          coverUrl: library.coverUrl,
                          description: library.description,
                          dateCreated: (new Date()).toUTCString(),
                          dateUpdated: (new Date()).toUTCString()
                        }
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
                  res.send(common.database500(errorConnectionUpdateUser));
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
            res.send(common.database500(errorAddLibrary));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.database500(errorConnectionAddLibrary));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

/* get lib details */
router.get('/current', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
    if (!errorConnectionGetData) {
      const sqlGetUser = `select * from user WHERE id = '${userId}';`;
      connectionGetData.query(sqlGetUser, (errorGetData, resultGetUser) => {
        connectionGetData.release();
        if (!errorGetData) {
          if (resultGetUser.length > 0) {
            if (resultGetUser[0].libraryId > 0) {
              global.sqlPool.getConnection((errorConnection, connection) => {
                if (!errorConnection) {
                  const sqlGetLibrary = `select * from library WHERE id = ${resultGetUser[0].libraryId};`;
                  const sqlGetBooks = `select * from book WHERE libraryId = ${resultGetUser[0].libraryId} and isDeleted = 0;`;
                  connection.query(
                    sqlGetLibrary + sqlGetBooks,
                    (error, resultGetLibraryAndBooks) => {
                      connection.release();
                      if (!error) {
                        const resultGetLibrary = resultGetLibraryAndBooks[0];
                        const resultGetBooks = resultGetLibraryAndBooks[1];
                        if (resultGetLibrary.length > 0) {
                          res.send({
                            code: 2000,
                            data: {
                              id: resultGetLibrary[0].id,
                              userId: resultGetLibrary[0].userId,
                              name: resultGetLibrary[0].name,
                              code: common.getLibraryJoinCode(
                                resultGetLibrary[0].userId,
                                resultGetLibrary[0].id
                              ),
                              coverUrl: resultGetLibrary[0].coverUrl,
                              description: resultGetLibrary[0].description,
                              dateCreated: resultGetLibrary[0].dateCreated,
                              dateUpdated: resultGetLibrary[0].dateUpdated,
                              books: resultGetBooks
                            }
                          });
                        }
                        else {
                          res.statusCode = 404;
                          res.send(common.customMessage500('No library.'));
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
                  res.send(common.database500(errorConnection));
                }
              });
            }
            else {
              res.send({
                code: 2000,
                data: {
                  id: 0,
                  userId,
                  name: '',
                  code: '',
                  coverUrl: '',
                  description: '',
                  dateCreated: '',
                  dateUpdated: '',
                  books: []
                }
              });
            }
          }
          else {
            res.statusCode = 404;
            res.send(common.customMessage500('No user.'));
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
});

/* self library list */
router.get('/my', signIn.verifyRequest, (req, res) => {
  /**
  var page = 1;
  var pageSize = 10;
  if (!isNaN(req.query.page)) {
    let pageInt = parseInt(req.query.page);
    if (pageInt > 0) {
      page = pageInt;
    }
  }
  */
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  global.sqlPool.getConnection((errorConnection, connection) => {
    if (!errorConnection) {
      // eslint-disable-next-line max-len
      /** let sql = "select minibook.member.id, minibook.member.libraryid, minibook.library.name, minibook.library.avatarurl, minibook.library.description, minibook.library.useropenid AS owneropenid, minibook.member.useropenid, minibook.member.realname, isadmin, minibook.member.datecreated, minibook.member.dateupdated from minibook.member, minibook.library WHERE minibook.member.useropenid = '" + req.params.openId + "' and minibook.library.id = minibook.member.libraryid order by dateupdated OFFSET " + ((page - 1) * pageSize) + " ROWS FETCH NEXT " + pageSize + " ROWS ONLY;";
      */
      // eslint-disable-next-line max-len
      /** let count = "select count(*) AS count from minibook.member, minibook.library WHERE minibook.member.useropenid = '" + req.params.openId + "' and minibook.library.id = minibook.member.libraryid;";
      */
      const sql = `select member.libraryId AS id, library.userId, library.name, library.coverUrl, library.description, member.id AS memberId, member.userId AS memberUserId, member.nickname AS memberNickname, isAdmin AS memberIsAdmin, library.dateCreated, library.dateUpdated from member, library WHERE member.userId = '${userId}' and library.id = member.libraryId order by dateUpdated;`;
      const count = '';
      connection.query(sql + count, (error, resultGetSelfLibraries) => {
        connection.release();
        if (!error) {
          res.send({
            code: 2000,
            data: resultGetSelfLibraries.map((l) => ({
              id: l.id,
              userId: l.userId,
              name: l.name,
              coverUrl: l.coverUrl,
              description: l.description,
              code: common.getLibraryJoinCode(l.userId, l.id),
              memberId: l.memberId,
              memberUserId: l.memberUserId,
              memberNickname: l.memberNickname,
              memberIsAdmin: l.memberIsAdmin,
              dateCreated: l.dateCreated,
              dateUpdated: l.dateUpdated
            })),
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
      res.send(common.database500(errorConnection));
    }
  });
});

/* self library list */
router.put('/:libraryId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.name) {
    const library = {
      id: Number(req.params.libraryId),
      name: req.body.name,
      description: req.body.description ? req.body.description : ''
    };
    global.sqlPool.getConnection((errorConnection, connection) => {
      if (!errorConnection) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${library.id} and isAdmin = 1;`;
        const sqlUpdateLibrary = `update library set name = '${library.name}', description = '${library.description}' WHERE id = '${library.id}';`;
        connection.query(sqlGetAdmin + sqlUpdateLibrary, (error, resultUpdateLibrary) => {
          connection.release();
          if (!error) {
            if (resultUpdateLibrary[0].length > 0) {
              if (resultUpdateLibrary[1].affectedRows > 0) {
                res.send({
                  code: 2000,
                  data: {
                    id: library.id,
                    name: library.name,
                    description: library.description,
                    userId: '',
                    code: common.getLibraryJoinCode(userId, library.id),
                    coverUrl: ''
                  }
                });
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('Update failed.'));
              }
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('Not admin.'));
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
        res.send(common.database500(errorConnection));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

/* get lib details */
router.get('/:libraryId', signIn.verifyRequest, (req, res) => {
  global.sqlPool.getConnection((errorConnection, connection) => {
    if (!errorConnection) {
      const sqlGetLibrary = `select * from library WHERE id = ${req.params.libraryId};`;
      const sqlGetBooks = `select * from book WHERE libraryId = ${req.params.libraryId} and isDeleted = 0;`;
      connection.query(sqlGetLibrary + sqlGetBooks, (error, resultGetLibraryAndBooks) => {
        connection.release();
        if (!error) {
          const resultGetLibrary = resultGetLibraryAndBooks[0];
          const resultGetBooks = resultGetLibraryAndBooks[1];
          if (resultGetLibrary.length > 0) {
            res.send({
              code: 2000,
              data: {
                id: resultGetLibrary[0].id,
                userId: resultGetLibrary[0].userId,
                name: resultGetLibrary[0].name,
                code: common.getLibraryJoinCode(resultGetLibrary[0].userId, resultGetLibrary[0].id),
                coverUrl: resultGetLibrary[0].coverUrl,
                description: resultGetLibrary[0].description,
                dateCreated: resultGetLibrary[0].dateCreated,
                dateUpdated: resultGetLibrary[0].dateUpdated,
                books: resultGetBooks
              }
            });
          }
          else {
            res.statusCode = 404;
            res.send(common.customMessage500('No library.'));
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
      res.send(common.database500(errorConnection));
    }
  });
});

router.delete('/:libraryId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const library = {
    id: req.params.libraryId,
    userId
  };
  const setDefault = !!req.body.setDefault;
  global.sqlPool.getConnection((errorConnectionAddLibrary, connectionAddLibrary) => {
    if (!errorConnectionAddLibrary) {
      const sqlAddLibrary = `insert into library(userId, name, coverUrl, description) values('${library.userId}', '${library.name}', '${library.coverUrl}', '${library.description}');`;
      connectionAddLibrary.query(sqlAddLibrary, (errorAddLibrary, resultsAddLibrary) => {
        connectionAddLibrary.release();
        if (!errorAddLibrary) {
          if (resultsAddLibrary.affectedRows > 0) {
            global.sqlPool.getConnection((errorConnectionUpdateUser, connectionUpdateUser) => {
              if (!errorConnectionUpdateUser) {
                let sqlUpdateUser = '';
                if (setDefault === true) {
                  sqlUpdateUser = `update user set libraryId = ${resultsAddLibrary.insertId} where id = '${library.userId}';`;
                }
                const sqlAddMember = `insert into member(libraryId, userId, nickname, isAdmin) values(${resultsAddLibrary.insertId}, '${library.userId}', '', 1);`;
                connectionUpdateUser.query(sqlUpdateUser + sqlAddMember, (error) => {
                  connectionUpdateUser.release();
                  if (!error) {
                    res.send({
                      code: 2000,
                      data: {
                        id: resultsAddLibrary.insertId,
                        userId: library.userId,
                        name: library.name,
                        code: common.getLibraryJoinCode(library.userId, resultsAddLibrary.insertId),
                        coverUrl: library.coverUrl,
                        description: library.description,
                        dateCreated: (new Date()).toUTCString(),
                        dateUpdated: (new Date()).toUTCString()
                      }
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
                res.send(common.database500(errorConnectionUpdateUser));
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
          res.send(common.database500(errorAddLibrary));
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(errorConnectionAddLibrary));
    }
  });
});

/* join lib by code */
router.post('/join/:code', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const libraryUserIdAndLibraryId = common.getLibraryByCode(req.params.code);
  const library = {
    userId: libraryUserIdAndLibraryId.split(':')[0],
    id: Number(libraryUserIdAndLibraryId.split(':')[1])
  };
  global.sqlPool.getConnection((errorConnectionGetLibrary, connectionGetLibrary) => {
    if (!errorConnectionGetLibrary) {
      const sqlGetLibrary = `select * from library WHERE userId like '${library.userId}' and id = ${library.id};`;
      const sqlGetMember = `select * from member WHERE userId = '${userId}' and libraryId = ${library.id};`;
      connectionGetLibrary.query(
        sqlGetLibrary + sqlGetMember,
        (errorGetLibrary, resultGetLibraryAndGetMember) => {
          connectionGetLibrary.release();
          if (!errorGetLibrary) {
            const resultGetLibrary = resultGetLibraryAndGetMember[0];
            const resultGetMember = resultGetLibraryAndGetMember[1];
            if (resultGetLibrary.length > 0) {
              if (resultGetMember.length > 0) {
                global.sqlPool.getConnection(
                  (errorConnectionSetCurrentLibrary, connectionSetCurrentLibrary) => {
                    if (!errorConnectionSetCurrentLibrary) {
                      const sqlSetCurrentLibrary = `update user set libraryId = ${library.id} where id = '${resultGetMember[0].userId}';`;
                      connectionSetCurrentLibrary.query(
                        sqlSetCurrentLibrary,
                        (error, resultSetCurrentLibrary) => {
                          connectionSetCurrentLibrary.release();
                          if (!error) {
                            if (resultSetCurrentLibrary.affectedRows > 0) {
                              res.send({
                                code: 2000,
                                data: {
                                  id: resultGetMember[0].id,
                                  libraryId: resultGetMember[0].libraryId,
                                  userId: resultGetMember[0].userId,
                                  nickname: resultGetMember[0].nickname,
                                  isAdmin: resultGetMember[0].isAdmin
                                }
                              });
                            }
                            else {
                              res.statusCode = 404;
                              res.send(common.customMessage500('Add member failed.'));
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
                      res.send(common.database500(errorConnectionSetCurrentLibrary));
                    }
                  }
                );
              }
              else {
                global.sqlPool.getConnection((errorConnectionAddMember, connectionAddMember) => {
                  if (!errorConnectionAddMember) {
                    const sqlAddMember = `insert into member(libraryId, userId, nickname, isAdmin) values(${library.id}, '${userId}', '', 0);`;
                    const sqlSetCurrentLibrary = `update user set libraryId = ${library.id} where id = '${userId}';`;
                    connectionAddMember.query(
                      sqlAddMember + sqlSetCurrentLibrary,
                      (error, resultAddMemberAndSetCurrentLibrary) => {
                        connectionAddMember.release();
                        if (!error) {
                          if (resultAddMemberAndSetCurrentLibrary[0].affectedRows > 0
                            && resultAddMemberAndSetCurrentLibrary[1].affectedRows > 0) {
                            res.send({
                              code: 2000,
                              data: {
                                id: resultAddMemberAndSetCurrentLibrary[1].insertId,
                                libraryId: resultGetLibrary[0].id,
                                userId,
                                nickname: '',
                                isAdmin: 0
                              }
                            });
                          }
                          else {
                            res.statusCode = 404;
                            res.send(common.customMessage500('Add member failed.'));
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
                    res.send(common.database500(errorConnectionAddMember));
                  }
                });
              }
            }
            else {
              res.statusCode = 404;
              res.send(common.customMessage500('No library.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetLibrary));
          }
        }
      );
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(errorConnectionGetLibrary));
    }
  });
});

router.get('/code/:code', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const libraryUserIdAndLibraryId = common.getLibraryByCode(req.params.code);
  const library = {
    userId: libraryUserIdAndLibraryId.split(':')[0],
    id: Number(libraryUserIdAndLibraryId.split(':')[1])
  };
  global.sqlPool.getConnection((errorConnectionGetLibrary, connectionGetLibrary) => {
    if (!errorConnectionGetLibrary) {
      const sqlGetLibrary = `select * from library WHERE userId like '${library.userId}' and id = ${library.id};`;
      const sqlGetMember = `select * from member WHERE userId = '${userId}' and libraryId = ${library.id};`;
      connectionGetLibrary.query(
        sqlGetLibrary + sqlGetMember,
        (errorGetLibrary, resultGetLibraryAndGetMember) => {
          connectionGetLibrary.release();
          if (!errorGetLibrary) {
            const resultGetLibrary = resultGetLibraryAndGetMember[0];
            const resultGetMember = resultGetLibraryAndGetMember[1];
            if (resultGetLibrary.length > 0) {
              if (resultGetMember.length > 0) {
                res.statusCode = 500;
                res.send(common.customMessage500('Already joined.'));
              }
              else {
                res.send({
                  code: 2000,
                  data: {
                    id: resultGetLibrary[0].id,
                    userId: resultGetLibrary[0].userId,
                    name: resultGetLibrary[0].name,
                    code: common.getLibraryJoinCode(
                      resultGetLibrary[0].userId,
                      resultGetLibrary[0].id
                    ),
                    coverUrl: resultGetLibrary[0].coverUrl,
                    description: resultGetLibrary[0].description,
                    dateCreated: resultGetLibrary[0].dateCreated,
                    dateUpdated: resultGetLibrary[0].dateUpdated
                  }
                });
              }
            }
            else {
              res.statusCode = 404;
              res.send(common.customMessage500('No library.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorGetLibrary));
          }
        }
      );
    }
    else {
      res.statusCode = 500;
      res.send(common.database500(errorConnectionGetLibrary));
    }
  });
});

module.exports = router;
