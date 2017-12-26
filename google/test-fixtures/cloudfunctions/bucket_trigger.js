/**
 * Background Cloud Function to be triggered by Cloud Storage.
 *
 * @param {object} event The Cloud Functions event.
 * @param {function} callback The callback function.
 */
exports.helloGCS = function (event, callback) {
    const file = event.data;

    if (file.resourceState === 'not_exists') {
        console.log(`File ${file.name} deleted.`);
    } else if (file.metageneration === '1') {
        // metageneration attribute is updated on metadata changes.
        // on create value is 1
        console.log(`File ${file.name} uploaded.`);
    } else {
        console.log(`File ${file.name} metadata updated.`);
    }

    callback();
};