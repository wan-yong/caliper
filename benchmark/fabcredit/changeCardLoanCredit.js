/**
* Changed from 
* Copyright 2017 HUAWEI. All Rights Reserved.
* Author: Yong Wan
* SPDX-License-Identifier: Apache-2.0
**/


'use strict';

const crypto = require('crypto');

module.exports.info  = 'change the Credit Card Loan'; //修改

let bc, contx;
let itemBytes = 1024;   // 修改，对应 config.json的"arguments": {  "itemBytes": 2048000 }
let ids = [];           // save the generated item ids

module.exports.ids = ids;

module.exports.init = function(blockchain, context, args) {
    if(args.hasOwnProperty('itemBytes') ) {
        itemBytes = args.itemBytes;
    }

    bc       = blockchain;
    contx    = context;
    return Promise.resolve();
};

module.exports.run = function() {
    const date   = new Date();
    const today  = (date.getMonth() + 1) + '/' + date.getDate() + '/' + date.getFullYear();
    const author = process.pid.toString();
    const buf    = crypto.randomBytes(itemBytes).toString('base64');
    const item = {
        'author' : author,
        'createtime' : today,
        'info' : '',
        'item' : buf
    };
    return bc.invokeSmartContract(contx, 'drm', 'v0', {verb : 'changeCardLoanCredit', item: JSON.stringify(item)}, 120);
};

module.exports.end = function(results) {
    for (let i in results){
        let stat = results[i];
        if(stat.IsCommitted()) {
            ids.push(stat.result.toString());
        }
    }
    return Promise.resolve();
};

