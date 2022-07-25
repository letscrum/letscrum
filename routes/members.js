const express = require('express');

const router = express.Router();
const signIn = require('../modules/signIn');
const common = require('../modules/common');

/* delete member */
router.delete('/:libraryId/:memberUserId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const member = {
    libraryId: Number(req.params.libraryId),
    userId: req.params.memberUserId ? req.params.memberUserId : userId
  };
  global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
    if (!errorConnectionGetData) {
      const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${member.libraryId} and isAdmin = 1;`;
      const sqlGetMember = `select * from member where userId = '${member.userId}' and libraryId = ${member.libraryId};`;
      connectionGetData.query(sqlGetAdmin + sqlGetMember, (errorGetData, resultGetUser) => {
        connectionGetData.release();
        if (!errorGetData) {
          if (resultGetUser[0].length > 0 || member.userId === userId) {
            if (resultGetUser[1].length > 0) {
              global.sqlPool.getConnection((errorConnectionDeleteData, connectionDeleteData) => {
                if (!errorConnectionDeleteData) {
                  const sqlDeleteMember = `delete from member where userId = '${member.userId}' and libraryId = ${member.libraryId};`;
                  const sqlGetUser = `select * from user where id = '${member.userId}' and libraryId = ${member.libraryId};`;
                  connectionDeleteData.query(
                    sqlDeleteMember + sqlGetUser,
                    (error, resultDeleteMember) => {
                      connectionDeleteData.release();
                      if (!error) {
                        if (resultDeleteMember[0].affectedRows > 0) {
                          if (resultDeleteMember[1].length > 0) {
                            global.sqlPool.getConnection(
                              (errorConnectionUpdateData, connectionUpdateData) => {
                                if (!errorConnectionUpdateData) {
                                  const sqlUpdateUser = `update user set libraryId = 0 where id = '${member.userId}';`;
                                  connectionUpdateData.query(
                                    sqlUpdateUser,
                                    (errorUpdateData, resultsUpdateUser) => {
                                      connectionUpdateData.release();
                                      if (!errorUpdateData) {
                                        if (resultsUpdateUser.affectedRows > 0) {
                                          res.send({
                                            code: 2000,
                                            data: {
                                              userId: member.userId,
                                              libraryId: member.libraryId
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
                                        res.send(common.database500(errorUpdateData));
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
                            res.send({
                              code: 2000,
                              data: {
                                userId: member.userId,
                                libraryId: member.libraryId
                              }
                            });
                          }
                        }
                        else {
                          res.statusCode = 500;
                          res.send(common.database500(error));
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
                  res.send(common.database500(errorConnectionDeleteData));
                }
              });
            }
            else {
              res.statusCode = 500;
              res.send(common.customMessage500('No member.'));
            }
          }
          else {
            res.statusCode = 500;
            res.send(common.customMessage500('Not admin user.'));
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

/* delete member */
router.put('/:libraryId/:memberUserId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.isAdmin === 0 || req.body.isAdmin === 1) {
    const member = {
      libraryId: Number(req.params.libraryId),
      userId: req.params.memberUserId,
      isAdmin: req.body.isAdmin,
      nickname: req.body.nickname ? req.body.nickname : ''
    };
    global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
      if (!errorConnectionGetData) {
        const sqlGetAdmin = `select * from member where userId = '${userId}' and libraryId = ${member.libraryId} and isAdmin = 1;`;
        const sqlGetMember = `select * from member where userId = '${member.userId}' and libraryId = ${member.libraryId};`;
        connectionGetData.query(
          sqlGetAdmin + sqlGetMember,
          (errorGetData, resultGetAdminAndMember) => {
            connectionGetData.release();
            if (!errorGetData) {
              if (resultGetAdminAndMember[0].length > 0) {
                if (resultGetAdminAndMember[1].length > 0) {
                  global.sqlPool.getConnection(
                    (errorConnectionUpdateData, connectionUpdateData) => {
                      if (!errorConnectionUpdateData) {
                        const sqlUpdateMember = `update member set isAdmin = ${member.isAdmin}, nickname = '${member.nickname}', dateUpdated = NOW() where libraryId = ${member.libraryId} and userId = '${member.userId}';`;
                        connectionUpdateData.query(sqlUpdateMember, (error, resultUpdateMember) => {
                          connectionUpdateData.release();
                          if (!error) {
                            if (resultUpdateMember.affectedRows > 0) {
                              res.send({
                                code: 2000,
                                data: {
                                  id: resultGetAdminAndMember[1][0].id,
                                  userId: member.userId,
                                  libraryId: member.libraryId,
                                  nickname: member.nickname,
                                  isAdmin: member.isAdmin,
                                  dateCreated: resultGetAdminAndMember[1][0].dateCreated,
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
                            res.send(common.database500(error));
                          }
                        });
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500(errorConnectionGetData));
                      }
                    }
                  );
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500('No member.'));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('Not admin user.'));
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

/* get all members */
router.get('/:libraryId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      // eslint-disable-next-line max-len
      // let sql = "select minibook.member.*, minibook.[user].nickname, minibook.[user].avatarurl from minibook.member, minibook.[user] WHERE minibook.member.libraryid = " + req.params.libraryId + " and minibook.member.useropenid = minibook.[user].openid order by minibook.member.isadmin desc, minibook.member.dateupdated OFFSET " + ((page - 1) * pageSize) + " ROWS FETCH NEXT " + pageSize + " ROWS ONLY;";
      // eslint-disable-next-line max-len
      // let count = "select count(*) from minibook.member WHERE libraryid = " + req.params.libraryId + ";";
      const sqlGetMembers = `select member.*, user.wechatNickname, user.wechatAvatarUrl from member, user WHERE member.libraryId = ${req.params.libraryId} and member.userId = user.id order by member.isAdmin desc, member.dateUpdated;`;
      const count = '';
      connection.query(sqlGetMembers + count, (error, resultGetMembers) => {
        connection.release();
        if (!error) {
          if (resultGetMembers.length > 0) {
            const selfAdmin = resultGetMembers.find(
              (member) => member.userId === userId && member.isAdmin === 1
            );
            if (selfAdmin) {
              res.send({
                code: 2000,
                data: resultGetMembers,
                pageSize: 0, // pageSize,
                currentPage: 1, // page,
                totalPage: 1 // Math.ceil(parseInt(result.recordsets[1][0].count) / pageSize),
              });
            }
            else {
              res.statusCode = 404;
              res.send(common.customMessage500('Not admin.'));
            }
          }
          else {
            res.statusCode = 404;
            res.send(common.customMessage500('Not member.'));
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
      res.send(common.database500(err));
    }
  });
});

/* add member */
router.post('/', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.libraryId && req.body.userId) {
    const member = {
      userId: req.body.userId,
      libraryId: Number(req.body.libraryId),
      nickname: req.body.nickname ? req.body.nickname : '',
      isAdmin: req.body.isAdmin ? req.body.isAdmin : 0
    };
    global.sqlPool.getConnection((err, connection) => {
      if (!err) {
        const sqlGetLibrary = `select * from member WHERE userId = '${userId}' and libraryId = ${member.libraryId} and isAdmin = 1;`;
        const sqlAddMember = `insert into member(libraryId, userId, nickname, isAdmin) values(${member.libraryId}, '${member.userId}', '${member.nickname}', ${member.isAdmin});`;
        connection.query(sqlGetLibrary + sqlAddMember, (error, resultGetLibraryAndAddMember) => {
          connection.release();
          if (!error) {
            const resultGetLibrary = resultGetLibraryAndAddMember[0];
            const resultAddMember = resultGetLibraryAndAddMember[1];
            if (resultGetLibrary.length > 0) {
              if (resultAddMember.affectedRows > 0) {
                res.send({
                  code: 2000,
                  data: {
                    id: resultAddMember.insertId,
                    libraryId: member.libraryId,
                    userId: member.userId,
                    nickname: member.nickname,
                    isAdmin: member.isAdmin,
                    dateCreated: (new Date()).toUTCString(),
                    dateUpdated: (new Date()).toUTCString()
                  }
                });
              }
              else {
                res.statusCode = 404;
                res.send(common.customMessage500('Add member failed.'));
              }
            }
            else {
              res.statusCode = 404;
              res.send(common.customMessage500('No permission.'));
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
        res.send(common.database500(err));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.get('/:libraryId/:memberUserEmail', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const member = {
    libraryId: Number(req.params.libraryId),
    userEmail: req.params.memberUserEmail
  };
  global.sqlPool.getConnection((errorConnectionGetData, connectionGetData) => {
    if (!errorConnectionGetData) {
      const sqlGetLibrary = `select * from library WHERE id = ${member.libraryId};`;
      const sqlGetAdmin = `select * from member WHERE userId = '${userId}' and libraryId = ${member.libraryId} and isAdmin = 1;`;
      connectionGetData.query(
        sqlGetLibrary + sqlGetAdmin,
        (errorGetLibrary, resultGetLibraryAndAdmin) => {
          connectionGetData.release();
          if (!errorGetLibrary) {
            const resultGetLibrary = resultGetLibraryAndAdmin[0];
            const resultGetAdmin = resultGetLibraryAndAdmin[1];
            if (resultGetLibrary.length > 0) {
              if (resultGetAdmin.length > 0) {
                global.sqlPool.getConnection((errorConnectionGetUser, connectionGetUser) => {
                  if (!errorConnectionGetUser) {
                    const sqlGetUser = `select * from user where email = '${member.userEmail}';`;
                    const sqlGetMember = `select * from member WHERE userId = '${member.userId}' and libraryId = ${member.libraryId};`;
                    connectionGetUser.query(
                      sqlGetUser + sqlGetMember,
                      (errorGetUser, resultsGetUserAndMember) => {
                        connectionGetUser.release();
                        if (!errorGetUser) {
                          const resultGetUser = resultsGetUserAndMember[0];
                          const resultGetMember = resultsGetUserAndMember[1];
                          if (resultGetUser.length > 0) {
                            if (resultGetMember.length <= 0) {
                              res.send({
                                code: 2000,
                                data: {
                                  id: resultGetUser[0].id,
                                  wechatNickname: resultGetUser[0].wechatNickname,
                                  wechatAvatarUrl: resultGetUser[0].wechatAvatarUrl,
                                  username: resultGetUser[0].username,
                                  email: resultGetUser[0].email,
                                  libraryId: Number(resultGetUser[0].libraryId),
                                  isActivated: resultGetUser[0].isActivated,
                                  isDisabled: resultGetUser[0].isDisabled,
                                  isDeleted: resultGetUser[0].isDeleted,
                                  dateCreated: resultGetUser[0].dateCreated,
                                  dateUpdated: resultGetUser[0].dateUpdated
                                }
                              });
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.customMessage500('Already a member.'));
                            }
                          }
                          else {
                            res.statusCode = 404;
                            res.send(common.customMessage500('No user.'));
                          }
                        }
                        else {
                          res.statusCode = 500;
                          res.send(common.database500(errorGetUser));
                        }
                      }
                    );
                  }
                  else {
                    res.statusCode = 500;
                    res.send(common.database500(errorConnectionGetUser));
                  }
                });
              }
              else {
                res.statusCode = 404;
                res.send(common.customMessage500('Not admin.'));
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
      res.send(common.database500(errorConnectionGetData));
    }
  });
});

module.exports = router;
