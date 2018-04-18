const grpc = require('grpc');
process.exitCode = 1;

const tests = grpc.load(__dirname + '/tests.proto').tests;
const testService = new tests.TestService('localhost:50051', grpc.credentials.createInsecure());
const promises = [];

for (let i = 0; i < 1; i++) {
    promises.push(new Promise((resolve, reject) => {
        testService.MapIncorrectKeyType({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'IncorrectType' && grpcerr.stack[1].type === 'MapElementKeyIncorrectType') { resolve(); return; }
            }
            reject('testService.MapIncorrectKeyType');
        })
    }));
    promises.push(new Promise((resolve, reject) => {
        testService.MapIncorrectValueType({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'IncorrectType' && grpcerr.stack[1].type === 'MapElementValueIncorrectType') { resolve(); return; }
            }
            reject('testService.MapIncorrectValueType');
        })
    }));
}

Promise.all(promises)
    .then(() => process.exit(0))
    .catch((ex) => {
        console.log(ex);
    });
