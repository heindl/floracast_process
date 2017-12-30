const glob = require("glob");
const GeoFire = require("geofire");
const admin = require("firebase");
const flags = require('flags');
const path = require('path');
// const LineByLineReader = require('line-by-line');
const fs = require('fs');
admin.initializeApp({
    apiKey: "AIzaSyCEfCQSdpiqEd1DvEiKY8wP6WZrGWqQ0-4",
    authDomain: "floracast-firestore.firebaseapp.com",
    databaseURL: "https://floracast-firestore.firebaseio.com",
    projectId: "floracast-firestore",
    storageBucket: "floracast-firestore.appspot.com",
    messagingSenderId: "1063757818890"
});

const db = admin.database().ref();

flags.defineString('cache_path', '', 'Path to directory of files written by Go process.');

const main = (cache_path) => {

    glob(path.join(cache_path, "/*.jsonl"), {}, function (err, files) {
        if (err) {
            console.error(err)
            process.exit()
        }
        let promises = [];
        for (let i = 0; i < files.length; i++) {
            promises.push(parse_file(files[i]))
        }
        Promise.all(promises).then((d)=>{
            process.exit()
        }).catch((err) => {
            console.error(err);
            process.exit()
        })
    });
};

const remove_past_records = (dateStr, taxonId) => {
    return new Promise((resolve, reject) => {
        const dateRef = db.child("predictions").child(dateStr);

        return dateRef.child(taxonId).remove().then(()=>{

            return dateRef.child("taxa").orderByKey().startAt(taxonId).endAt(taxonId+"~").once("value").then((snapshots)=>{
                let promises = [];
                snapshots.forEach((s) => {
                    promises.push(s.ref.remove())
                });
                Promise.all(promises).then(resolve).catch(reject)
            }).catch(reject)

        }).catch((err)=>{reject(err)});
    })
};

let keysets = 0;
let completed = 0;

const parse_file = (filename) => {
    return new Promise((resolve, reject) => {

        // Should end in .jsonl so remove that and split to get date and taxon.
        const [dateStr, taxonId] = filename.split("/").pop().slice(0, -6).split("-");

        return remove_past_records(dateStr, taxonId).then(()=>{

            const taxonRef = new GeoFire(db.child("predictions").child(dateStr).child(taxonId));
            const taxaRef = new GeoFire(db.child("predictions").child(dateStr).child("taxa"));

            let keySet = {};
            fs.readFileSync(filename).toString().split(/\r|\n/).forEach((line) => {
                if (line.trim().length === 0) {
                    return
                }
                const p = JSON.parse(line);
                let k = `${taxonId},${p.WildernessAreaID},${p.ScaledPredictionValue.toFixed(6).replace(".", "|")},${p.ScarcityValue.toFixed(6).replace(".", "|")}`;
                keySet[k] = [p.Location.latitude, p.Location.longitude]
            });

            keysets += 1
            console.log("keysets", keysets)

            Promise.all([
                taxonRef.set(keySet),
                taxaRef.set(keySet)
            ]).then((d)=>{
                completed +=1
                console.log("completed", completed)
                resolve()
            }).catch(reject)

        }).catch(reject);

    })
};

flags.parse();
main(flags.get('cache_path'));