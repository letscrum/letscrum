const jwt = require('jsonwebtoken');

var signIn = module.exports;

signIn.signedInResponse = (
  userId,
  wechatOpenId,
  email,
  callback
) => {
  const response = {
    statusCode: 200,
    code: 2000,
    data: null
  };
  const payload = {
    userId,
    wechatOpenId,
    email
  };
  jwt.sign(
    payload,
    global.jwtSecret,
    { expiresIn: global.jwtExpiresIn },
    (errorSignIn, accessToken) => {
      if (errorSignIn) {
        response.statusCode = 500;
        response.code = 5000;
        response.data = {
          message: 'Access JWT sign failed.'
        };
        callback(response);
      }
      else {
        jwt.sign(
          payload,
          global.jwtSecret,
          { expiresIn: global.jwtExpiresIn * 2 },
          (errorSignInElse, refreshToken) => {
            if (errorSignInElse) {
              response.statusCode = 500;
              response.code = 5000;
              response.data = {
                message: 'Refresh JWT sign failed.'
              };
              callback(response);
            }
            else {
              response.statusCode = 200;
              response.code = 2000;
              response.data = {
                userId,
                wechatOpenId,
                email,
                accessToken,
                refreshToken
              };
              callback(response);
            }
          }
        );
      }
    }
  );
};

signIn.verifyRequest = (req, res, next) => {
  signIn.verify(req.headers.authorization, (resVerify) => {
    if (resVerify.statusCode !== 200) {
      res.statusCode = resVerify.statusCode;
      res.send({
        code: resVerify.code,
        data: resVerify.data
      });
    }
    else {
      next();
    }
  });
};

signIn.verify = (accessToken, callback) => {
  const response = {
    statusCode: 200,
    code: 2000,
    data: null
  };
  jwt.verify(accessToken, global.jwtSecret, (err, accessDecoded) => {
    if (!err) {
      response.statusCode = 200;
      response.code = 2000;
      response.data = accessDecoded;
      callback(response);
    }
    else if (err.name === 'TokenExpiredError') {
      // eslint-disable-next-line no-console
      console.error(err);
      response.statusCode = 400;
      response.code = 4000;
      response.data = {
        message: 'Access JWT expired.'
      };
      callback(response);
    }
    else {
      // eslint-disable-next-line no-console
      console.error(err);
      response.statusCode = 500;
      response.code = 5000;
      response.data = {
        message: 'Access JWT verify failed.'
      };
      callback(response);
    }
  });
};

signIn.getUserIdInToken = (accessToken, callback) => {
  jwt.verify(accessToken, global.jwtSecret, (err, accessDecoded) => {
    if (!err) {
      callback(accessDecoded.userId);
    }
    else {
      callback('');
    }
  });
};

signIn.refresh = (refreshToken, callback) => {
  let response = {
    statusCode: 200,
    code: 2000,
    data: null
  };
  jwt.verify(refreshToken, global.jwtSecret, (err, refreshDecoded) => {
    if (!err) {
      signIn.signedInResponse(
        refreshDecoded.userId,
        refreshDecoded.wechatOpenId,
        refreshDecoded.email,
        (responseSignIn) => {
          response = responseSignIn;
          callback(response);
        }
      );
    }
    else if (err.name === 'TokenExpiredError') {
      // eslint-disable-next-line no-console
      console.error(err);
      response.statusCode = 400;
      response.code = 4000;
      response.data = {
        message: 'Refresh JWT expired.'
      };
      callback(response);
    }
    else {
      // eslint-disable-next-line no-console
      console.error(err);
      response.statusCode = 500;
      response.code = 5000;
      response.data = {
        message: 'Refresh JWT verify failed.'
      };
      callback(response);
    }
  });
};
