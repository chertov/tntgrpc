const grpc = require('grpc');
process.exitCode = 1;

const tests = grpc.load(__dirname + '/tests.proto').tests;
const testService = new tests.TestService('localhost:50051', grpc.credentials.createInsecure());
const promises = [];

for (let i = 0; i < 1; i++) {
    promises.push(new Promise((resolve, reject) => {

        const request = {
            outsiders: [
                { name: 'out1' },
                { name: 'out2' },
                { name: 'out3' },
                {
                    name: 'out4',
                    outer: true,
                    mapInMap: {
                        outsider1: { outer: true, name: "outsider 1" },
                        outsider2: { outer: false, name: "outsider 100500!!!" },
                        outsider3: { outer: true, name: "outsider 3" },
                    }
                },
                { name: 'out5' },
            ],
        };
        testService.ItsOk1(request, (err, reply) => {
            if (err) {
                console.log('err.details', err.details);
                reject('testService.ItsOk1');
                return;
            }
            if (reply.enemies["enemy100500"].name !== "outsider 100500!!!") {
                reject('testService.ItsOk1');
                return;
            }
            resolve();
        })
    }));
    promises.push(new Promise((resolve, reject) => {
        testService.ItsOk2({}, (err, reply) => {
            if (err) {
                console.log('err.details', err.details);
                reject('testService.ItsOk2');
                return;
            }
            resolve();
        })
    }));
}

Promise.all(promises)
    .then(() => process.exit(0))
    .catch((ex) => {
        console.log(ex);
    });
