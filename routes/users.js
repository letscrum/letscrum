const express = require('express');
var qrcode = require('qrcode');

const router = express.Router();
const signIn = require('../modules/signIn');
const common = require('../modules/common');

/* update user profile */
router.put('/profile', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  if (req.body.wechatNickname && req.body.wechatAvatarUrl) {
    const user = {
      wechatNickname: req.body.wechatNickname,
      wechatAvatarUrl: req.body.wechatAvatarUrl
    };
    global.sqlPool.getConnection((err, connection) => {
      if (!err) {
        const sql = `update user set wechatNickname = '${user.wechatNickname}', wechatAvatarUrl = '${user.wechatAvatarUrl}' where id = '${userId}';`;
        connection.query(sql, (error, resultUpdateProfile) => {
          connection.release();
          if (!error) {
            if (resultUpdateProfile.affectedRows > 0) {
              res.send({
                code: 2000,
                data: {
                  userId,
                  wechatNickname: user.wechatNickname,
                  wechatAvatarUrl: user.wechatAvatarUrl
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

/* update user profile */
router.put('/library/:libraryId', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  const library = {
    id: req.params.libraryId
  };
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      const sql = `update user set libraryId = ${library.id} where id = '${userId}';`;
      connection.query(sql, (error, resultSetLibrary) => {
        connection.release();
        if (!error) {
          if (resultSetLibrary.affectedRows > 0) {
            res.send({
              code: 2000,
              data: {
                userId,
                librayrId: library.id
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
});

/* get user */
router.get('/verify', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  global.sqlPool.getConnection((err, connection) => {
    if (!err) {
      const sql = `select * from user where id = '${userId}';`;
      connection.query(sql, (error, result) => {
        connection.release();
        if (!error) {
          if (result.length > 0) {
            res.send({
              code: 2000,
              data: result[0]
            });
          }
          else {
            res.statusCode = 500;
            res.send(common.customMessage500('No user.'));
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

/* get user */
router.get('/refresh', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  signIn.refresh(req.headers.authorization, (response) => {
    if (response.code === 2000) {
      global.sqlPool.getConnection((err, connection) => {
        if (!err) {
          const sql = `select * from user where id = '${userId}';`;
          connection.query(sql, (error, result) => {
            connection.release();
            if (!error) {
              if (result.length > 0) {
                res.send({
                  code: 2000,
                  data: {
                    id: result[0].id,
                    wechatOpenId: result[0].wechatOpenId,
                    wechatSessionKey: result[0].wechatSessionKey,
                    wechatNickname: result[0].wechatNickname,
                    wechatAvatarUrl: result[0].wechatAvatarUrl,
                    username: result[0].username,
                    email: result[0].email,
                    password: result[0].password,
                    libraryId: Number(result[0].libraryId),
                    isActivated: result[0].isActivated,
                    isDisabled: result[0].isDisabled,
                    isDeleted: result[0].isDeleted,
                    dateCreated: result[0].dateCreated,
                    dateUpdated: result[0].dateUpdated,
                    accessToken: response.data.accessToken,
                    refreshToken: response.data.refreshToken
                  }
                });
              }
              else {
                res.statusCode = 500;
                res.send(common.customMessage500('No user.'));
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
      res.send(response);
    }
  });
});

router.get('/qr', signIn.verifyRequest, (req, res) => {
  let userId = '';
  signIn.getUserIdInToken(req.headers.authorization, (id) => {
    userId = id;
  });
  qrcode.toDataURL(userId, (err, url) => {
    if (!err) {
      res.send({
        code: 2000,
        data: {
          id: userId,
          idBase64: url
        }
      });
    }
    else {
      res.statusCode = 500;
      res.send(common.customMessage500('QR generate fail.'));
    }
  });
});

module.exports = router;
