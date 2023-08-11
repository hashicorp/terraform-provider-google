/**
 * HTTP Cloud Function for testing volume mount Secrets.
 */
const fs = require('fs')
exports.echoSecret = (req, res) => {
    const path = '/etc/secrets/test-secret'
    fs.access(path, fs.F_OK, (err) => {
        if (err) {
            console.error(err)
            res.status(200).send(err)
            return
        }
        fs.readFile(path, 'utf8', function (err, data) {
            res.status(200).send("Secret: " + data)

        });
    })
};