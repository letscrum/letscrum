const express = require('express');

const router = express.Router();
const request = require('request');
const { v4: uuidv4 } = require('uuid');
const signIn = require('../modules/signIn');
const common = require('../modules/common');

router.get('/notice', (_req, res) => {
  res.send({
    code: 2000,
    data: {
      show: 1,
      time: 5000,
      message: 'this is a message example, that you can dismiss it when you see it, maybe you can\'t see it, just wait for 5 secends, it will disappared.'
    }
  });
});

/* user signin */
router.get('/signin/wechat', (req, res) => {
  let code = null;
  const url = 'https://api.weixin.qq.com/sns/jscode2session';
  if (!req.query.code) {
    res.statusCode = 500;
    res.send(common.customMessage500('No Wechat code be found.'));
  }
  else {
    code = req.query.code;
    const loginUrl = `${url}?appid=${global.appId}&secret=${global.appSecret}&js_code=${code}&grant_type=authorization_code`;
    let wechatOpenId = null;
    let wechatSessionKey = null;
    request(loginUrl, (errorRequest, response, body) => {
      if (!errorRequest && response.statusCode === 200) {
        const data = JSON.parse(body);
        if (data.errcode) {
          res.statusCode = 500;
          res.send(common.customData500(data));
        }
        else {
          wechatOpenId = data.openid;
          wechatSessionKey = data.session_key;
          global.sqlPool.getConnection((errorConnectionGetUser, connectionGetUser) => {
            if (!errorConnectionGetUser) {
              const sqlGetUser = `select * from user where wechatOpenId = '${wechatOpenId}';`;
              connectionGetUser.query(sqlGetUser, (errorGetUser, resultsGetUser) => {
                connectionGetUser.release();
                if (!errorGetUser) {
                  if (resultsGetUser.length > 0) {
                    const user = resultsGetUser[0];
                    signIn.signedInResponse(
                      user.id,
                      wechatOpenId,
                      user.email,
                      (responseSignIn) => {
                        const resData = {
                          ...responseSignIn.data,
                          wechatSessionKey,
                          wechatNickname: user.wechatNickname,
                          wechatAvatarUrl: user.wechatAvatarUrl,
                          libraryId: user.libraryId
                        };
                        res.statusCode = responseSignIn.statusCode;
                        res.send({
                          code: responseSignIn.code,
                          data: resData
                        });
                      }
                    );
                  }
                  else {
                    global.sqlPool.getConnection((errorConnectionAddUser, connectionAddUser) => {
                      if (!errorConnectionAddUser) {
                        const userId = uuidv4();
                        const sqlAddUser = `insert into user(id, wechatOpenId, wechatSessionKey) values('${userId}', '${wechatOpenId}', '${wechatSessionKey}');`;
                        connectionAddUser.query(sqlAddUser, (errorAddUser, resultsAddUser) => {
                          connectionAddUser.release();
                          if (!errorAddUser) {
                            if (resultsAddUser.affectedRows > 0) {
                              signIn.signedInResponse(
                                userId,
                                wechatOpenId,
                                '',
                                (responseSignIn) => {
                                  const resData = {
                                    ...responseSignIn.data,
                                    wechatSessionKey,
                                    wechatNickname: '',
                                    wechatAvatarUrl: '',
                                    libraryId: 0
                                  };
                                  res.statusCode = responseSignIn.statusCode;
                                  res.send({
                                    code: responseSignIn.code,
                                    data: resData
                                  });
                                }
                              );
                            }
                            else {
                              res.statusCode = 500;
                              res.send(common.database500());
                            }
                          }
                          else {
                            res.statusCode = 500;
                            res.send(common.database500(errorAddUser));
                          }
                        });
                      }
                      else {
                        res.statusCode = 500;
                        res.send(common.database500(errorConnectionAddUser));
                      }
                    });
                  }
                }
                else {
                  res.statusCode = 500;
                  res.send(common.database500(errorGetUser));
                }
              });
            }
            else {
              res.statusCode = 500;
              res.send(common.database500(errorConnectionGetUser));
            }
          });
        }
      }
      else {
        res.statusCode = 500;
        res.send(common.customMessage500('WeChat API error.'));
      }
    });
  }
});

