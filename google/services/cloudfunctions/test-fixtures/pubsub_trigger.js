/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * Background Cloud Function to be triggered by Pub/Sub.
 *
 * @param {object} event The Cloud Functions event.
 * @param {function} callback The callback function.
 */
exports.helloPubSub = function (event, callback) {
    const pubsubMessage = event.data;
    const name = pubsubMessage.data ? Buffer.from(pubsubMessage.data, 'base64').toString() : 'World';

    console.log(`Hello, ${name}!`);

    callback();
};