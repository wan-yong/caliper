/**
* Changed from 
* Copyright
* Author: Yong Wan
* SPDX-License-Identifier: Apache-2.0
**/


'use strict';

const crypto = require('crypto');

module.exports.info  = 'change the Credit Card Loan'; //修改

let bc, contx;
let Key = 'PersonalCredit0' //修改，保存目标的Key值

module.exports.Key = Key;

let CreditCardLoan = '-1000';   // 修改，对应 config.json的"arguments": {  "CreditCardLoan": 2048000 }

module.exports.init = function(blockchain, context, args) {
    if(args.hasOwnProperty('Key') ) {
        Key = args.Key;
    }

    if(args.hasOwnProperty('CreditCardLoan') ) {
        CreditCardLoan = args.CreditCardLoan;
    }

    bc       = blockchain;
    contx    = context;

    bc.invokeSmartContract(contx, 'fabcredit', 'v0', {verb : 'initLedger'}, 120);
    
    return Promise.resolve();
};

module.exports.run = function() {
    
    const item = {
        'Key' : Key,
        'CreditCardLoan' : CreditCardLoan,
    };
    return bc.invokeSmartContract(contx, 'fabcredit', 'v0', {verb : 'changeCardLoanCredit', item: JSON.stringify(item)}, 120);
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