router.post('/signup', (req, res) => {
  if (req.body.type) {
    const signUpType = req.body.type ? req.body.type : 'email';
    if (signUpType === 'email') {
      if (req.body.email && req.body.password && req.body.nickname) {
        const user = {
          id: uuidv4(),
          email: req.body.email,
          password: req.body.password,
          wechatNickname: req.body.nickname
        };
        global.sqlPool.getConnection((errorConnectionGetUser, connectionGetUser) => {
          if (!errorConnectionGetUser) {
            const sqlGetUser = `select * from user where email = '${user.email}';`;
            connectionGetUser.query(sqlGetUser, (errorGetUser, resultsGetUser) => {
              connectionGetUser.release();
              if (!errorGetUser) {
                if (resultsGetUser.length > 0) {
                  res.statusCode = 500;
                  res.send(common.customMessage500('Email used.'));
                }
                else {
                  global.sqlPool.getConnection((err, connection) => {
                    if (!err) {
                      const sqlAddUser = `insert into user(id, wechatNickname, email, password) values('${user.id}', '${user.wechatNickname}', '${user.email}', '${user.password}');`;
                      connection.query(sqlAddUser, (error, resultAddUser) => {
                        connection.release();
                        if (!error) {
                          if (resultAddUser.affectedRows > 0) {
                            signIn.signedInResponse(
                              user.id,
                              '',
                              user.email,
                              (responseSignIn) => {
                                const resData = {
                                  ...responseSignIn.data,
                                  wechatSessionKey: '',
                                  wechatNickname: user.wechatNickname,
                                  wechatAvatarUrl: '',
                                  libraryId: 0
                                };
                                res.statusCode = responseSignIn.statusCode;
                                res.send({
                                  code: responseSignIn.code,
                                  data: resData
                                });
                              }
                            );
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
                      res.send(common.database500(err));
                    }
                  });
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.database500(errorGetUser));
              }
            });
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorConnectionGetUser));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.customMessage500('Request body error.'));
      }
    }
    else {
      res.statusCode = 500;
      res.send(common.customMessage500('Wrong type.'));
    }
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.post('/signin', (req, res) => {
  if (req.body.type) {
    const signUpType = req.body.type ? req.body.type : 'email';
    if (signUpType === 'email') {
      if (req.body.email && req.body.password) {
        const user = {
          email: req.body.email,
          password: req.body.password
        };
        global.sqlPool.getConnection((errorConnectionGetUser, connectionGetUser) => {
          if (!errorConnectionGetUser) {
            const sqlGetUser = `select * from user where email = '${user.email}';`;
            connectionGetUser.query(sqlGetUser, (errorGetUser, resultsGetUser) => {
              connectionGetUser.release();
              if (!errorGetUser) {
                if (resultsGetUser.length > 0) {
                  const currentUser = resultsGetUser[0];
                  if (currentUser.password === user.password) {
                    signIn.signedInResponse(
                      currentUser.id,
                      '',
                      currentUser.email,
                      (responseSignIn) => {
                        const resData = {
                          ...responseSignIn.data,
                          wechatSessionKey: '',
                          wechatNickname: currentUser.wechatNickname,
                          wechatAvatarUrl: currentUser.wechatAvatarUrl,
                          libraryId: currentUser.libraryId
                        };
                        res.statusCode = responseSignIn.statusCode;
                        res.send({
                          code: responseSignIn.code,
                          data: resData
                        });
                      }
                    );
                  }
                  else {
                    res.statusCode = 500;
                    res.send(common.customMessage500('Wrong password.'));
                  }
                }
                else {
                  res.statusCode = 500;
                  res.send(common.customMessage500('No user.'));
                }
              }
              else {
                res.statusCode = 500;
                res.send(common.database500(errorGetUser));
              }
            });
          }
          else {
            res.statusCode = 500;
            res.send(common.database500(errorConnectionGetUser));
          }
        });
      }
      else {
        res.statusCode = 500;
        res.send(common.customMessage500('Request body error.'));
      }
    }
    else {
      res.statusCode = 500;
      res.send(common.customMessage500('Wrong type.'));
    }
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.post('/suggestion', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.email && req.body.suggestion) {
    const suggestion = {
      userId,
      email: req.body.email,
      suggestion: req.body.suggestion
    };
    global.sqlPool.getConnection((err, connection) => {
      if (!err) {
        const sqlAddSuggestion = `insert into suggestion(userId, email, suggestion) values('${suggestion.userId}', '${suggestion.email}', '${suggestion.suggestion}');`;
        connection.query(sqlAddSuggestion, (error, resultAddSuggestion) => {
          connection.release();
          if (!error) {
            if (resultAddSuggestion.affectedRows > 0) {
              res.send({
                code: 2000,
                data: {
                  id: resultAddSuggestion.insertId,
                  userId: suggestion.userId,
                  email: suggestion.email,
                  suggestion: suggestion.suggestion,
                  datecreated: (new Date()).toUTCString(),
                  dateupdated: (new Date()).toUTCString()
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
        res.send(common.database500(err));
      }
    });
  }
  else {
    res.statusCode = 500;
    res.send(common.customMessage500('Request body error.'));
  }
});

router.get('/test', signIn.verifyRequest, (_req, res) => {
  res.send('123');
});

module.exports = router;
