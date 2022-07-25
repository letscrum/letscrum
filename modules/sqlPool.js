var mysql = require('mysql');

var sqlPool = mysql.createPool({
  multipleStatements: true,
  connectionLimit: 100,
  host: 'localhost',
  user: 'root',
  password: process.env.MYSQL_PASSWORD,
  database: 'letscrum'
});
module.exports = sqlPool;
