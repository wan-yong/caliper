/**
* Copyright 2017 HUAWEI. All Rights Reserved.
*
* SPDX-License-Identifier: Apache-2.0
*
*/


'use strict';

module.exports.info  = 'querying';

let bc, contx;
let Key;
module.exports.init = function(blockchain, context, args) {
    const changeCardLoanCredit = require('./changeCardLoanCredit.js');
    bc      = blockchain;
    contx   = context;
    Key = changeCardLoanCredit.Key;
    return Promise.resolve();
};

module.exports.run = function() {
    const thiskey  = Key[Math.floor(Math.random()*(Key.length))];
    return bc.queryState(contx, 'fabcredit', 'v0', thiskey);
};

module.exports.end = function(results) {
    return Promise.resolve();
};
