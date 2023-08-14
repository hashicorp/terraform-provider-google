/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * Background Cloud Function to be triggered by Firestore.
 *
 * @param {object} event The Cloud Functions event.
 * @param {function} callback The callback function.
 */
exports.helloFirestore = function (event, callback) {
    const messageId = event.params.messageId;

    console.log(`Received message ${messageId}`);

    callback();
};
