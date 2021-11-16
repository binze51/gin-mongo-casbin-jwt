// 账号表
db.getCollection("bz_accounts").createIndex({ "accountName": 1 }, { "unique": true }, {
    "background": true
})
//登录状态表
db.getCollection("bz_loginlist").createIndex({ "platform": 1, "accountName": 1, "lock": 1 }, { "unique": true }, {
    "background": true
})