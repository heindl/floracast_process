const glob = require("glob");
const GeoFire = require("geofire");
// const firebase = require("firebase");
const admin = require("firebase");
admin.initializeApp({
    apiKey: "AIzaSyCEfCQSdpiqEd1DvEiKY8wP6WZrGWqQ0-4",
    authDomain: "floracast-firestore.firebaseapp.com",
    databaseURL: "https://floracast-firestore.firebaseio.com",
    projectId: "floracast-firestore",
    storageBucket: "floracast-firestore.appspot.com",
    messagingSenderId: "1063757818890"
});
const db = admin.database().ref();

// const speciesReader = require('readline').createInterface({
//     input: require('fs').createReadStream('/tmp/predictions-drexhdfzez/species.txt')
// });
//
// const taxaRef = new GeoFire(db.child("Species").child("Predictions"));
//
// speciesReader.on('line', function (line) {
//
//     let obj = JSON.parse(line);
//
//     return taxaRef.set(obj).then(function(){
//         // console.log("done ref predictions")
//     }).catch(function(error) {
//         console.error(error)
//     });
// });


glob("/tmp/predictions-udmugziqny/*.taxon", {}, function (err, files) {
    if (err) {
        return console.error(err)
    }
    for (let i = 0; i<files.length; i++) {
        let taxon = files[i].split("/")[3].slice(0, -6);

        let a = {};
        let lines = require('fs').readFileSync(files[i], 'utf-8').split('\n');

        for (let j = 0; j<lines.length; j++) {
            if (lines[j] === "") {
                continue
            }
           let parsed = JSON.parse(lines[j])
            Object.assign(a, parsed);
        }

        let taxonRef = new GeoFire(db.child("Predictions").child(taxon));
        taxonRef.set(a).then(function(){
            // console.log("done ref predictions")
        }).catch(function(error) {
            console.error(error)
        });
    }

    // files is an array of filenames.
    // If the `nonull` option is set, and nothing
    // was found, then files is ["**/*.js"]
    // er is an error object or null.
});