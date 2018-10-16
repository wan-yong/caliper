/**
* Copyright 2018 IBM. All Rights Reserved.
*
* SPDX-License-Identifier: Apache-2.0
*
*/

'use strict';

module.exports.info  = 'create stocks on manufacture';

let prodCode;
let cOnManuQty;
let bc, contx;
module.exports.init = function(blockchain, context, args) {
    if(!args.hasOwnProperty('prodCode')) {
        return Promise.reject(new Error('falcon.createOnManu - "prodCode" is missed in the arguments'));
    }

    if(!args.hasOwnProperty('cOnManuQty')) {
        return Promise.reject(new Error('falcon.createOnManu - "cOnManuQty" is missed in the arguments'));
    }

    prodCode = args.prodCode.toString();
    cOnManuQty = args.cOnManuQty.toString();
    bc = blockchain;
    contx = context;
    return Promise.resolve();
};

module.exports.run = function() {
    return bc.invokeSmartContract(contx, 'falcon', 'v0', {verb: 'createOnManu', ProdCode: prodCode, ProdCount: cOnManuQty}, 30);
};

module.exports.end = function() {
    return Promise.resolve();
};


